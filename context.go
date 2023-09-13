package govel

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Textf Prints a formatted text string.
func (c *Context) Textf(statusCode int, format string, a ...interface{}) {
	c.statusCode = statusCode

	fmt.Fprintf(c.Buf, format, a...)
}

// Test prints a blank text.
func (c *Context) Text(statusCode int, s string) {
	c.statusCode = statusCode

	fmt.Fprint(c.Buf, s)
}

func (c *Context) Json(s int, j interface{}) error {
	c.ContentType("application/json")
	c.Status(s)

	if j != nil {
		return json.NewEncoder(c.Buf).Encode(j)
	}

	return nil
}

func (c *Context) RawJson(j interface{}) {
	c.ContentType("application/json")

	if j != nil {
		json.NewEncoder(c.Buf).Encode(j)
	}
}

// ContentType sets the Content-Type header.
func (c *Context) ContentType(content string) {
	c.Headers["Content-Type"] = content
}

// Bytes sends bytes to the response.
func (c *Context) Bytes(s int, b []byte) {
	c.statusCode = s

	c.Buf.Write(b)
}

// Status sets the status code.
func (c *Context) Status(s int) *Context {
	c.statusCode = s

	return c
}

// SetFormValues sets the values of the struct from a map by the "form" tag.
func (c *Context) SetFormValues(ptr interface{}, formValues map[string]interface{}) {
	value := reflect.ValueOf(ptr)

	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		panic("ptr is not a pointer to a struct")
	}

	value = value.Elem()

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)

		tag := field.Tag.Get("form")

		if tag != "" {
			if fieldValue, ok := formValues[tag]; ok {

				fieldValue := reflect.ValueOf(fieldValue)

				if !fieldValue.Type().AssignableTo(field.Type) {
					panic(fmt.Sprintf("value for field '%s' is not assignable to its type", field.Name))

				}

				value.Field(i).Set(fieldValue.Convert(field.Type))
			}
		}
	}

}

// RemoteAddr returns the address of the client.
func (c *Context) RemoteAddr() (ip string, port int) {
	remote := c.Request.RemoteAddr

	addr := strings.Split(remote, ":")

	ip = addr[0]
	port, _ = strconv.Atoi(addr[1])

	return ip, port
}

// Params returns all route parameters.
func (c *Context) Params() map[string]string {
	return mux.Vars(c.Request)
}

// Param returns a route parameter by its name.
func (c *Context) Param(key string) string {
	return c.Params()[key]
}

// Query returns a query param by its name.
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// SetHeader sets a header by its name.
func (c *Context) SetHeader(key string, value string) {
	c.Request.Header.Set(key, value)
}

// GetHeader returns a header by its name.
func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

// Session always returns a session, even if empty.
//
// The session will not be saved if no changes are made.
func (c *Context) Session(session string) (Session, error) {
	s, err := Store.Get(c.Request, session)

	originalState := make(map[interface{}]interface{})

	for key, value := range s.Values {
		originalState[key] = value
	}

	sessionStruct := Session{session: s, originalState: originalState}

	c.sessions = append(c.sessions, sessionStruct)

	return sessionStruct, err
}

// NewForm creates a new "Form" model for HTTP requests with "application/x-www-form-urlencoded".
func (c *Context) NewForm() (Form, error) {
	w := c.ResponseWriter
	r := c.Request

	// parse the form

	err := r.ParseForm()

	if err != nil {
		return Form{writer: w, request: r}, err
	}

	// return the form
	return Form{writer: w, request: r}, nil
}

// Creates a new "MultipartFormData" model for HTTP requests with "multipart/form-data".
//
// For a detailed explanation of the value of the "maxMemory" variable see https://pkg.go.dev/net/http#Request.ParseMultipartForm.
func (c *Context) NewMultiPartFormDataForm(maxMemory int64) (MultipartFormData, error) {
	w := c.ResponseWriter
	r := c.Request

	// parse the form

	err := r.ParseMultipartForm(maxMemory)

	if err != nil {
		return MultipartFormData{writer: w, request: r}, err
	}

	// return the form
	return MultipartFormData{writer: w, request: r}, nil
}

// Redirect redirects the user to a specific URL.
func (c *Context) Redirect(url string, code int) {
	c.statusCode = code
	c.Headers["Location"] = url
}

// RedirectToRoute redirects the user to a specific route.
//
// Route is called internally to get the route.
func (c *Context) RedirectToRoute(name string, params SMap, code int) {
	route := Route(name, params)

	c.Redirect(route, code)
}
