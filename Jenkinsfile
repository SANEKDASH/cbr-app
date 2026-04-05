pipeline {
    agent {
	label 'docker && linux'
    }

    options {
	gitLabConnection('rest-api-app-gitlab-connection')
	timeout(time: 5, unit: "MINUTES")
	disableConcurrentBuilds()
    }

    environment {
	DOCKER_CREDS = credentials('docker-hub-creds')

	DOCKER_HUBREPO = 'rest-api-app'
	GIT_SHA = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()
	DOCKER_IMAGE_NAME = "${DOCKER_CREDS_USR}/${DOCKERHUB_REPO}:${GIT_SHA}"
    }

    stages {
	stage('lint') {
	    steps {
		script {
		    docker.image('hadolint/hadolint:v2.14.0-debian').inside {
			sh 'hadolint Dockerfile'
		    }
		}
	    }
	}

	stage('build') {
	    steps {
		gitlabCommitStatus('build') {
		    sh 'docker build -t "${DOCKER_IMAGE_NAME}" .'
		}
	    }
	}

	stage('test') {
	    environment {
		PORT = '8090'
		AUTHOR = 'a.dashchinsky'
		VERSION = '1.0.0'
	    }

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
				    sh "chmod +x ${WORKSPACE}/tests/info/test_info.sh"
				    sh "${WORKSPACE}/tests/info/test_info.sh"
				}
			    }
			}
			stage ('/info/currency') {
			    steps {
				gitlabCommitStatus('/info/currency') {
				    sh "chmod +x ${WORKSPACE}/tests/currency/test_currency.sh"
				    sh "${WORKSPACE}/tests/currency/test_currency.sh"
				}
			    }
			}
		    }
		}
	    }
	}

	stage('deploy') {
	    when {
		branch 'master'
	    }
	    steps {
		gitlabCommitStatus('deploy') {
		    script {
			docker.withRegistry('', 'docker-hub-creds') {
			    sh 'docker push "${DOCKER_IMAGE_NAME}"'
			}
		    }
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
