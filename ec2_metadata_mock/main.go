package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

func rawData() (map[string]interface{}, error) {
	file := "./example_metadata.json"
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		msg := fmt.Sprintf("Unable to read '%s'. %s", file, err.Error())
		return map[string]interface{}{}, errors.New(msg)
	}

	var jsonData interface{}
	err = json.Unmarshal(raw, &jsonData)
	if err != nil {
		msg := fmt.Sprintf("Unable to parse JSON. %s", err.Error)
		return map[string]interface{}{}, errors.New(msg)
	}

	return jsonData.(map[string]interface{}), nil
}

func urlData(parent string, data map[string]interface{}) map[string]string {
	results := map[string]string{}

	for k, v := range data {
		val := reflect.Indirect(reflect.ValueOf(v))
		if val.Kind() == reflect.Interface {
			val = val.Elem()
		}

		switch val.Kind() {
		case reflect.String:
			results[parent+k+"/"] = val.String()

		case reflect.Map:
			mapKeysOutput := []string{}
			for subKey, subVal := range val.Interface().(map[string]interface{}) {
				subVal := reflect.Indirect(reflect.ValueOf(subVal))
				switch subVal.Kind() {
				case reflect.Map:
					mapKeysOutput = append(mapKeysOutput, subKey+"/")
				case reflect.String:
					mapKeysOutput = append(mapKeysOutput, subKey)
				}
			}

			sort.Strings(mapKeysOutput)
			results[parent+k+"/"] = strings.Join(mapKeysOutput, "\n")

			result := urlData(parent+k+"/", val.Interface().(map[string]interface{}))

			for resultKey, resultVal := range result {
				results[resultKey] = resultVal
			}
		}
	}

	return results
}

func main() {

	data, err := rawData()
	if err != nil {
		fmt.Printf("Unable to load raw data. Error: %s\n", err.Error())
		os.Exit(1)
	}

	uData := urlData("/", data)

	urls := []string{}

	for k, _ := range uData {
		urls = append(urls, k)
	}

	sort.Strings(urls)

	router := mux.NewRouter()

	for _, k := range urls {
		url := k
		data := uData[url]
		router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, data)
		}).Methods("GET")
	}

	http.ListenAndServe(":8080", router)
}
