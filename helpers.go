package govel

import (
	"bytes"
	"math"
	"net/http"
)

func getErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func searchRoute(l []routeNamed, e string) int {
	left := 0
	right := len(l) - 1

	for left <= right {

		middleIndex := int(math.Floor((float64(left+right) / 2)))
		middleElement := l[middleIndex].Route

		if middleElement == e {
			return middleIndex
		}

		if e < middleElement {
			right = middleIndex - 1
		} else {
			left = middleIndex + 1
		}

	}

	return -1
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	c := &Context{ResponseWriter: w, Request: r}

	// set the initial values
	c.Buf = new(bytes.Buffer)
	c.Headers = make(map[string]string)
	c.statusCode = 200
	c.SharedPayload = make(map[string]interface{})

	return c
}
