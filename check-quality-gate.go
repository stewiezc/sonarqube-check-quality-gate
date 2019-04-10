package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

func main() {
	// pull in environment variables
	sonarToken := getEnv("SONAR_TOKEN", "")
	ciProjectDir := getEnv("CI_PROJECT_DIR", "/tmp/app")
	reportTaskPath := ciProjectDir + "/.scannerwork/report-task.txt"

	// pull in variables from report-task.txt
	_, err := os.Stat(reportTaskPath)
	if err != nil {
		log.Fatal("report-task.txt not found in path:", reportTaskPath)
	}

	sonardata, err := ini.Load(reportTaskPath)
	if err != nil {
		log.Fatal(err)
	}

	organization := sonardata.Section("").Key("organization").String()
	projectKey := sonardata.Section("").Key("projectKey").String()
	serverUrl := sonardata.Section("").Key("serverUrl").String()
	serverVersion := sonardata.Section("").Key("serverVersion").String()
	dashboardUrl := sonardata.Section("").Key("dashboardUrl").String()
	ceTaskId := sonardata.Section("").Key("ceTaskId").String()

	sonarAnalysisUrl := serverUrl + "/api/qualitygates/project_status?analysisId="

	// print out some info if debug
	debug := flag.Bool("d", false, "some debug output")
	flag.Parse()
	if *debug == true {
		fmt.Println("sonarToken:", sonarToken)
		fmt.Println("sonarAnalysisUrl:", sonarAnalysisUrl)
		fmt.Println("ciProjectDir:", ciProjectDir)
		fmt.Println("reportTaskPath:", reportTaskPath)
		fmt.Println("organization:", organization)
		fmt.Println("projectKey:", projectKey)
		fmt.Println("serverUrl:", serverUrl)
		fmt.Println("serverVersion:", serverVersion)
		fmt.Println("dashboardUrl:", dashboardUrl)
		fmt.Println("ceTaskId:", ceTaskId)
	}

	// discover sonar info from report-task.txt file (ini)

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
