package koala

import (
	"strings"
)

// 路由定义
type Route struct {
	slice   []string
	method  string
	handler func(context *Context)
}

func (app *App) addRoute(pattern string, method string, handler func(context *Context)) {
	slice := strings.Split(pattern, "/")
	app.routes = append(app.routes, Route{slice: slice[1:], method: method, handler: handler})
}

func (app *App) route(context *Context) {
	requestPath := context.Request.URL.Path
	isFound := false
	for i := 0; i < len(app.routes); i++ {
		route := app.routes[i]
		if route.method == context.Request.Method {
			url := strings.Split(requestPath, "/")[1:]
			if len(url) != len(route.slice) {
				continue
			}
			if len(route.slice[0]) == 0 {
				if len(url[0]) != 0 {
					continue
				}
			} else {
				matched := true
				for i := 0; i < len(route.slice); i++ {
					if route.slice[i][0] == ':' {
						context.urlSlice[route.slice[i][1:]] = url[i]
					} else if route.slice[i] == url[i] {
						continue
					} else {
						matched = false
						break
					}
				}
				if !matched {
					continue
				}
			}
			isFound = true
			route.handler(context)
			break
		}
	}

	if !isFound {
		app.notFoundHandler(context)
	}
}
