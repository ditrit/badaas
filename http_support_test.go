package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const BaseUrl = "http://localhost:8000"

func (t *TestContext) requestGET(url string) error {
	response, err := t.httpClient.Get(fmt.Sprintf("%s%s", BaseUrl, url))
	if err != nil {
		return err
	}

	t.storeResponseInContext(response)
	return nil
}

func (t *TestContext) iSigninAsWithPassword(email, password string) error {
	url := fmt.Sprintf("%s%s", BaseUrl, "/login")
	payload := strings.NewReader(fmt.Sprintf("{\"email\": %q,\"password\": %q}", email, password))

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	res, err := t.httpClient.Do(req)
	if err != nil {
		return err
	}

	t.storeResponseInContext(res)
	return nil
}

func (t *TestContext) storeResponseInContext(response *http.Response) {
	t.statusCode = response.StatusCode

	buffer, err := io.ReadAll(response.Body)
	if err != nil {
		log.Panic(err)
	}
	response.Body.Close()
	json.Unmarshal(buffer, &t.json)
}

func (t *TestContext) assertStatusCode(_ context.Context, expectedStatusCode int) error {
	if t.statusCode != expectedStatusCode {
		return fmt.Errorf("expect status code %d but is %d", expectedStatusCode, t.statusCode)
	}
	return nil
}
func (t *TestContext) assertResponseFieldIsEquals(field string, expectedValue string) error {
	value := t.json[field].(string)
	if !assertValue(value, expectedValue) {
		return fmt.Errorf("expect response field %s is %s but is %s", field, expectedValue, value)
	}
	return nil
}

func assertValue(value string, expectedValue string) bool {
	return expectedValue == value
}

func (t *TestContext) assertResponseFieldMatchesRegex(field string, regex string) error {
	re, err := regexp.Compile(regex)
	if err != nil {
		return fmt.Errorf("regex did not compile (ERROR=%s)", err.Error())
	}
	value := t.json[field].(string)
	if !re.Match([]byte(value)) {
		return fmt.Errorf("%q do not match the regex %q", value, regex)
	}
	return nil
}
