pipeline {
    agent any
    
    environment {
        IMAGE_NAME = 'go-hello-operator'
        VAULT_ADDR = 'http://vault.qwerfvcxza.site'  // UPDATE THIS!
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
                        vaultCredentialId: 'b1fd7cdc-25c5-4a31-9e87-3dce22039c11'  // Must match Jenkins credential
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
                    echo "✅ Vault secrets loaded"
                }
            }
        }
        
        stage('Build with Dockerfile') {
            steps {
                script {
                    // THIS IS THE KEY - Let Dockerfile handle everything!
                    // Your Dockerfile already:
                    // 1. Downloads dependencies (go mod download)
                    // 2. Builds the Go binary (go build)
                    // 3. Creates optimized container
                    sh """
                        docker build -t ${IMAGE_NAME}:${env.BUILD_NUMBER} .
                        docker tag ${IMAGE_NAME}:${env.BUILD_NUMBER} ${IMAGE_NAME}:latest
                    """
                    echo "✅ Docker image built using Dockerfile"
                }
            }
        }
        
        stage('Test Container') {
            steps {
                script {
                    sh """
                        # Test that the container works
                        docker run -d --name test-app \
                          -e VAULT_ADDR="${VAULT_ADDR}" \
                          -e VAULT_TOKEN="${VAULT_TOKEN}" \
                          -e ENVIRONMENT="${ENVIRONMENT}" \
                          -p 8080:8080 ${IMAGE_NAME}:${env.BUILD_NUMBER} &
                        sleep 10
                        
                        # Verify the application responds
                        echo "Testing application..."
                        curl -f http://localhost:8080 || exit 1
                        
                        # Cleanup test container
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
                        passwordVariable: 'DOCKER_PASSWORD'
                    )]) {
                        sh """
                            docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
                            docker tag ${IMAGE_NAME}:${env.BUILD_NUMBER} $DOCKER_USERNAME/${IMAGE_NAME}:${env.BUILD_NUMBER}
                            docker tag ${IMAGE_NAME}:latest $DOCKER_USERNAME/${IMAGE_NAME}:latest
                            docker push $DOCKER_USERNAME/${IMAGE_NAME}:${env.BUILD_NUMBER}
                            docker push $DOCKER_USERNAME/${IMAGE_NAME}:latest
                        """
                    }
                }
            }
        }
    }
    
    post {
        always {
            // Safe cleanup without docker commands that might fail
            sh '''
                docker rm -f test-app || true
                echo "Cleanup completed"
            '''
            cleanWs()
        }
        success {
            echo "🎉 Build ${env.BUILD_NUMBER} succeeded!"
        }
        failure {
            echo "💥 Build ${env.BUILD_NUMBER} failed!"
        }
    }
}