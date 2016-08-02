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

package generator

import (
	"bytes"
	"regexp"
	"testing"
)

type reportGenTest struct {
	name   string
	tmpl   string
	input  interface{}
	output string
}

var wBodyHeader string = `<header class="top">
<div class="header">
  <div class="container">
     <div class="logo"><img src="images/logo.png" alt="Report logo"></div>
        <h2 class="project">Project: projname</h2>
      </div>
  </div>
</header>`

var wChartDiv string = `<div class="report-overview">
  <div class="report_chart">
    <div class="chart">
      <nvd3 options="options" data="data"></nvd3>
    </div>
    <div class="total-specs"><span class="value">41</span><span class="txt">Total specs</span></div>
  </div>`

var wResCntDiv string = `
  <div class="report_test-results">
    <ul>
      <li class="fail"><span class="value">2</span><span class="txt">Failed</span></li>
      <li class="pass"><span class="value">39</span><span class="txt">Passed</span></li>
      <li class="skip"><span class="value">0</span><span class="txt">Skipped</span></li>
    </ul>
  </div>`

var wEnvLi string = `<div class="report_details"><ul>
      <li>
        <label>Environment </label>
        <span>default</span>
      </li>`

var wTagsLi string = `
      <li>
        <label>Tags </label>
        <span>foo</span>
      </li>`

var wSuccRateLi string = `
      <li>
        <label>Success Rate </label>
        <span>34%</span>
      </li>`

var wExecTimeLi string = `
     <li>
        <label>Total Time </label>
        <span>00:01:53</span>
      </li>`

var wTimestampLi string = `
     <li>
        <label>Generated On </label>
        <span>Jun 3, 2016 at 12:29pm</span>
      </li>
    </ul>
  </div>
</div>`

var wSidebarAside string = `<aside class="sidebar">
  <h3 class="title">Specifications</h3>
  <div class="searchbar">
    <input id="searchSpecifications" placeholder="Type specification or tag name" type="text" />
    <i class="fa fa-search"></i>
  </div>
  <div id="listOfSpecifications">
    <ul id="scenarios" class="spec-list">
		<a href="passing_spec.html">
    	<li class='passed spec-name'>
	      <span id="scenarioName" class="scenarioname">Passing Spec</span>
	      <span id="time" class="time">00:01:04</span>
    	</li>
		</a>
		<a href="failing_spec.html">
    	<li class='failed spec-name'>
	      <span id="scenarioName" class="scenarioname">Failing Spec</span>
	      <span id="time" class="time">00:00:30</span>
    	</li>
		</a>
		<a href="skipped_spec.html">
    	<li class='skipped spec-name'>
	      <span id="scenarioName" class="scenarioname">Skipped Spec</span>
	      <span id="time" class="time">00:00:00</span>
    	</li>
		</a>
    </ul>
  </div>
</aside>`

var wCongratsDiv string = `<div class="congratulations details">
  <p>Congratulations! You've gone all <span class="green">green</span> and saved the environment!</p>
</div>`

var wHookFailureWithScreenhotDiv string = `<div class="error-container failed">
<div collapsable class="error-heading">BeforeSuite Failed:<span class="error-message"> SomeError</span></div>
  <div class="toggleShow" data-toggle="collapse" data-target="#hookFailureDetails">
    <span>[Show details]</span>
  </div>
  <div class="exception-container" id="hookFailureDetails">
      <div class="exception">
        <pre class="stacktrace">Stack trace</pre>
      </div>
      <div class="screenshot-container">
        <a href="data:image/png;base64,iVBO" rel="lightbox">
          <img src="data:image/png;base64,iVBO" class="screenshot-thumbnail" />
        </a>
      </div>
  </div>
</div>`

var wHookFailureWithoutScreenhotDiv string = `<div class="error-container failed">
  <div collapsable class="error-heading">BeforeSuite Failed:<span class="error-message"> SomeError</span></div>
  <div class="toggleShow" data-toggle="collapse" data-target="#hookFailureDetails">
    <span>[Show details]</span>
  </div>
  <div class="exception-container" id="hookFailureDetails">
      <div class="exception">
        <pre class="stacktrace">Stack trace</pre>
      </div>
  </div>
</div>`

var wSpecHeaderStartWithTags string = `<header class="curr-spec">
  <h3 class="spec-head" title="/tmp/gauge/specs/foobar.spec">Spec heading</h3>
  <span class="time">00:01:01</span>`

var wTagsDiv string = `<div class="tags scenario_tags contentSection">
  <strong>Tags:</strong>
  <span> tag1</span>
  <span> tag2</span>
</div>`

