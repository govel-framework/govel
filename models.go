package govel

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/sessions"
)

type panicHandler func(*Context, interface{})

type Context struct {
	ResponseWriter http.ResponseWriter

	Request *http.Request

	// SharedPayload shares data between middlewares and routes.
	//
	// SharedPayload gets cleaned up after each request.
	SharedPayload map[string]interface{}

	// Buf is a buffer that is used to store the response body.
	Buf *bytes.Buffer

	// Headers stores the response headers.
	Headers map[string]string

	statusCode int

	sessions []Session
}

/*
* Router and routes structs and functions
 */

type routeModel struct {
	unique_id   string
	path        string
	action      routeFunction
	middlewares middlewaresFunctions
	name        string
	method      string
	pathUpdated bool
}

type routeFunction func(c *Context)

type middlewaresFunctions []middlewareFunction

type middlewareFunction func(c *Context) int

type routeNamed struct {
	Route string
	Url   string
}

type groupModel struct {
	parent      *groupModel
	prefix      string
	routes      map[string]*routeModel
	middlewares middlewaresFunctions
	name        string
	subGroups   []*groupModel
}

/*
* Yaml file struct
 */

type sqlStruct struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Driver   string `yaml:"driver"`
}

type staticStruct struct {
	Path string `yaml:"path"`
	Dir  string `yaml:"dir"`
}

type keysStruct struct {
	Sessions string `yaml:"sessions"`
}

type configYamlFile struct {
	Port   int          `yaml:"port"`
	Sql    sqlStruct    `yaml:"sql"`
	Static staticStruct `yaml:"static"`
	Keys   keysStruct   `yaml:"keys"`
}

/*
* Forms and validations
 */

type Form struct {
	writer  http.ResponseWriter
	request *http.Request

	skip          bool
	validationVar interface{}
}

type MultipartFormData struct {
	writer  http.ResponseWriter
	request *http.Request

	skip          bool
	validationVar interface{}
}

type formInterface interface {
	skipValidation() bool

	setSkipValidation(bool)

	setVariableForValidation(interface{})

	getVariableForValidation() interface{}

	Get(string) string
}

type FormFile struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
}

// Map is an alias for map[string]interface{}.
type Map map[string]interface{}

// SMap is an alias for map[string]string.
type SMap map[string]string

// OnError is an alias for map[string]SMap.
type OnError map[string]SMap

// validateFunc is the template for a function used during form validation.
//
// form is only used when a function needs direct access to the form.
type validateFunc func(key string, value string, form formInterface) error

type validationError struct {
	Format string
	Key    string
}

func (e *validationError) Error() string {
	return fmt.Sprintf(e.Format, e.Key)
}

/*
 * Sessions
 */

type Session struct {
	session *sessions.Session

	originalState map[interface{}]interface{}
}
