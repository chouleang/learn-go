pipeline {
    agent any
    
    environment {
        IMAGE_NAME = 'go-hello-operator'
        VAULT_ADDR = 'http://vault.qwerfvcxza.site'  // Replace with your actual Vault URL
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Load Vault Secrets') {
            steps {
                withVault(
                    configuration: [
                        vaultUrl: "${VAULT_ADDR}",
                        vaultCredentialId: 'b1fd7cdc-25c5-4a31-9e87-3dce22039c11'  // Name of your Jenkins credential
                    ],
                    vaultSecrets: [
                        [
                            path: 'secret/jenkins/go-operator',
                            secretValues: [
                                [envVar: 'VAULT_TOKEN', vaultKey: 'vault-token'],
                                [envVar: 'ENVIRONMENT', vaultKey: 'environment'],
                                [envVar: 'DOCKER_PASSWORD', vaultKey: 'docker-password']
                            ]
                        ]
                    ]
                ) {
                    echo "✅ Vault secrets loaded:"
                    sh 'echo "Environment: $ENVIRONMENT"'
                    sh 'echo "Docker password length: ${#DOCKER_PASSWORD}"'
                    // Don't echo the actual token for security
                }
            }
        }
        
        stage('Build and Test') {
            steps {
                sh '''
                    echo "Building Go application..."
                    go mod download
                    go test ./... -v
                    go build -o main .
                '''
            }
        }
        
        stage('Build Docker Image') {
            steps {
                script {
                    // Build with Vault environment variables
                    sh """
                        docker build -t ${IMAGE_NAME}:${env.BUILD_NUMBER} .
                        docker tag ${IMAGE_NAME}:${env.BUILD_NUMBER} ${IMAGE_NAME}:latest
                    """
                }
            }
        }
        
        stage('Test Container') {
            steps {
                script {
                    sh """
                        docker run -d --name test-app \
                          -e VAULT_ADDR="${VAULT_ADDR}" \
                          -e VAULT_TOKEN="${VAULT_TOKEN}" \
                          -e ENVIRONMENT="${ENVIRONMENT}" \
                          -p 8080:8080 ${IMAGE_NAME}:${env.BUILD_NUMBER} &
                        sleep 15
                        echo "Testing container..."
                        curl -f http://localhost:8080 || exit 1
                        docker stop test-app
                        docker rm test-app
                    """
                }
            }
        }
        
        stage('Push to Registry') {
            when {
                branch 'main'
            }
            steps {
                script {
                    withCredentials([usernamePassword(
                        credentialsId: 'dockerhub-creds',
                        usernameVariable: 'DOCKER_USERNAME',
                        passwordVariable: 'DOCKER_PASSWORD'  // This comes from Vault!
                    )]) {
                        sh """
                            docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
                            docker tag ${IMAGE_NAME}:${env.BUILD_NUMBER} ${DOCKER_USERNAME}/${IMAGE_NAME}:${env.BUILD_NUMBER}
                            docker push ${DOCKER_USERNAME}/${IMAGE_NAME}:${env.BUILD_NUMBER}
                        """
                    }
                }
            }
        }
    }
    
    post {
        always {
            sh '''
                docker rm -f test-app || true
                docker rmi ${IMAGE_NAME}:${env.BUILD_NUMBER} || true
                cleanWs()
            '''
        }
        success {
            echo "🎉 Build ${env.BUILD_NUMBER} succeeded!"
        }
        failure {
            echo "💥 Build ${env.BUILD_NUMBER} failed!"
        }
    }
}