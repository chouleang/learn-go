pipeline {
    agent any

    tools {
        go 'go1.21'
    }

    environment {
        DOCKER_IMAGE = 'chouleang/go-hello-operator'
        DOCKER_TAG = "build-${BUILD_NUMBER}"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Dependencies') {
            steps {
                sh '''
                echo "=== Environment Setup ==="
                go version
                go mod download
                '''
            }
        }
        
        stage('Build Docker Image') {
            steps {
                sh """
                echo "=== Building Docker Image ==="
                docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
                echo "Image built: ${DOCKER_IMAGE}:${DOCKER_TAG}"
                """
            }
        }
        
        stage('Push to Docker Hub') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'docker-hub-cred',
                    usernameVariable: 'DOCKER_USERNAME',
                    passwordVariable: 'DOCKER_PASSWORD'
                )]) {
                    sh """
                    echo "=== Pushing to Docker Hub ==="
                    echo "Logging in as \$DOCKER_USERNAME..."
                    echo \$DOCKER_PASSWORD | docker login -u \$DOCKER_USERNAME --password-stdin
                    
                    echo "Pushing images..."
                    docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
                    
                    echo "Logging out..."
                    docker logout
                    echo "✅ Push completed!"
                    """
                }
            }
        }
    }
    
    post {
        always {
            echo '=== Pipeline Completed ==='
            cleanWs()
        }
        success {
            echo "✅ SUCCESS: Image ${DOCKER_IMAGE}:${DOCKER_TAG} pushed to Docker Hub"
        }
        failure {
            echo '❌ Pipeline failed!'
        }
    }
}
