pipeline {
    agent {
        docker {
            image 'google/cloud-sdk:alpine'
            args '-v /var/run/docker.sock:/var/run/docker.sock'
            // Run as root to install packages
            args '-u root:root'
        }
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
                    # Install Go as root
                    apk add --no-cache go
                    go version
                '''
            }
        }
        
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        // ... rest of your stages ...
    }
}