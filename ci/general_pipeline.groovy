pipeline {
    agent any

    options {
	gitLabConnection('rest-api-app-gitlab-connection')
	timeout(time: 5, unit: "MINUTES")
	disableConcurrentBuilds()
    }

    environment {
	DOCKER_CREDS = credentials('docker-hub-creds')

	DOCKERHUB_REPO = 'rest-api-app'
	GIT_SHA = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()
	DOCKER_IMAGE_NAME = "${DOCKER_CREDS_USR}/${DOCKERHUB_REPO}:${GIT_SHA}"

	PORT = '8090'
	AUTHOR = 'a.dashchinsky'
	VERSION = "${GIT_SHA}"
    }

    stages {
	stage('lint') {
	    parallel {
		stage('hadolint') {
		    agent { label 'staging' }
		    steps {
			gitlabCommitStatus('lint') {
			    script {
				docker.image('hadolint/hadolint:v2.14.0-debian').inside {
				    sh 'hadolint Dockerfile'
				}
			    }
			}
		    }
		}

		stage ('gosec') {
		    agent { label 'staging'}
		    steps {
			gitlabCommitStatus('gosec') {
			    script {
				docker.image('securego/gosec:2.24.6').inside {
				    sh 'gosec GOFLAGS="-mod=readonly" GONOSUMDB="*" -no-fail ./...'
				}
			    }
			}
		    }
		}
	    }
	}

	stage('build') {
	    agent { label 'staging' }
	    steps {
		gitlabCommitStatus('build') {
		    sh 'docker rm -f ${DOCKERHUB_REPO} || true'
		    sh 'docker build -t ${DOCKER_IMAGE_NAME} .'
		}
	    }
	}

	stage('test') {
	    agent { label 'staging' }
	    stages {
		stage('setup') {
		    steps {
			gitlabCommitStatus('setup') {
			    sh 'docker compose down --volumes --remove-orphans'
			    sh 'docker compose up -d'
			}
		    }
		}

		stage('run tests') {
		    parallel {
			stage ('/info') {
			    steps {
				gitlabCommitStatus('/info') {
				    sh "${WORKSPACE}/tests/info/test_info.sh"
				}
			    }
			}
			stage ('/info/currency') {
			    steps {
				gitlabCommitStatus('/info/currency') {
				    sh "${WORKSPACE}/tests/currency/test_currency.sh"
				}
			    }
			}
		    }
		}
	    }
	}

	stage('dockerhub push') {
	    agent { label 'staging' }
	    when {
	    	branch 'master'
	    }

	    steps {
		gitlabCommitStatus('dockerhub push') {
		    script {
			docker.withRegistry('', 'docker-hub-creds') {
			    sh 'docker push "${DOCKER_IMAGE_NAME}"'
			}
		    }
		}
	    }
	}

	stage('deploy') {
	    agent { label 'production' }

	    environment {
		PRODUCTION_HOST = credentials('PRODUCTION_HOST')
		PRODUCTION_USER = credentials('PRODUCTION_USER')
		DEPLOY_PATH = credentials('PRODUCTION_REST_API_APP_DEPLOY_PATH')
	    }

	    when {
	    	branch 'master'
	    }

	    steps {
		gitlabCommitStatus('deploy') {
		    sh 'docker rm -f ${DOCKERHUB_REPO} || true'
		    script {
			docker.withRegistry('', 'docker-hub-creds') {
			    sh 'docker pull "${DOCKER_IMAGE_NAME}"'
			}
		    }

		    sh 'docker compose -f ${WORKSPACE}/docker-compose.yml down --volumes --remove-orphans'
		    sh 'docker compose -f ${WORKSPACE}/docker-compose.yml up -d'
		}
	    }
	}
    }

    post {
	always {
	    sh 'docker compose down --volumes --remove-orphans'
	    sh 'docker rmi "${DOCKER_IMAGE_NAME}" || true'
	}
	cleanup {
	    cleanWs deleteDirs: true, notFailBuild: true
	}
    }

}
