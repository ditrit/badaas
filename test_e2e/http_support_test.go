package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v16"
	integrationtests "github.com/ditrit/badaas/test_integration"
	"github.com/elliotchance/pie/v2"
)

const BaseUrl = "http://localhost:8000"

func (t *TestContext) requestGet(url string) error {
	return t.request(url, http.MethodGet, nil, nil)
}

func (t *TestContext) requestWithJson(url, method string, jsonTable *godog.Table) error {
	return t.request(url, method, nil, jsonTable)
}

func (t *TestContext) request(url, method string, query map[string]string, jsonTable *godog.Table) error {
	var payload io.Reader
	var err error
	if jsonTable != nil {
		payload, err = buildJSONFromTable(jsonTable)
		if err != nil {
			return err
		}
	}

	method, err = checkMethod(method)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(method, BaseUrl+url, payload)
	if err != nil {
		return fmt.Errorf("failed to build request ERROR=%s", err.Error())
	}

	q := request.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	request.URL.RawQuery = q.Encode()

	response, err := t.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to run request ERROR=%s", err.Error())
	}
	t.storeResponseInContext(response)
	return nil
}

func (t *TestContext) storeResponseInContext(response *http.Response) {
	t.statusCode = response.StatusCode

	err := json.NewDecoder(response.Body).Decode(&t.json)
	if err != nil {
		t.json = map[string]any{}
	}
}

func (t *TestContext) assertStatusCode(expectedStatusCode int) error {
	if t.statusCode != expectedStatusCode {
		return fmt.Errorf("expect status code %d but is %d", expectedStatusCode, t.statusCode)
	}
	return nil
}

func (t *TestContext) assertResponseFieldIsEquals(field string, expectedValue string) error {
	fields := strings.Split(field, ".")
	jsonMap := t.json.(map[string]any)

	for _, field := range fields[:len(fields)-1] {
		intValue, present := jsonMap[field]
		if !present {
			return fmt.Errorf("expected response field %s to be %s but it is not present", field, expectedValue)
		}
		jsonMap = intValue.(map[string]any)
	}

	lastValue, present := jsonMap[pie.Last(fields)]
	if !present {
		return fmt.Errorf("expected response field %s to be %s but it is not present", field, expectedValue)
	}

	if !assertValue(lastValue, expectedValue) {
		return fmt.Errorf("expected response field %s to be %s but is %v", field, expectedValue, lastValue)
	}

	return nil
}

func assertValue(value any, expectedValue string) bool {
	switch value.(type) {
	case string:
		return expectedValue == value
	case int:
		expectedValueInt, err := strconv.Atoi(expectedValue)
		if err != nil {
			panic(err)
		}
		return expectedValueInt == value
	case float64:
		expectedValueFloat, err := strconv.ParseFloat(expectedValue, 64)
		if err != nil {
			panic(err)
		}
		return expectedValueFloat == value
	default:
		panic("unsupported format")
	}
}

// build a map from a godog.Table
func buildMapFromTable(table *godog.Table) (map[string]any, error) {
	data := make(map[string]any, 0)
	for indexRow, row := range table.Rows {
		if indexRow == 0 {
			for indexCell, cell := range row.Cells {
				if cell.Value != []string{"key", "value", "type"}[indexCell] {
					return nil, fmt.Errorf("should have %q as first line of the table", "| key | value | type |")
				}
			}
		} else {
			key := row.Cells[0].Value
			valueAsString := row.Cells[1].Value
			valueType := row.Cells[2].Value

			switch valueType {
			case stringValueType:
				data[key] = valueAsString
			case booleanValueType:
				boolean, err := strconv.ParseBool(valueAsString)
				if err != nil {
					return nil, fmt.Errorf("can't parse %q as boolean for key %q", valueAsString, key)
				}
				data[key] = boolean
			case integerValueType:
				integer, err := strconv.ParseInt(valueAsString, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("can't parse %q as integer for key %q", valueAsString, key)
				}
				data[key] = integer
			case floatValueType:
				floatingNumber, err := strconv.ParseFloat(valueAsString, 64)
				if err != nil {
					return nil, fmt.Errorf("can't parse %q as float for key %q", valueAsString, key)
				}
				data[key] = floatingNumber
			case jsonValueType:
				jsonMap := map[string]string{}
				err := json.Unmarshal([]byte(valueAsString), &jsonMap)
				if err != nil {
					return nil, fmt.Errorf("can't parse %q as json for key %q", valueAsString, key)
				}
				data[key] = jsonMap
			case nullValueType:
				data[key] = nil
			default:
				return nil, fmt.Errorf("type %q does not exists, please use %v", valueType, []string{stringValueType, booleanValueType, integerValueType, floatValueType, nullValueType})
			}

		}
	}

	return data, nil
}

// build a json payload in the form of a reader from a godog.Table
func buildJSONFromTable(table *godog.Table) (io.Reader, error) {
	data, err := buildMapFromTable(table)
	if err != nil {
		panic("should not return an error")
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		panic("should not return an error")
	}
	return strings.NewReader(string(bytes)), nil
}

const (
	stringValueType  = "string"
	booleanValueType = "boolean"
	integerValueType = "integer"
	floatValueType   = "float"
	nullValueType    = "null"
	jsonValueType    = "json"
)

