package main

import (
	"reflect"
	"sort"
	"strings"
)

func parsedData(parent string, data map[string]interface{}) map[string]string {
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
			for subKey, valVal := range val.Interface().(map[string]interface{}) {
				valVal := reflect.Indirect(reflect.ValueOf(valVal))
				switch valVal.Kind() {
				case reflect.Map:
					mapKeysOutput = append(mapKeysOutput, subKey+"/")
				case reflect.String:
					mapKeysOutput = append(mapKeysOutput, subKey)
				}
			}

			sort.Strings(mapKeysOutput)
			results[parent+k+"/"] = strings.Join(mapKeysOutput, "\n")

			result := parsedData(parent+k+"/", val.Interface().(map[string]interface{}))

			for y, z := range result {
				results[y] = z
			}
		}
	}

	return results
}

func main() {
}