var wSpecCommentsWithTableTag string = `<span></span>
<span>This is an executable specification file. This file follows markdown syntax.</span>
<span></span>
<span>To execute this specification, run</span>
<span>gauge specs</span>
<span></span>
<table class="data-table">
  <tr>
    <th>Word</th>
    <th>Count</th>
  </tr>
  <tbody>
    <tr class='passed'>
      <td>Gauge</td>
      <td>3</td>
    </tr>
    <tr class='failed'>
      <td>Mingle</td>
      <td>2</td>
    </tr>
    <tr class='skipped'>
      <td>foobar</td>
      <td>1</td>
    </tr>
  </tbody>
</table>
<span>Comment 1</span>
<span>Comment 2</span>
<span>Comment 3</span>`

var wSpecCommentsWithoutTableTag string = `<span></span>
<span>This is an executable specification file. This file follows markdown syntax.</span>
<span></span>
<span>To execute this specification, run</span>
<span>gauge specs</span>
<span></span>`

var wScenarioContainerStartPassDiv string = `<div class='scenario-container passed'>`
var wScenarioContainerStartFailDiv string = `<div class='scenario-container failed'>`
var wScenarioContainerStartSkipDiv string = `<div class='scenario-container skipped'>`

var wscenarioHeaderStartDiv string = `<div class="scenario-head">
  <h3 class="head borderBottom">Scenario Heading</h3>
  <span class="time">00:01:01</span>`

var wPassStepStartDiv string = `<div class='step'>
  <h5 class='execution-time'><span class='time'>Execution Time : 00:03:31</span></h5>
  <div class='step-info passed'>
    <ul>
      <li class='step'>
        <div class='step-txt'>`

var wFailStepStartDiv string = `<div class='step'>
  <h5 class='execution-time'><span class='time'>Execution Time : 00:03:31</span></h5>
  <div class='step-info failed'>
    <ul>
      <li class='step'>
        <div class='step-txt'>`

var wSkipStepStartDiv string = `<div class='step'>
  <div class='step-info skipped'>
    <ul>
      <li class='step'>
        <div class='step-txt'>`

var wStepEndDiv string = `<span>Say</span><span class='parameter'>"hi"</span><span>to</span><span class='parameter'>"gauge"</span>
          <div class='inline-table'>
            <div>
              <table>
                <tr>
                  <th>Word</th>
                  <th>Count</th>
                </tr>
                <tbody>
                  <tr>
                    <td>Gauge</td>
                    <td>3</td>
                  </tr>
                  <tr>
                    <td>Mingle</td>
                    <td>2</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </li>
    </ul>
  </div>
</div>
`

var re *regexp.Regexp = regexp.MustCompile("[ ]*[\n\t][ ]*")

var reportGenTests = []reportGenTest{
	{"generate body header with project name", bodyHeaderTag, &overview{ProjectName: "projname"}, wBodyHeader},
	{"generate report overview with tags", reportOverviewTag, &overview{"projname", "default", "foo", 34, "00:01:53", "Jun 3, 2016 at 12:29pm", 41, 2, 39, 0},
		wChartDiv + wResCntDiv + wEnvLi + wTagsLi + wSuccRateLi + wExecTimeLi + wTimestampLi},
	{"generate report overview without tags", reportOverviewTag, &overview{"projname", "default", "", 34, "00:01:53", "Jun 3, 2016 at 12:29pm", 41, 2, 39, 0},
		wChartDiv + wResCntDiv + wEnvLi + wSuccRateLi + wExecTimeLi + wTimestampLi},
	{"generate sidebar with appropriate pass/fail/skip class", sidebarDiv, &sidebar{
		IsPreHookFailure: false,
		Specs: []*specsMeta{
			newSpecsMeta("Passing Spec", "00:01:04", false, false, nil, "passing_spec.html"),
			newSpecsMeta("Failing Spec", "00:00:30", true, false, nil, "failing_spec.html"),
			newSpecsMeta("Skipped Spec", "00:00:00", false, true, nil, "skipped_spec.html"),
		}}, wSidebarAside},
	{"do not generate sidebar if presuitehook failure", sidebarDiv, &sidebar{
		IsPreHookFailure: true,
		Specs:            []*specsMeta{},
	}, ""},
	{"generate congratulations bar if all specs are passed", congratsDiv, &overview{}, wCongratsDiv},
	{"don't generate congratulations bar if some spec failed", congratsDiv, &overview{Failed: 1}, ""},
	{"generate hook failure div with screenshot", hookFailureDiv, newHookFailure("BeforeSuite", "SomeError", "iVBO", "Stack trace"), wHookFailureWithScreenhotDiv},
	{"generate hook failure div without screenshot", hookFailureDiv, newHookFailure("BeforeSuite", "SomeError", "", "Stack trace"), wHookFailureWithoutScreenhotDiv},
	{"generate spec header with tags", specHeaderStartTag, &specHeader{"Spec heading", "00:01:01", "/tmp/gauge/specs/foobar.spec", []string{"foo", "bar"}}, wSpecHeaderStartWithTags},
	{"generate div for tags", tagsDiv, &specHeader{Tags: []string{"tag1", "tag2"}}, wTagsDiv},
	{"generate spec comments with data table (if present)", specCommentsAndTableTag, newSpec(true), wSpecCommentsWithTableTag},
	{"generate spec comments without data table", specCommentsAndTableTag, newSpec(false), wSpecCommentsWithoutTableTag},
	{"generate passing scenario container", scenarioContainerStartDiv, &scenario{ExecStatus: pass}, wScenarioContainerStartPassDiv},
	{"generate failed scenario container", scenarioContainerStartDiv, &scenario{ExecStatus: fail}, wScenarioContainerStartFailDiv},
	{"generate skipped scenario container", scenarioContainerStartDiv, &scenario{ExecStatus: skip}, wScenarioContainerStartSkipDiv},
	{"generate scenario header", scenarioHeaderStartDiv, &scenario{Heading: "Scenario Heading", ExecTime: "00:01:01"}, wscenarioHeaderStartDiv},
	{"generate pass step start div", stepStartDiv, newStep(pass), wPassStepStartDiv},
	{"generate fail step start div", stepStartDiv, newStep(fail), wFailStepStartDiv},
	{"generate skipped step start div", stepStartDiv, newStep(skip), wSkipStepStartDiv},
}

