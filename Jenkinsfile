pipeline {
    agent any

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
        
        stage('Install gcloud') {
            steps {
                script {
                    sh '''
                        # Install gcloud reliably
                        set -e
                        
                        echo "Installing Google Cloud SDK..."
                        
                        # Check what OS we're on
                        if [ -f /etc/alpine-release ]; then
                            # Alpine Linux
                            apk add --no-cache curl python3
                            curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-linux-x86_64.tar.gz
                            tar -xf google-cloud-cli-linux-x86_64.tar.gz -C /opt/
                        else
                            # Debian/Ubuntu
                            apt-get update && apt-get install -y curl python3
                            curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-linux-x86_64.tar.gz
                            tar -xf google-cloud-cli-linux-x86_64.tar.gz -C /opt/
                        fi
                        
                        # Install without prompts
                        /opt/google-cloud-sdk/install.sh --quiet --usage-reporting false --command-completion false --path-update false
                        
                        # Verify installation
                        /opt/google-cloud-sdk/bin/gcloud --version
                        echo "Google Cloud SDK installed successfully!"
                    '''
                }
            }
        }
        
        stage('Deploy to GKE Singapore') {
            steps {
                script {
                    withCredentials([file(credentialsId: 'gcp-service-account-key', variable: 'GCP_KEY')]) {
                        sh """
                        # Use full path to gcloud
                        /opt/google-cloud-sdk/bin/gcloud auth activate-service-account --key-file=${GCP_KEY}
                        
                        # Configure kubectl to use our Singapore GKE cluster
                        /opt/google-cloud-sdk/bin/gcloud container clusters get-credentials ${GKE_CLUSTER} --zone ${GKE_ZONE}
                        
                        # Update the deployment with new image
                        kubectl set image deployment/go-hello-operator go-hello-operator=${DOCKER_IMAGE}:${DOCKER_TAG}
                        
                        # Wait for rollout to complete
                        kubectl rollout status deployment/go-hello-operator
                        
                        # Get the service external IP
                        EXTERNAL_IP=\$(kubectl get service go-hello-operator-service -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
                        echo "Application deployed to Singapore GKE"
                        echo "Access your app at: http://\$EXTERNAL_IP"
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
