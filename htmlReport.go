// Copyright 2015 ThoughtWorks, Inc.

// This file is part of getgauge/html-report.

// getgauge/html-report is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// getgauge/html-report is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with getgauge/html-report.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/getgauge/common"
	"github.com/getgauge/gauge/logger"
	"github.com/getgauge/html-report/gauge_messages"
	"github.com/getgauge/html-report/generator"
	"github.com/getgauge/html-report/listener"
)

const (
	reportTemplateDir           = "report-template"
	defaultReportsDir           = "reports"
	gaugeReportsDirEnvName      = "gauge_reports_dir" // directory where reports are generated by plugins
	overwriteReportsEnvProperty = "overwrite_reports"
	resultJsFile                = "result.js"
	htmlReport                  = "html-report"
	SETUP_ACTION                = "setup"
	EXECUTION_ACTION            = "execution"
	GAUGE_HOST                  = "localhost"
	GAUGE_PORT_ENV              = "plugin_connection_port"
	PLUGIN_ACTION_ENV           = "html-report_action"
	timeFormat                  = "2006-01-02 15.04.05"
)

var projectRoot string
var pluginDir string

type nameGenerator interface {
	randomName() string
}

type timeStampedNameGenerator struct {
}

func (T timeStampedNameGenerator) randomName() string {
	return time.Now().Format(timeFormat)
}

func findPluginAndProjectRoot() {
	projectRoot = os.Getenv(common.GaugeProjectRootEnv)
	if projectRoot == "" {
		fmt.Printf("Environment variable '%s' is not set. \n", common.GaugeProjectRootEnv)
		os.Exit(1)
	}

	var err error
	pluginDir, err = os.Getwd()
	if err != nil {
		fmt.Printf("Error finding current working directory: %s \n", err)
		os.Exit(1)
	}
}

func createExecutionReport() {
	os.Chdir(projectRoot)
	listener, err := listener.NewGaugeListener(GAUGE_HOST, os.Getenv(GAUGE_PORT_ENV))
	if err != nil {
		fmt.Println("Could not create the gauge listener")
		os.Exit(1)
	}
	listener.OnSuiteResult(createReport)
	listener.Start()
}

func addDefaultPropertiesToProject() {
	defaultPropertiesFile := getDefaultPropertiesFile()

	reportsDirProperty := &(common.Property{
		Comment:      "The path to the gauge reports directory. Should be either relative to the project directory or an absolute path",
		Name:         gaugeReportsDirEnvName,
		DefaultValue: defaultReportsDir})

	overwriteReportProperty := &(common.Property{
		Comment:      "Set as false if gauge reports should not be overwritten on each execution. A new time-stamped directory will be created on each execution.",
		Name:         overwriteReportsEnvProperty,
		DefaultValue: "true"})

	if !common.FileExists(defaultPropertiesFile) {
		fmt.Printf("Failed to setup html report plugin in project. Default properties file does not exist at %s. \n", defaultPropertiesFile)
		return
	}
	if err := common.AppendProperties(defaultPropertiesFile, reportsDirProperty, overwriteReportProperty); err != nil {
		fmt.Printf("Failed to setup html report plugin in project: %s \n", err)
		return
	}
	fmt.Println("Succesfully added configurations for html-report to env/default/default.properties")
}

func getDefaultPropertiesFile() string {
	return filepath.Join(projectRoot, "env", "default", "default.properties")
}

func createReport(suiteResult *gauge_messages.SuiteExecutionResult) {
	reportsDir := getReportsDirectory(getNameGen())
	err := generator.GenerateReports(suiteResult.GetSuiteResult(), reportsDir)
	if err != nil {
		logger.Fatalf("Failed to generate reports: %s", err.Error())
	}
	err = copyReportTemplateFiles(reportsDir)
	if err != nil {
		logger.Fatalf("Error copying template directory :%s\n", err.Error())
	}
	fmt.Printf("Successfully generated html-report to => %s\n", reportsDir)
}

func getNameGen() nameGenerator {
	var nameGen nameGenerator
	if shouldOverwriteReports() {
		nameGen = nil
	} else {
		nameGen = timeStampedNameGenerator{}
	}
	return nameGen
}

func getReportsDirectory(nameGen nameGenerator) string {
	reportsDir, err := filepath.Abs(os.Getenv(gaugeReportsDirEnvName))
	if reportsDir == "" || err != nil {
		reportsDir = defaultReportsDir
	}
	createDirectory(reportsDir)
	var currentReportDir string
	if nameGen != nil {
		currentReportDir = filepath.Join(reportsDir, htmlReport, nameGen.randomName())
	} else {
		currentReportDir = filepath.Join(reportsDir, htmlReport)
	}
	createDirectory(currentReportDir)
	return currentReportDir
}

func copyReportTemplateFiles(reportDir string) error {
	reportTemplateDir := filepath.Join(pluginDir, reportTemplateDir)
	_, err := common.MirrorDir(reportTemplateDir, reportDir)
	return err
}

func shouldOverwriteReports() bool {
	envValue := os.Getenv(overwriteReportsEnvProperty)
	if strings.ToLower(envValue) == "true" {
		return true
	}
	return false
}

func createDirectory(dir string) {
	if common.DirExists(dir) {
		return
	}
	if err := os.MkdirAll(dir, common.NewDirectoryPermissions); err != nil {
		fmt.Printf("Failed to create directory %s: %s\n", defaultReportsDir, err)
		os.Exit(1)
	}
}
