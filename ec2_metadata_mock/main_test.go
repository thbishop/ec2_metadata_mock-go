package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestParsedData(t *testing.T) {
	input := map[string]interface{}{
		"user-data": "foo",
		"meta-data": map[string]interface{}{
			"ami-id": "ami-12345678",
			"block-device-mapping": map[string]interface{}{
				"root":       "/dev/sda1",
				"ephemeral0": "sdb",
				"ephemeral1": "sdg",
			},
			"hostname": "foo-bar",
		},
	}

	expected := map[string]string{
		"/user-data/":                                 "foo",
		"/meta-data/":                                 "ami-id\nblock-device-mapping/\nhostname",
		"/meta-data/ami-id/":                          "ami-12345678",
		"/meta-data/block-device-mapping/":            "ephemeral0\nephemeral1\nroot",
		"/meta-data/block-device-mapping/ephemeral0/": "sdb",
		"/meta-data/block-device-mapping/ephemeral1/": "sdg",
		"/meta-data/block-device-mapping/root/":       "/dev/sda1",
		"/meta-data/hostname/":                        "foo-bar",
	}

	results := parsedData("/", input)

	var expectedKeys, resultsKeys []string

	for k, _ := range results {
		resultsKeys = append(resultsKeys, k)
	}

	for k, _ := range expected {
		expectedKeys = append(expectedKeys, k)
	}

	sort.Strings(resultsKeys)
	sort.Strings(expectedKeys)

	if !reflect.DeepEqual(resultsKeys, expectedKeys) {
		t.Fatalf("Expected '%v' got '%v'", expectedKeys, resultsKeys)
	}

	for k, _ := range results {
		if !reflect.DeepEqual(results[k], expected[k]) {
			t.Fatalf("Values at key '%s' do not match. Expected '%s' got '%s'", k, expected[k], results[k])
		}
	}
}
