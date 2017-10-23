package koala

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
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

var VaildSuffix = []string{
	".doc", ".docx",
	".ppt", ".pptx",
	".xls", ".xlsx",
	".txt", ".pdf",
}

func SavePostFile(r *http.Request, key string, dir string) (string, string, string, error) {
	file, handle, err := r.FormFile(key)
	if err != nil {
		return "", "", "", err
	}
	filename := handle.Filename
	suffix := path.Ext(filename)
	flag := false
	for _, s := range VaildSuffix {
		if suffix == s {
			flag = true
		}
	}
	if !flag {
		return "", "", "", errors.New("不支持的文件后缀名")
	}
	attachPath := dir + filename
	filepath := "./static/upload" + dir + filename
	os.MkdirAll(path.Dir(filepath), 0777)
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}
	_, err = io.Copy(f, file)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}
	defer f.Close()
	defer file.Close()
	return attachPath, filename, suffix, nil
}
