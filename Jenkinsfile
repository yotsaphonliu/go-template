#!/usr/bin/env groovy

library(
  identifier: 'service-pipeline@main',
  retriever: modernSCM(
    [
      $class: 'GitSCMSource',
      remote: 'repo.devops.git',
      credentialsId: 'GitLab'
    ]
  )
)

serviceService {
    APP_NAME = 'oa'
}