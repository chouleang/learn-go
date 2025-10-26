pipeline {
    agent any
    
    stages {
        stage('Test Docker Commands') {
            steps {
                script {
                    // Test 1: Run commands inside a container
                    docker.image('alpine:latest').inside {
                        sh 'echo "Hello from Alpine container!"'
                        sh 'cat /etc/os-release'
                    }
                    
                    // Test 2: Pull and inspect image
                    def testImage = docker.image('hello-world:latest')
                    testImage.pull()
                    echo "Image ID: ${testImage.id}"
                }
            }
        }
    }
}