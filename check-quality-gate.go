package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// pull in environment variables
	sonarToken := getEnv("SONAR_TOKEN", "")
	sonarHostUrl := getEnv("SONAR_HOST_URL", "https://sonarcloud.io")
	sonarAnalysisUrl := getEnv("SONAR_ANALYSIS_URL", sonarHostUrl+"/api/qualitygates/project_status?analysisId=")
	ciProjectDir := getEnv("CI_PROJECT_DIR", "/root/src")

	// print out some info if debug
	debug := flag.Bool("d", false, "some debug output")
	flag.Parse()
	if *debug == true {
		fmt.Println("sonarToken:", sonarToken)
		fmt.Println("sonarHostUrl:", sonarHostUrl)
		fmt.Println("sonarAnalysisUrl:", sonarAnalysisUrl)
		fmt.Println("ciProjectDir:", ciProjectDir)
	}

	// discover sonar info from report-task.txt file (ini)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
