/*
 * In this file are all the models and functions for testing.
 */
package govel

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Main struct for testing.
type Test struct {
	t t
}

type t interface {
	Error(args ...any)
}

type formForTesting struct {
	fields  url.Values
	request *http.Request
}

type multiPartFormDataForTesting struct {
	fields map[string]string
	files  map[string]string
	url    string
	method string
}

type testingValue struct {
	key   string
	value interface{}
}

// FromTestingKey returns the value inside the "key" inside the "testing" key of the .yaml file as a testingValue struct.
func (t *Test) FromTestingKey(key string) testingValue {
	value, err := configFileKeys["testing"].(map[interface{}]interface{})

	if err == false {
		t.t.Error("Key \"testing\" does not exist in .yaml file.")
	}

	return testingValue{value: value[key], key: key}
}

// String returns the value as a string. It panics if the value is not a string.
func (tv testingValue) String() string {
	string, err := tv.value.(string)

	if err == false {
		panic(fmt.Sprintf("Key %s is not a valid string", tv.key))
	}

	return string
}

// Int returns the value as an int. It panics if the value is not an int.
func (tv testingValue) Int() int {
	int, err := tv.value.(int)

	if err == false {
		panic(fmt.Sprintf("Key %s is not a valid int", tv.key))
	}

	return int
}

// Interface returns the value as an inteface{}.
func (tv testingValue) Interface() interface{} {
	return tv.value
}

/**
 * All functions to work with http requests.
 */

// UrlOk makes a GET HTTTP request to url and checks that its status code equals 200.
//
// If the status code is not 200 it will return false.
func (t *Test) UrlOk(url string) bool {
	response, err := http.Get(url)

	if err != nil {
		return false
	}

	if response.StatusCode != 200 {
		return false
	}

	return true
}

// TestRoute makes an HTTP GET request to the route URL (if it exists)
// and checks if the response status code is 200.
//
// The request will make it to the path http://127.0.0.1:port/path.
// So the server must be running locally.
func (t *Test) TestRoute(route string) bool {
	port := configFileKeys["port"]
	url := Route(route, nil)

	if url == "" {
		return false
	}

	full_url := fmt.Sprintf("http://127.0.0.1:%d%s", port, url)

	response, err := http.Get(full_url)

	if err != nil {
		return false
	}

	if response.StatusCode != 200 {
		return false
	}

	return true
}

/**
 * All functions to work with forms.
 */

// NewTestingForm returns a new formForTesting for making requests of type application/x-www-form-urlencoded.
func (t *Test) NewTestingForm(httpMethod string, to string) (*formForTesting, error) {
	req, err := http.NewRequest(httpMethod, to, nil)

	if err != nil {
		return &formForTesting{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return &formForTesting{
		request: req,
		fields:  url.Values{},
	}, nil
}

// AddField adds a new input for the form.
func (form *formForTesting) AddField(key string, value string) {
	form.fields.Set(key, value)
}

// Submit submits the form
func (form *formForTesting) Submit() (*http.Response, error) {
	form.request.Body = io.NopCloser(strings.NewReader(form.fields.Encode()))

	client := &http.Client{}

	resp, err := client.Do(form.request)

	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil
}

// NewMultiPartFormDataForm returns a new multiPartFormDataForTesting for making requests of type multipart/form-data.
func (t *Test) NewMultiPartFormData(httpMethod string, url string) multiPartFormDataForTesting {
	return multiPartFormDataForTesting{
		url:    url,
		method: httpMethod,
		files:  make(map[string]string),
		fields: make(map[string]string),
	}
}

// AddField adds a new input for the form.
func (form *multiPartFormDataForTesting) AddField(key string, value string) {
	form.fields[key] = value
}

// AddField adds a new file for the form.
func (form *multiPartFormDataForTesting) AddFile(fileName string, filePath string) {
	form.files[fileName] = filePath
}

// Submit submits the form.
func (form *multiPartFormDataForTesting) Submit() (*http.Response, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// iterate over the fields
	for key, value := range form.fields {
		w.WriteField(key, value)
	}

	// iterate over the files
	for key, value := range form.files {
		file, err := os.Open(value)

		if err != nil {
			return nil, err
		}

		fw, err := w.CreateFormFile(key, value)

		if err != nil {
			return nil, err
		}

		_, err = io.Copy(fw, file)

		if err != nil {
			return nil, err
		}

		file.Close()
	}

	w.Close()

	req, err := http.NewRequest(form.method, form.url, &buf)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