func TestExecute(t *testing.T) {
	testReportGen(reportGenTests, t)
}

func testReportGen(reportGenTests []reportGenTest, t *testing.T) {
	buf := new(bytes.Buffer)
	for _, test := range reportGenTests {
		gen(test.tmpl, buf, test.input)

		got := removeNewline(buf.String())
		want := removeNewline(test.output)

		if got != want {
			t.Errorf("%s:\nwant:\n%q\ngot:\n%q\n", test.name, want, got)
		}
		buf.Reset()
	}
}

func removeNewline(s string) string {
	return re.ReplaceAllLiteralString(s, "")
}

func newHookFailure(name, errMsg, screenshot, stacktrace string) *hookFailure {
	return &hookFailure{
		HookName:   name,
		ErrMsg:     errMsg,
		Screenshot: screenshot,
		StackTrace: stacktrace,
	}
}

func newOverview() *overview {
	return &overview{
		ProjectName: "gauge-testsss",
		Env:         "default",
		SuccRate:    95,
		ExecTime:    "00:01:53",
		Timestamp:   "Jun 3, 2016 at 12:29pm",
	}
}

func newSpecsMeta(name, execTime string, failed, skipped bool, tags []string, fileName string) *specsMeta {
	return &specsMeta{
		SpecName:   name,
		ExecTime:   execTime,
		Failed:     failed,
		Skipped:    skipped,
		Tags:       tags,
		ReportFile: fileName,
	}
}

func newSpec(withTable bool) *spec {
	t := &table{
		Headers: []string{"Word", "Count"},
		Rows: []*row{
			{
				Cells: []string{"Gauge", "3"},
				Res:   pass,
			},
			{
				Cells: []string{"Mingle", "2"},
				Res:   fail,
			},
			{
				Cells: []string{"foobar", "1"},
				Res:   skip,
			},
		},
	}

	c1 := []string{"\n", "This is an executable specification file. This file follows markdown syntax.", "\n", "To execute this specification, run", "\tgauge specs", "\n"}
	c2 := []string{"Comment 1", "Comment 2", "Comment 3"}

	if withTable {
		return &spec{
			CommentsBeforeTable: c1,
			Table:               t,
			CommentsAfterTable:  c2,
		}
	}

	return &spec{
		CommentsBeforeTable: c1,
	}
}

func newStep(s status) *step {
	return &step{
		Fragments: []*fragment{
			{FragmentKind: textFragmentKind, Text: "Say "},
			{FragmentKind: staticFragmentKind, Text: "hi"},
			{FragmentKind: textFragmentKind, Text: " to "},
			{FragmentKind: dynamicFragmentKind, Text: "gauge"},
			{FragmentKind: tableFragmentKind,
				Table: &table{
					Headers: []string{"Word", "Count"},
					Rows: []*row{
						{Cells: []string{"Gauge", "3"}},
						{Cells: []string{"Mingle", "2"}},
					},
				},
			},
		},
		Res: &result{
			Status:   s,
			ExecTime: "00:03:31",
		},
	}
}