// check if the method is allowed and sanitize the string
func checkMethod(method string) (string, error) {
	allowedMethods := []string{http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace}
	sanitizedMethod := strings.TrimSpace(strings.ToUpper(method))
	if !pie.Contains(
		allowedMethods,
		sanitizedMethod,
	) {
		return "", fmt.Errorf("%q is not a valid HTTP method (please choose between %v)", method, allowedMethods)
	}

	return sanitizedMethod, nil
}

func (t *TestContext) objectExists(entityType string, jsonTable *godog.Table) error {
	err := t.request(
		"/eav/objects/"+entityType,
		http.MethodPost,
		nil,
		jsonTable,
	)
	if err != nil {
		return err
	}

	err = t.assertStatusCode(http.StatusCreated)
	if err != nil {
		return err
	}

	return nil
}

func (t *TestContext) objectExistsWithRelation(entityType string, relationAttribute string, jsonTable *godog.Table) error {
	id, present := t.json.(map[string]any)["id"]
	if !present {
		panic("object id not available")
	}

	jsonTable.Rows = append(jsonTable.Rows, &messages.PickleTableRow{
		Cells: []*messages.PickleTableCell{
			{
				Value: relationAttribute,
			},
			{
				Value: id.(string),
			},
			{
				Value: stringValueType,
			},
		},
	})
	return t.objectExists(entityType, jsonTable)
}

func (t *TestContext) queryWithObjectID(entityType string) error {
	id, present := t.json.(map[string]any)["id"]
	if !present {
		panic("object id not available")
	}

	err := t.requestGet(
		"/eav/objects/" + entityType + "/" + id.(string),
	)
	if err != nil {
		return err
	}

	err = t.assertStatusCode(http.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

func (t *TestContext) queryObjectsWithConditions(entityType string, jsonTable *godog.Table) error {
	err := t.requestWithJson(
		"/eav/objects/"+entityType,
		http.MethodGet,
		jsonTable,
	)
	if err != nil {
		return err
	}

	err = t.assertStatusCode(http.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

func (t *TestContext) queryAllObjects(entityType string) error {
	err := t.requestGet(
		"/eav/objects/" + entityType,
	)
	if err != nil {
		return err
	}

	err = t.assertStatusCode(http.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

func (t *TestContext) thereAreObjects(expectedAmount int, entityType string) error {
	amount := len(t.json.([]any))
	if amount != expectedAmount {
		return fmt.Errorf("expect amount %d, but there are %d objects of type %s", expectedAmount, amount, entityType)
	}

	return nil
}

func (t *TestContext) thereIsObjectWithAttributes(expectedEntityType string, jsonTable *godog.Table) error {
	objectList := t.json.([]any)
	expectedValues, err := buildMapFromTable(jsonTable)
	if err != nil {
		log.Fatalln(err)
	}

	for _, object := range objectList {
		objectMap := object.(map[string]any)
		objectAttrs := objectMap["attrs"].(map[string]any)

		if objectMap["type"] == expectedEntityType {
			if t.areAllAttributesEqual(objectAttrs, expectedValues) {
				return nil
			}
		}
	}

	return fmt.Errorf("object with attributes %v not found in %v", expectedValues, objectList)
}

func (t *TestContext) deleteWithObjectID(entityType string) error {
	id, present := t.json.(map[string]any)["id"]
	if !present {
		panic("object id not available")
	}

	err := t.request(
		"/eav/objects/"+entityType+"/"+id.(string),
		http.MethodDelete,
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	err = t.assertStatusCode(http.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

func (t *TestContext) modifyWithAttributes(entityType string, jsonTable *godog.Table) error {
	id, present := t.json.(map[string]any)["id"]
	if !present {
		panic("object id not available")
	}

	err := t.request(
		"/eav/objects/"+entityType+"/"+id.(string),
		http.MethodPut,
		nil,
		jsonTable,
	)
	if err != nil {
		return err
	}

	err = t.assertStatusCode(http.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

func (t *TestContext) saleExists(productInt int, code int, description string) {
	product := &integrationtests.Product{
		Int: productInt,
	}

	sale := &integrationtests.Sale{
		Code:        code,
		Description: description,
		Product:     *product,
	}
	err := t.db.Create(sale).Error
	if err != nil {
		log.Fatalln(err)
	}
}

func (t *TestContext) querySalesWithConditions(jsonTable *godog.Table) error {
	err := t.requestWithJson(
		"/objects/sale",
		http.MethodGet,
		jsonTable,
	)
	if err != nil {
		return err
	}

	err = t.assertStatusCode(http.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

func (t *TestContext) thereIsSaleWithAttributes(jsonTable *godog.Table) error {
	objectList := t.json.([]any)
	expectedValues, err := buildMapFromTable(jsonTable)
	if err != nil {
		log.Fatalln(err)
	}

	for _, object := range objectList {
		if t.areAllAttributesEqual(object.(map[string]any), expectedValues) {
			return nil
		}
	}

	return fmt.Errorf("object with attributes %v not found in %v", expectedValues, objectList)
}

func (t *TestContext) areAllAttributesEqual(objectMap, expectedValues map[string]any) bool {
	allEqual := true
	for attributeName, expectedValue := range expectedValues {
		actualValue, isPresent := objectMap[attributeName]
		if !isPresent || actualValue != expectedValue {
			allEqual = false
		}
	}

	return allEqual
}
