/*
* This file contains all the functions and structs that are used to add modules to the framework.
 */
package govel

type initModuleFunc func(config map[interface{}]interface{}) error

// InitModules will add the modules to the framework.
func InitModules(m ...initModuleFunc) {
	modules = m
}
