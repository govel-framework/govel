package govel

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Here are global variables or variables that will be used outside the function where they were declared

var (
	// Router
	router = mux.NewRouter().StrictSlash(true)

	// Routes
	savedRoutes = make(map[string]*routeModel)
	routes      []routeNamed

	// A global "Panic handler" function
	panicHandlerFunc panicHandler

	// Global middlewares
	globalMiddlewares middlewaresFunctions

	// indicates if the current route is inside a group
	inGroup            bool
	currentGroupConfig = &groupModel{routes: make(map[string]*routeModel)}

	// Represents the configuracion of the .yaml file in a map.
	configFileKeys map[interface{}]interface{}

	// modules to be initialized
	modules []initModuleFunc

	// Store is the session store.
	Store *sessions.CookieStore
)
