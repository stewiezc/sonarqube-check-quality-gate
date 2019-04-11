package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/ini.v1"
)

type TaskResponse struct {
	Task struct {
		ID                 string        `json:"id"`
		Type               string        `json:"type"`
		ComponentID        string        `json:"componentId"`
		ComponentKey       string        `json:"componentKey"`
		ComponentName      string        `json:"componentName"`
		ComponentQualifier string        `json:"componentQualifier"`
		AnalysisID         string        `json:"analysisId"`
		Status             string        `json:"status"`
		SubmittedAt        string        `json:"submittedAt"`
		SubmitterLogin     string        `json:"submitterLogin"`
		StartedAt          string        `json:"startedAt"`
		ExecutedAt         string        `json:"executedAt"`
		ExecutionTimeMs    int           `json:"executionTimeMs"`
		Logs               bool          `json:"logs"`
		HasScannerContext  bool          `json:"hasScannerContext"`
		Organization       string        `json:"organization"`
		WarningCount       int           `json:"warningCount"`
		Warnings           []interface{} `json:"warnings"`
	} `json:"task"`
}

type AnalysisResponse struct {
	ProjectStatus struct {
		Status     string `json:"status"`
		Conditions []struct {
			Status         string `json:"status"`
			MetricKey      string `json:"metricKey"`
			Comparator     string `json:"comparator"`
			PeriodIndex    int    `json:"periodIndex"`
			ErrorThreshold string `json:"errorThreshold"`
			ActualValue    string `json:"actualValue"`
		} `json:"conditions"`
		Periods []struct {
			Index     int    `json:"index"`
			Mode      string `json:"mode"`
			Date      string `json:"date"`
			Parameter string `json:"parameter"`
		} `json:"periods"`
		IgnoredConditions bool `json:"ignoredConditions"`
	} `json:"projectStatus"`
}

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
	ceTaskUrl := sonardata.Section("").Key("ceTaskUrl").String()

	// print out some info if debug
	debug := flag.Bool("d", false, "some debug output")
	flag.Parse()
	if *debug == true {
		fmt.Println("sonarToken:", sonarToken)
		fmt.Println("ciProjectDir:", ciProjectDir)
		fmt.Println("reportTaskPath:", reportTaskPath)
		fmt.Println("organization:", organization)
		fmt.Println("projectKey:", projectKey)
		fmt.Println("serverUrl:", serverUrl)
		fmt.Println("serverVersion:", serverVersion)
		fmt.Println("dashboardUrl:", dashboardUrl)
		fmt.Println("ceTaskId:", ceTaskId)
		fmt.Println("ceTaskUrl:", ceTaskUrl)
	}

	// query sonar task to discover the analysis id
	analysisid := getTask(sonarToken, ceTaskUrl)
	sonarAnalysisUrl := serverUrl + "/api/qualitygates/project_status?analysisId=" + analysisid

	if *debug == true {
		fmt.Println("analysisId:", analysisid)
		fmt.Println("sonarAnalysisUrl:", sonarAnalysisUrl)
	}

	// query sonar analysisid to discover if quality gate passed or failed
	gateStatus := getAnalysis(sonarToken, sonarAnalysisUrl)
	if gateStatus == "ERROR" {
		log.Fatal("Quality gate FAILED with ", gateStatus, " dashboard available here: ", dashboardUrl)
	} else {
		fmt.Println("Quality gate PASSED with", gateStatus, "dashboard available here:", dashboardUrl)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getTask(sonarToken string, ceTaskUrl string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ceTaskUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(sonarToken, "")

	var analysisid string
	for done := false; !done; {
		resp, getErr := client.Do(req)
		if getErr != nil {
			log.Fatal(err)
		}

		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		var taskresponse TaskResponse
		jsonErr := json.Unmarshal(body, &taskresponse)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		taskStatus := taskresponse.Task.Status
		if taskStatus == "SUCCESS" {
			analysisid = taskresponse.Task.AnalysisID
			done = true
		} else if taskStatus == "FAILED" {
			log.Fatal("task has failed")
		}
	}

	return analysisid
}

func getAnalysis(sonarToken string, analysisUrl string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", analysisUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(sonarToken, "")

	resp, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var analysisresponse AnalysisResponse
	jsonErr := json.Unmarshal(body, &analysisresponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	gateStatus := analysisresponse.ProjectStatus.Status
	return gateStatus
}
