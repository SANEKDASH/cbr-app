pipeline {
    agent any

    options {
	gitLabConnection('rest-api-app-gitlab-connection')
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
		sh 'docker compose build'
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
	    parallel {
		stage ('/info') {
		    steps {
			updateGitlabCommitStatus name: '/info', state: 'running'
			sh 'docker compose down --volumes --remove-orphans'
			sh 'docker compose up -d'
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
			always {
			    sh 'docker compose down --volumes --remove-orphans'
			}
		    }
		}
		stage ('/info/currency') {
		    steps {
			echo 'running /info/currency tests'
		    }
		    post {
			success {
			    updateGitlabCommitStatus name: '/info/currency', state: 'success'
			}
			failure {
			    updateGitlabCommitStatus name: '/info/currency', state: 'failed'
			}
			always {
			    sh 'docker compose down --volumes --remove-orphans'
			}
		    }

		}
	    }
	}

	stage('deploy') {
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
