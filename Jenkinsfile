pipeline {
    agent {
	label 'docker && linux'
    }

    options {
	gitLabConnection('rest-api-app-gitlab-connection')
    }

    environment {
	DOCKERHUB_USER = 'sanekdash'
	DOCKERHUB_REPO = 'rest-api-app'

	GIT_SHA = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()

	DOCKER_IMAGE_TAG = "${env.BUILD_NUMBER}"
	DOCKER_IMAGE_NAME = "${DOCKERHUB_USER}/${DOCKERHUB_REPO}:${GIT_SHA}"
    }

    stages {
	stage('lint') {
	    steps {
		updateGitlabCommitStatus name: 'lint', state: 'running'

		script {
		    def lintRes = sh(script: ''' docker run --rm \
				     -v ${WORKSPACE}:${WORKSPACE} \
				     -w ${WORKSPACE} \
				     hadolint/hadolint:latest-debian \
				     hadolint Dockerfile ''',
				     returnStatus: true)

		    if (lintRes != 0) {
			error('Hadolint check failed.')
		    } else {
			echo 'Hadolint check passed.'
		    }

		}
	    }
	    post {
		success {
		    updateGitlabCommitStatus name: 'lint', state: 'success'
		}

		failure {
		    updateGitlabCommitStatus name: 'lint', state: 'failed'
		}
	    }

	}

	stage('build') {
	    steps {
		updateGitlabCommitStatus name: 'build', state: 'running'
		sh 'docker build -t ${DOCKER_IMAGE_NAME} .'
	    }

	    post {
		success {
		    updateGitlabCommitStatus name: 'build', state: 'success'
		}

		failure {
		    updateGitlabCommitStatus name: 'build', state: 'failed'
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
			sh 'docker compose down --volumes --remove-orphans'
			sh 'docker compose up -d'
		    }
		    post {
			success {
			    updateGitlabCommitStatus name: 'setup', state: 'success'
			}
			failure {
			    updateGitlabCommitStatus name: 'setup', state: 'failed'
			}
		    }
		}

		stage('run tests') {
		    parallel {
			stage ('/info') {
			    steps {
				updateGitlabCommitStatus name: '/info', state: 'running'
				sh 'chmod +x ./tests/info/test_info.sh'
				sh './tests/info/test_info.sh'
			    }
			    post {
				success {
				    updateGitlabCommitStatus name: 'info', state: 'success'
				}
				failure {
				    updateGitlabCommitStatus name: 'info', state: 'failed'
				}
			    }
			}
			stage ('/info/currency') {
			    steps {
				updateGitlabCommitStatus name: '/info/currency', state: 'running'
				sh 'chmod +x ./tests/currency/test_currency.sh'
				sh './tests/currency/test_currency.sh'
			    }
			    post {
				success {
				    updateGitlabCommitStatus name: '/info/currency', state: 'success'
				}
				failure {
				    updateGitlabCommitStatus name: '/info/currency', state: 'failed'
				}
			    }
			}
		    }
		}
	    }

	    post {
		always {
		    sh 'docker compose down --volumes --remove-orphans'
		}
	    }
	}

	stage('docker-hub push') {
	    // when {
	    // 	branch 'master'
	    // }
	    steps {
		updateGitlabCommitStatus name: 'docker-hub push', state: 'running'
		script {
		    withCredentials([usernamePassword(
			credentialsId: 'docker-hub-creds',
			usernameVariable: 'DOCKER_USER',
			passwordVariable: 'DOCKER_PASS')]) {
			sh '''
	                    echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin
        	            docker push ${DOCKER_IMAGE_NAME}
                	    docker logout
                	'''
		    }
		}
	    }

	    post {
		success {
		    updateGitlabCommitStatus name: 'deploy', state: 'success'
		}

		failure {
		    updateGitlabCommitStatus name: 'deploy', state: 'failed'
		}
	    }
	}

	stage('deploy') {
	    // when {
	    // 	branch 'master'
	    // }

	    steps {
		echo 'deploy stage'
	    }

	    post {
		success {
		    updateGitlabCommitStatus name: 'deploy', state: 'success'
		}

		failure {
		    updateGitlabCommitStatus name: 'deploy', state: 'failed'
		}
	    }
	}

    }
}
