# sonarqube-check-quality-gate
Check if quality gates have pass/failed

## Variables
this is meant to be run in a pipeline after sonar-scanner outputs the report-task.txt file. It is configured with environment variables

`SONAR_TOKEN`

`SONAR_HOST_URL` - default https://sonarcloud.io

`SONAR_ANALYSIS_URL` - default to SONAR_HOST_URL/api/qualitygates/project_status?analysisId=

`CI_PROJECT_DIR` - starting directory. This works with gitlab. 