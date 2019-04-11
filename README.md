# sonarqube-check-quality-gate
Check if quality gates have pass/failed

## variables
this is meant to be run in a pipeline after sonar-scanner outputs the report-task.txt file. It is configured with environment variables

`SONAR_TOKEN`

`CI_PROJECT_DIR` - starting directory. This works with gitlab. 

the expectation is that sonar-scanner runs in your application code directory and writes out `$CI_PROJECT_DIR/.scannerwork/report-task.txt`

## print out some debug
pass a -d

`check-quality-gate -d`

## example run
```
./check-quality-gate 
2019/04/10 22:40:41 Quality gate FAILED with ERROR dashboard available here: https://sonarcloud.io/dashboard?id=sample-project
```
