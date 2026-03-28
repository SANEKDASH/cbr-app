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
			error('Hadolint check failed')
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
		echo 'build stage'
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
	    steps {
		updateGitlabCommitStatus name: 'test', state: 'running'

		echo 'test stage'
	    }
	    post {
		success {
		    updateGitlabCommitStatus name: 'test', state: 'success'
		}
		failure {
		    updateGitlabCommitStatus name: 'test', state: 'failed'
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
