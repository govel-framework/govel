package govel

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	yaml "gopkg.in/yaml.v2"
)

var (
	router = mux.NewRouter().StrictSlash(true)

	savedRoutes = make(map[string]*routeModel)
	routes      []routeNamed

	// A global "Panic handler" function.
	panicHandlerFunc panicHandler

	// modules to be initialized
	modules []initModuleFunc

	// Store is the session store.
	Store *sessions.CookieStore

	// Represents the configuracion of the .yaml file in a map.
	configFileKeys map[interface{}]interface{}

	// Global middlewares
	globalMiddlewares middlewaresFunctions

	// indicates if the current route is inside a group
	inGroup            bool
	currentGroupConfig = &groupModel{routes: make(map[string]*routeModel)}
)

// LoadConfigFileForTests returns a test model for testing.
func LoadConfigFileForTests(routes_func func(), t t, configFilePath string) Test {
	if routes_func != nil {
		routes_func()
	}

	readYamlAndGetPort(configFilePath)

	return Test{
		t: t,
	}

}

// LoadConfigFile reads the .yaml file and starts the web server based on its configuration.
func LoadConfigFIle(file string) {
	// load every route
	for _, m := range savedRoutes {

		router.
			Path(m.path).
			Handler(callFunction(m.action, m.middlewares)).
			Methods(m.method)

		newRoute := routeNamed{Route: m.name, Url: m.path}

		routes = append(routes, newRoute)
	}

	// main funcs
	port := fmt.Sprintf("%s%s", ":", readYamlAndGetPort(file))

	log.Print("Server is running on port " + port)

	// start the server

	http.Handle("/", router)

	getErr(http.ListenAndServe(port, nil))
}

// readYamlAndGetPort gets the port from the yaml file, parses the file and sets the rest of the configuration.
func readYamlAndGetPort(file string) string {
	// structs
	yamlConfig := configYamlFile{}

	// read .yaml file and return port

	fileContent, err := os.ReadFile(file)

	getErr(err)

	err = yaml.Unmarshal(fileContent, &yamlConfig)

	getErr(err)

	// check the required fields
	if yamlConfig.Port == 0 {
		panic("The port is required.")
	}

	configFileKeys = make(map[interface{}]interface{})

	yaml.Unmarshal(fileContent, configFileKeys)

	// set the rest of the configuraiton
	if yamlConfig.Static.Dir != "" && yamlConfig.Static.Path != "" {
		s := http.StripPrefix(yamlConfig.Static.Path, http.FileServer(http.Dir(yamlConfig.Static.Dir+"/")))

		router.PathPrefix(yamlConfig.Static.Path).Handler(s)
	}

	if yamlConfig.Keys.Sessions != "" {
		Store = sessions.NewCookieStore([]byte(yamlConfig.Keys.Sessions))
	}

	// initilize the modules

	for _, module := range modules {
		// load the module
		err := module(configFileKeys)

		if err != nil {
			panic(fmt.Sprintf("Cannot initialize module: %s", err.Error()))
		}
	}

	// "clean" the vars
	currentGroupConfig = nil
	inGroup = false
	modules = nil
	savedRoutes = nil

	// return the port
	return strconv.Itoa(yamlConfig.Port)
}

// functions

func Get(path string, action routeFunction) *routeModel {

	m := routeModel{unique_id: time.Now().String(), path: path, action: action, method: "GET"}
	m.update()

	return &m
}

func Post(path string, action routeFunction) *routeModel {

	m := routeModel{unique_id: time.Now().String(), path: path, action: action, method: "POST"}
	m.update()

	return &m
}

func Put(path string, action routeFunction) *routeModel {

	m := routeModel{unique_id: time.Now().String(), path: path, action: action, method: "PUT"}
	m.update()

	return &m
}

func Delete(path string, action routeFunction) *routeModel {

	m := routeModel{unique_id: time.Now().String(), path: path, action: action, method: "DELETE"}
	m.update()

	return &m
}

// Middlewares adds one or multiple middlewares to a route.
func (m *routeModel) Middlewares(action ...middlewareFunction) *routeModel {
	m.middlewares = action

	m.update()

	return m
}

// Name adds a name to a route.
func (m *routeModel) Name(name string) *routeModel {
	m.name = name

	m.update()

	return m
}

// Saves or updates the route.
func (m *routeModel) update() {

	if inGroup {

		if !m.pathUpdated {
			m.path = currentGroupConfig.prefix + m.path
			m.pathUpdated = true
		}

		currentGroupConfig.routes[m.unique_id] = m
	}

	savedRoutes[m.unique_id] = m
}

// Gets route's url by its name.
func Route(r string, data SMap) string {
	if r == "" {
		return ""
	}

	// sort routes
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Route < routes[j].Route
	})

	index := searchRoute(routes, r)

	if index == -1 {
		return ""
	} else {
		url := routes[index].Url

		// from here we check if the route needs any parameters
		re := regexp.MustCompile(`(?U)\{.*\}`)

		matches := re.FindAllString(url, -1)

		if len(matches) == 0 {
			return url
		} else {
			if data == nil {
				panic(fmt.Sprintf("Route %s requires parameters. data cannot be nil", r))
			}
		}

		re = regexp.MustCompile(`\{|\}`)

		for _, i := range matches {
			key := re.ReplaceAllString(i, "")

			text, exists := data[key]

			if !exists {
				panic(fmt.Sprintf("Route %s requires %s", r, key))
			}

			url = strings.Replace(url, "{"+key+"}", text, -1)
		}

		return url
	}
}

// GetKeyFromYAML returns the value of the key from the YAML config file.
//
// If the key is empty, it returns the whole YAML config file as a map.
func GetKeyFromYAML(key string) interface{} {
	if key == "" {
		return configFileKeys
	}

	return configFileKeys[key]
}

// Sets a "404 url not found" function.
func Set404NotFound(function routeFunction) {
	router.NotFoundHandler = httpHandler(function)
}

// Sets a general "method not allowed" function.
func SetMethodNotAllowed(function routeFunction) {
	router.MethodNotAllowedHandler = httpHandler(function)
}

// Sets a general "handle panic" function.
func SetPanicHandler(function panicHandler) {
	panicHandlerFunc = function
}

// Sets global middlewares.
func SetGlobalMiddlewares(function ...middlewareFunction) {
	globalMiddlewares = function
}
