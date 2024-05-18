pipeline {
    agent any

    stages {
        stage('Prepare') {
            steps {
                deleteDir()  // Ensures the workspace is clean before starting the build
            }
        }
        stage('Clone Repository') {
            steps {
                git branch: 'main', url: 'https://github.com/angelhvargas/redfishcli.git'
            }
        }
        stage('Unit Tests') {
            steps {
                sh 'make test', label: 'Running unit tests'
                sh 'make integration-test', label: 'Running integration tests'
            }
        }

        stage('Integration Tests') {
            environment {
                BMC_HOSTNAME = 'test-bmc-hostname'
                BMC_USERNAME = 'test-bmc-username'
                BMC_PASSWORD = 'test-bmc-password'
            }
            steps {
                sh 'go test -tags=integration -v ./...'
            }
        }

        stage('Build') {
            steps {
                // No longer use 'dir' to change to a non-existent subdirectory
                sh 'pwd'
                sh 'ls -la'  // List the files to confirm structure (this can be removed later)
                sh 'go build -o redfishcli main.go'  // Build using the main.go at this level
            }
        }
        stage('Execute') {
            steps {
                sh './redfishcli'  // Execute the binary
            }
        }
    }
    post {
        always {
            cleanWs()  // Clean the workspace after every run
        }
    }
}
