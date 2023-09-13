package govel

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

// callFunction is the function between the request and the action.
func callFunction(action routeFunction, middlewares middlewaresFunctions) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c := newContext(rw, r)

		defer func() {
			// finish the request

			// save the sessions if any
			if len(c.sessions) > 0 {
				for _, session := range c.sessions {
					// check if the sessions has new values
					if reflect.DeepEqual(session.originalState, session.session.Values) {
						continue
					}

					session.session.Save(c.Request, c.ResponseWriter)
				}
			}

			// set the headers
			for key, value := range c.Headers {
				c.ResponseWriter.Header().Set(key, value)
			}

			// write the rest of the response
			c.ResponseWriter.WriteHeader(c.statusCode)
			c.ResponseWriter.Write(c.Buf.Bytes())

			// recover from panic
			if r := recover(); r != nil {
				if panicHandlerFunc != nil {
					panicHandlerFunc(c, r)
				} else {

					// get the stack trace of the panic
					actionFunctionName := runtime.FuncForPC(reflect.ValueOf(action).Pointer()).Name()

					stack := make([]uintptr, 1024)

					length := runtime.Callers(2, stack)

					var skip int

					for i := 0; i < length; i++ {
						pc := stack[i]

						funcPtr := runtime.FuncForPC(pc)

						if funcPtr.Name() == actionFunctionName {
							skip = i + 1

							break
						}

					}

					pc, file, line, _ := runtime.Caller(skip)
					caller := runtime.FuncForPC(pc)

					// format the message
					panicMsg := fmt.Sprintf(`%sgovel: %sPanic message recovered: %s%s
	%sOrigin function: %s%s
	%sOrigin file: %s%s
	%sLine: %s%d

`, colorBlue, colorRed, colorReset, r, colorBlue, colorReset, caller.Name(), colorBlue, colorReset, file, colorBlue, colorReset, line)

					fmt.Print(panicMsg)
				}
			}

			// close the request body
			c.Request.Body.Close()

		}()

		// Middlewares must return two values.
		// If it returns 0, the request will continue as normal, but if it returns 1, the request will abort.

		var cancel int = 0

		if globalMiddlewares != nil {
			for _, middleware := range globalMiddlewares {
				if middleware(c) != cancel {
					cancel = 1
					break
				}
			}
		}

		if middlewares != nil && cancel == 0 {
			for _, middleware := range middlewares {
				if middleware(c) != cancel {
					cancel = 1
					break
				}
			}
		}

		if cancel == 0 {
			action(c)
		}
	}
}

func httpHandler(function routeFunction) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c := newContext(rw, r)

		function(c)
	}
}
