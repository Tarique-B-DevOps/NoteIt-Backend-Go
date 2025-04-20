pipeline {
    agent any
    
    environment {
        PROJECT_ID = "${env.GCP_PROJECT_ID}"
        REPOSITORY_NAME = "${env.GCR_NOTEIT_BACKEND_REPO}"
        GCR_URI = "gcr.io/${PROJECT_ID}/${REPOSITORY_NAME}"
    }
    
    stages {
        stage('Build Docker Image with Arguments') {
            steps {
                script {
                    docker.build("${GCR_URI}:latest")
                }
            }
        }

        stage('Scan Image with Trivy') {
            steps {
                script {
                    def trivyOutput = sh(script: "trivy image ${GCR_URI}:latest", returnStdout: true).trim()

                    echo trivyOutput

                    if (trivyOutput.contains("CRITICAL") || trivyOutput.contains("HIGH")) {
                        echo "High or Critical vulnerabilities found in the Docker image."

                        // Uncomment below line to abort pipeline in case of vulnerabilities
                        // error "Aborting the pipeline due to critical or high vulnerabilities."
                    } else {
                        echo "No critical or high vulnerabilities found. Proceeding with the pipeline."
                    }
                }
            }
        }
        
        stage('Push to GCR') {
            steps {
                script {
                    docker.withRegistry('https://gcr.io', 'gcp-staging-project-sa') {
                        docker.image("${GCR_URI}:latest").push()
                    }
                }
            }
        }
    }

    post {
    always {
        slackSend color: "#FFFF00", message: """Docker build started: ${env.JOB_NAME} ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)
        Docker image: ${GCR_URI}:latest
        """
    }
    success {
        slackSend color: "#00FF00", message: "Docker image built and pushed successfully - Job: ${env.JOB_NAME} - Build: ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
    }
    failure {
        slackSend failOnError: true, color: "#FF0000", message: "Docker build/push failed - Job: ${env.JOB_NAME} - Build: ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
    }
}

}
