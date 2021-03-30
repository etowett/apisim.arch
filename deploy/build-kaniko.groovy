#!groovy

import groovy.json.JsonOutput

pipeline {
    agent {
        kubernetes {
            label "build-apisim-${BUILD_NUMBER}"
            defaultContainer 'jnlp'
            yaml """
apiVersion: v1
kind: Pod
metadata:
  labels:
    jenkins-build: apisim-build
    some-label: "build-apisim-${BUILD_NUMBER}"
spec:
  containers:
  - name: kaniko
    image: gcr.io/kaniko-project/executor:v1.5.1-debug
    imagePullPolicy: IfNotPresent
    command:
    - /busybox/cat
    tty: true
    volumeMounts:
      - name: jenkins-docker-cfg
        mountPath: /kaniko/.docker
  - name: runner
    image: ektowett/jenkins-slave:v0.1.0
    command:
    - cat
    tty: true
  volumes:
  - name: jenkins-docker-cfg
    projected:
      sources:
      - secret:
          name: docker-credentials
          items:
            - key: .dockerconfigjson
              path: config.json
"""
        }
    }

    environment {
        GITHUB_ACCESS_TOKEN  = credentials('github-token')
    }

    stages {

        stage('Checkout Code') {
            steps {
              checkout scm
            }
        }

        stage('Build with Kaniko') {
          steps {
            container(name: 'kaniko', shell: '/busybox/sh') {
              withEnv(['PATH+EXTRA=/busybox']) {
                sh '''#!/busybox/sh -xe
                  /kaniko/executor \
                    --dockerfile Dockerfile \
                    --context `pwd`/ \
                    --verbosity debug \
                    --insecure \
                    --skip-tls-verify \
                    --destination ektowett/apisim:latest

                  /kaniko/executor \
                    --dockerfile Dockerfile-migrate \
                    --context `pwd`/  \
                    --verbosity debug \
                    --insecure \
                    --skip-tls-verify \
                    --destination ektowett/apisim-migrate:latest
                '''
              }
            }
          }
        }
    }
}
