#!groovy

import groovy.json.JsonOutput

helm_chart_version = "0.1.0"

pipeline {
    agent any

    environment {
      GITHUB_ACCESS_TOKEN  = credentials('github-token')
    }

    stages {
    	stage('Checkout Code') {
          steps {
            checkout scm
          }
        }

        stage('Build the deploy image') {
            steps {
                sh """
                    docker build -t apisim .
                    docker build -t apisim-migrate . -f Dockerfile-migrate
                """
            }
        }

        stage('Publish to docker registry') {
            steps {
                container('runner') {
                    withDockerRegistry([credentialsId: 'docker-registry-token', url: 'https://index.docker.io/v1/']) {
                        sh """
                            docker tag apisim ektowett/apisim:${env.BRANCH_NAME}-${GIT_COMMIT.take(10)}
                            docker tag apisim-migrate ektowett/apisim-migrate:${env.BRANCH_NAME}-${GIT_COMMIT.take(10)}

                            docker push ektowett/apisim:${env.BRANCH_NAME}-${GIT_COMMIT.take(10)}
                            docker push ektowett/apisim-migrate:${env.BRANCH_NAME}-${GIT_COMMIT.take(10)}
                        """
                    }
                }
            }
        }

        // stage('Deploy service') {
        //     steps {
        //         container('runner') {
        //             dir("backend") {
        //                 sh """
        //                     curl -LO https://kip0127-helm.s3.eu-west-1.amazonaws.com/app-0.1.0.tgz
        //                     tar -xzf app-0.1.0.tgz
        //                     helm upgrade -i --debug apisim ./app \
        //                         --version ${helm_chart_version} \
        //                         --set image.tag=${env.BRANCH_NAME}-${GIT_COMMIT.take(10)} \
        //                         --set hook.image.tag=${env.BRANCH_NAME}-${GIT_COMMIT.take(10)} \
        //                         -f ./helm/${env.ENV}.yaml
        //                     kubectl rollout status deployment.apps/apisim
        //                 """
        //             }
        //         }
        //     }
        // }
    }
}
