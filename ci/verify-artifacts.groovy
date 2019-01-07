#!/usr/bin/env groovy

pipeline {
  agent none

  options {
    timeout(time: 6, unit: 'HOURS')
  }

  stages {
    stage('Verify artifacts') {
      parallel {
        stage("Verify Linux artifacts") {
          agent { label 'py36' }

          steps {
            withCredentials([
              [$class: 'StringBinding',
              credentialsId: 'd146870f-03b0-4f6a-ab70-1d09757a51fc',
              variable: 'GITHUB_TOKEN']
            ]) {
                sh '''
                  bash -exc " \
                    cd ci; \
                    python3 -m venv env; \
                    source env/bin/activate; \
                    pip install -r requirements.txt; \
                    ./verify-artifacts.py"
                '''
            }
          }
        }

        stage("Verify macOS artifacts") {
          agent { label 'mac-hh-yosemite' }
          steps {
            withCredentials([
              [$class: 'StringBinding',
              credentialsId: 'd146870f-03b0-4f6a-ab70-1d09757a51fc',
              variable: 'GITHUB_TOKEN']
            ]) {
                sh '''
                  bash -exc " \
                    cd ci; \
                    python3 -m venv env; \
                    source env/bin/activate; \
                    pip install -r requirements.txt; \
                    ./verify-artifacts.py"
                '''
            }
          }
        }

        stage("Verify Windows artifacts") {
          agent {
            node {
              label 'windows'
              customWorkspace 'C:\\windows\\workspace'
            }
          }

          steps {
            withCredentials([
              [$class: 'StringBinding',
              credentialsId: 'd146870f-03b0-4f6a-ab70-1d09757a51fc',
              variable: 'GITHUB_TOKEN']
            ]) {
                bat '''
                  bash -exc " \
                    cd ci; \
                    python -m venv env; \
                    env/Scripts/python -m pip install -U pip; \
                    env/Scripts/pip install -r requirements.txt; \
                    env/Scripts/python ./verify-artifacts.py"
                '''
            }
          }
        }
      }
    }
  }
}
