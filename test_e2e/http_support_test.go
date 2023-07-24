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
	"github.com/elliotchance/pie/v2"

	"github.com/ditrit/badaas/testintegration/models"
)

const BaseURL = "http://localhost:8000"

func (t *TestContext) requestGet(url string) error {
	return t.request(url, http.MethodGet, nil)
}

func (t *TestContext) requestWithJSON(url, method string, jsonTable *godog.Table) error {
	return t.request(url, method, jsonTable)
}

func (t *TestContext) request(url, method string, jsonTable *godog.Table) error {
	var payload io.Reader

	if jsonTable != nil {
		payload = buildJSONFromTable(jsonTable)
	}

	method, err := checkMethod(method)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(method, BaseURL+url, payload)
	if err != nil {
		return fmt.Errorf("failed to build request ERROR=%s", err.Error())
	}

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

	err = response.Body.Close()
	if err != nil {
		log.Fatalln(err)
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

	jsonMap, ok := t.json.(map[string]any)
	if !ok {
		log.Fatalln("json is not a map")
	}

	for _, field := range fields[:len(fields)-1] {
		intValue, present := jsonMap[field]
		if !present {
			return fmt.Errorf("expected response field %s to be %s but it is not present", field, expectedValue)
		}

		intValueMap, ok := intValue.(map[string]any)
		if !ok {
			log.Fatalln("intValue is not a map")
		}

		jsonMap = intValueMap
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

func verifyHeader(row *messages.PickleTableRow) error {
	for indexCell, cell := range row.Cells {
		if cell.Value != []string{"key", "value", "type"}[indexCell] {
			return fmt.Errorf("should have %q as first line of the table", "| key | value | type |")
		}
	}

	return nil
}

func getTableValue(key, valueAsString, valueType string) (any, error) {
	switch valueType {
	case stringValueType:
		return valueAsString, nil
	case booleanValueType:
		boolean, err := strconv.ParseBool(valueAsString)
		if err != nil {
			return nil, fmt.Errorf("can't parse %q as boolean for key %q", valueAsString, key)
		}

		return boolean, nil
	case integerValueType:
		integer, err := strconv.ParseInt(valueAsString, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("can't parse %q as integer for key %q", valueAsString, key)
		}

		return integer, nil
	case floatValueType:
		floatingNumber, err := strconv.ParseFloat(valueAsString, 64)
		if err != nil {
			return nil, fmt.Errorf("can't parse %q as float for key %q", valueAsString, key)
		}

		return floatingNumber, nil
	case jsonValueType:
		jsonMap := map[string]string{}

		err := json.Unmarshal([]byte(valueAsString), &jsonMap)
		if err != nil {
			return nil, fmt.Errorf("can't parse %q as json for key %q", valueAsString, key)
		}

		return jsonMap, nil
	default:
		return nil, fmt.Errorf(
			"type %q does not exists, please use %v",
			valueType,
			[]string{stringValueType, booleanValueType, integerValueType, floatValueType},
		)
	}
}

// build a map from a godog.Table
func buildMapFromTable(table *godog.Table) (map[string]any, error) {
	data := make(map[string]any, 0)

	err := verifyHeader(table.Rows[0])
	if err != nil {
		return nil, err
	}

	for _, row := range table.Rows[1:] {
		key := row.Cells[0].Value
		valueAsString := row.Cells[1].Value
		valueType := row.Cells[2].Value

		value, err := getTableValue(key, valueAsString, valueType)
		if err != nil {
			return nil, err
		}

		data[key] = value
	}

	return data, nil
}

// build a json payload in the form of a reader from a godog.Table
func buildJSONFromTable(table *godog.Table) io.Reader {
	data, err := buildMapFromTable(table)
	if err != nil {
		panic("should not return an error")
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		panic("should not return an error")
	}

	return strings.NewReader(string(bytes))
}

const (
	stringValueType  = "string"
	booleanValueType = "boolean"
	integerValueType = "integer"
	floatValueType   = "float"
	jsonValueType    = "json"
)

// check if the method is allowed and sanitize the string
func checkMethod(method string) (string, error) {
	allowedMethods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

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
	jsonTable.Rows = append(jsonTable.Rows, &messages.PickleTableRow{
		Cells: []*messages.PickleTableCell{
			{
				Value: relationAttribute,
			},
			{
				Value: t.getIDFromJSON(),
			},
			{
				Value: stringValueType,
			},
		},
	})

	return t.objectExists(entityType, jsonTable)
}

func (t *TestContext) getIDFromJSON() string {
	id, present := t.json.(map[string]any)["id"]
	if !present {
		log.Fatalln("object id not available")
	}

	idString, ok := id.(string)
	if !ok {
		log.Fatalln("id in json is not a string")
	}

	return idString
}

func (t *TestContext) saleExists(productInt int, code int, description string) {
	product := &models.Product{
		Int: productInt,
	}

	sale := &models.Sale{
		Code:        code,
		Description: description,
		Product:     *product,
	}

	if err := t.db.Create(sale).Error; err != nil {
		log.Fatalln(err)
	}
}

func (t *TestContext) querySalesWithConditions(jsonTable *godog.Table) error {
	err := t.requestWithJSON(
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
	expectedValues, err := buildMapFromTable(jsonTable)
	if err != nil {
		log.Fatalln(err)
	}

	objectMapList := t.getObjectMapListFromJSON()
	for _, objectMap := range objectMapList {
		if t.areAllAttributesEqual(objectMap, expectedValues) {
			return nil
		}
	}

	return fmt.Errorf("object with attributes %v not found in %v", expectedValues, objectMapList)
}

func (t *TestContext) getListFromJSON() []any {
	objectList, ok := t.json.([]any)
	if !ok {
		log.Fatalln("json is not a list")
	}

	return objectList
}

func (t *TestContext) getObjectMapListFromJSON() []map[string]any {
	objectList := t.getListFromJSON()

	return pie.Map(objectList, func(object any) map[string]any {
		objectMap, ok := object.(map[string]any)
		if !ok {
			log.Fatalln("object in json list is not a map")
		}

		return objectMap
	})
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
