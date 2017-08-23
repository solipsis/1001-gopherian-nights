// Simple command line tool to find out what tests failed
// on an Atlassian Bamboo instance
package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var username = os.Getenv("BAMBOO_USERNAME")
var password = os.Getenv("BAMBOO_PASSWORD")

func main() {

	environments := fetchEnvironments()

	// Get the failing tests
	for _, env := range environments {
		arr := strings.Split(env, "-")
		projectKey, planKey, jobKey := arr[0], arr[1], arr[2]
		plan := fmt.Sprintf("https://bamboo-auto.compyanyName.com/rest/api/latest/result/%s-%s-JOBTITLE-%s?expand=testResults.failedTests", projectKey, planKey, jobKey)
		// print tests that are failing on bamboo
		fmt.Println(findFailingTests(plan))
	}

}

type result struct {
	PlanName    string `xml:"planName"`
	FailedTests []struct {
		ClassName  string `xml:"className,attr"`
		MethodName string `xml:"methodName,attr"`
	} `xml:"testResults>failedTests>testResult"`
}

type env struct {
	Environments []struct {
		Key string `xml:"key,attr"`
	} `xml:"results>result"`
}

func (r *result) String() string {
	var b bytes.Buffer
	for _, v := range r.FailedTests {
		b.WriteString(fmt.Sprintf("\tTestClass: %s\n", v.ClassName))
		b.WriteString(fmt.Sprintf("\t\tTest: %s\n", v.MethodName))
	}
	return b.String()
}

func fetchEnvironments() []string {

	req, err := http.NewRequest("GET", "https://bamboo-auto.companyName.com/rest/api/latest/result/", nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.SetBasicAuth(username, password)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error parsing body: %v", err)
	}

	r := new(env)
	err = xml.Unmarshal(body, r)
	if err != nil {
		log.Fatalf("Error unmarshaling: %v", err)
	}
	arr := []string{}
	for _, s := range r.Environments {
		arr = append(arr, s.Key)
	}
	return arr
}

func findFailingTests(url string) *result {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.SetBasicAuth(username, password)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error parsing body: %v", err)
	}

	r := new(result)
	err = xml.Unmarshal(body, r)
	if err != nil {
		log.Fatalf("Error unmarshaling: %v", err)
	}
	return r
}
