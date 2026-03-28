pipeline {
    agent any

    options {
	gitLabConnection('rest-api-app-gitlab-connection')
}

    stages {
	stage('lint') {
	    steps {
		updateGitlabCommitStatus name: 'lint', state: 'running'

		echo 'lint stage'
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
