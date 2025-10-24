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
<<<<<<< Updated upstream
                sh """
                echo "=== Building Docker Image ==="
                docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
                echo "Image built: ${DOCKER_IMAGE}:${DOCKER_TAG}"
                """
=======
                script {
                    sh """
                    docker --version
                    docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
                    """
                }
>>>>>>> Stashed changes
            }
        }
        
        stage('Push to Docker Hub') {
            steps {
<<<<<<< Updated upstream
                withCredentials([usernamePassword(
=======
                script {
                    withCredentials([usernamePassword(
>>>>>>> Stashed changes
                    credentialsId: 'docker-hub-cred',
                    usernameVariable: 'DOCKER_USERNAME',
                    passwordVariable: 'DOCKER_PASSWORD'
                )]) {
                    sh """
<<<<<<< Updated upstream
                    echo "=== Pushing to Docker Hub ==="
                    echo "Logging in as \$DOCKER_USERNAME..."
                    echo \$DOCKER_PASSWORD | docker login -u \$DOCKER_USERNAME --password-stdin
                    
                    echo "Pushing images..."
                    docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
                    
                    echo "Logging out..."
                    docker logout
                    echo "✅ Push completed!"
=======
                    echo \$DOCKER_PASSWORD | docker login -u \$DOCKER_USERNAME --password-stdin
                    docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
                    docker logout
                    echo "✅ Successfully pushed to Docker Hub!"
>>>>>>> Stashed changes
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
