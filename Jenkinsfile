pipeline {
    agent {
        docker {
            image 'google/cloud-sdk:alpine'
            args '-v /var/run/docker.sock:/var/run/docker.sock -u root:root'
        }
    }

    tools {
        go 'go1.21'
    }

    environment {
        DOCKER_IMAGE = 'chouleang/go-hello-operator'
        DOCKER_TAG = "build-${BUILD_NUMBER}"
        GKE_CLUSTER = 'go-hello-cluster'
        GKE_ZONE = 'asia-southeast1-a'
    }

    stages {
        stage('Setup') {
            steps {
                sh '''
                    # Install Go and other dependencies
                    apk add --no-cache go
                    go version
                    echo "All dependencies installed successfully!"
                '''
            }
        }
        
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Dependencies') {
            steps {
                sh '''
                go version
                go mod download
                '''
            }
        }
        
        stage('Build Docker Image') {
            steps {
                sh """
                docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
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
                    echo \$DOCKER_PASSWORD | docker login -u \$DOCKER_USERNAME --password-stdin
                    docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
                    docker logout
                    echo "Successfully pushed to Docker Hub!"
                    """
                }
            }
        }
        
        stage('Deploy to GKE Singapore') {
            steps {
                script {
                    withCredentials([file(credentialsId: 'gcp-service-account-key', variable: 'GCP_KEY')]) {
                        sh """
                        # Authenticate and configure
                        gcloud auth activate-service-account --key-file=${GCP_KEY}
                        gcloud container clusters get-credentials ${GKE_CLUSTER} --zone ${GKE_ZONE}
                        
                        echo "=== Current deployment image ==="
                        kubectl get deployment go-hello-operator -o jsonpath='{.spec.template.spec.containers[0].image}'
                        echo ""
                        
                        echo "=== Updating deployment to use new image ==="
                        # Update the deployment with new image
                        kubectl set image deployment/go-hello-operator go-hello-operator=${DOCKER_IMAGE}:${DOCKER_TAG} --record
                        
                        echo "=== Waiting for rollout ==="
                        kubectl rollout status deployment/go-hello-operator --timeout=300s
                        
                        echo "=== Updated deployment image ==="
                        kubectl get deployment go-hello-operator -o jsonpath='{.spec.template.spec.containers[0].image}'
                        echo ""
                        
                        echo "=== New pods ==="
                        kubectl get pods -l app=go-hello-operator
                        
                        # Get the service external IP
                        EXTERNAL_IP=\$(kubectl get service go-hello-operator-service -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
                        echo "Application deployed to Singapore GKE"
                        echo "Access your app at: http://\$EXTERNAL_IP"
                        echo "New image: ${DOCKER_IMAGE}:${DOCKER_TAG}"
                        """
                    }
                }
            }
        }
    }
    
    post {
        always {
            echo 'Pipeline completed!'
        }
        success {
            echo "SUCCESS: Image ${DOCKER_IMAGE}:${DOCKER_TAG} deployed to Singapore GKE"
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}