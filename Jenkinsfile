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
                    docker.withRegistry('https://gcr.io', 'gcp-devops-proj-servive-account') {
                        docker.image("${GCR_URI}:latest").push()
                    }
                }
            }
        }
    }
}
