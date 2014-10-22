package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

func TestSetupRouter(t *testing.T) {
	input := map[string]string{
		"/":                                           "meta-data\nuser-data\n",
		"/user-data/":                                 "foo",
		"/meta-data/":                                 "ami-id\nblock-device-mapping/\nhostname",
		"/meta-data/ami-id/":                          "ami-12345678",
		"/meta-data/block-device-mapping/":            "ephemeral0\nephemeral1\nroot",
		"/meta-data/block-device-mapping/ephemeral0/": "sdb",
		"/meta-data/block-device-mapping/ephemeral1/": "sdg",
		"/meta-data/block-device-mapping/root/":       "/dev/sda1",
		"/meta-data/hostname/":                        "foo-bar",
	}

	router := setupRouter(input)

	ts := httptest.NewServer(router)
	defer ts.Close()

	for k, _ := range input {
		res, err := http.Get(ts.URL + k)
		if err != nil {
			t.Fatalf("Error on GET to '%s': %s", k, err.Error())
		}
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatalf("Error reading body from '%s': %s", k, err.Error())
		}

		if !reflect.DeepEqual(string(body), input[k]) {
			t.Fatalf("Error with body from '%s'.\nExpected\n'%s'\nGot\n'%s'", k, input[k], string(body))
		}
	}

}

func TestUrlData(t *testing.T) {
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
		"/":                                           "meta-data\nuser-data\n",
		"/user-data/":                                 "foo",
		"/meta-data/":                                 "ami-id\nblock-device-mapping/\nhostname",
		"/meta-data/ami-id/":                          "ami-12345678",
		"/meta-data/block-device-mapping/":            "ephemeral0\nephemeral1\nroot",
		"/meta-data/block-device-mapping/ephemeral0/": "sdb",
		"/meta-data/block-device-mapping/ephemeral1/": "sdg",
		"/meta-data/block-device-mapping/root/":       "/dev/sda1",
		"/meta-data/hostname/":                        "foo-bar",
	}

	results := urlData("/", input)

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
