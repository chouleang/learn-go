pipeline{
    agent any

    environment {
        DOCKER_IMAGE = 'chouleang/go-hello-operator'
        DOCKER_TAG = "build-${BUILD_NUMBER}"
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'main' , url: 'https://github.com/chouleang/learn-go.git'
            }
        }

        stage('Dependencies') {
            steps {
                sh 'go version'
                sh 'go mod download'
                sh 'go mod verify'
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    docker.build("${DOCKER_IMAGE}:${DOCKER_TAG}")
                }
            }
        }

        stage('Push to Docker Hub') {
            steps {
                script {
                    docker.withRegistry('', 'docker-hub-cred') {
                        docker.image("${DOCKER_IMAGE}:${DOCKER_TAG}").push()
                        docker.image("${DOCKER_IMAGE}:latest").push()
                    }
                }
            }
        }

        stage('Test Image') {
            steps {
                script {
                    // Test the built image locally 
                    docker.image("${DOCKER_IMAGE}:${DOCKER_TAG}").inside('-p 8080:8080') {
                        sh 'sleep 5'
                        sh 'curl -f http://localhost:8080/health || exit 1'
                    }
                }
            }
        }
    }

    post {
        always{
            echo 'Pipeline completed'
            cleanWs() // Clean workspace after build
        }
        success {
            echo "SUCCESS: Image ${DOCKER_IMAGE}:${DOCKER_TAG} pushed to Docker Hub"
            sh "echo 'Build ${BUILD_NUMBER} completed successfully'"
        }
        failure {
            echo 'Pipeline failed! Check the logs above'
        }
    }
}
