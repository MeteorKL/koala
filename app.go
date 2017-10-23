package koala

import (
	"net/http"
	"github.com/MeteorKL/koala/logger"
)

type App struct {
	routes          []Route
	notFoundHandler func(context *Context)
	renderPath      string
}

func NewApp() *App {
	return &App{
		notFoundHandler: func(context *Context) {
			context.Writer.WriteHeader(404)
			context.Writer.Write([]byte("404 Page Not Found!"))
		},
	}
}

func (app *App) SetNotFound(handler func(context *Context)) {
	app.notFoundHandler = handler
}

func (app *App) SetRenderPath(renderPath string) {
	app.renderPath = renderPath
}

func (app *App) Run(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info(r.Method + " " + r.URL.Path)
		app.route(initContext(app, w, r))
	})
	logger.Info("Listening on port " + addr)
	if err := http.ListenAndServe(":"+addr, nil); err != nil {
		logger.Fatal("ListenAndServe:")
		logger.Fatal(err)
	}
}

func (app *App) Static(pattern string, dir string) {
	http.Handle(pattern, http.StripPrefix(pattern, http.FileServer(http.Dir(dir))))
}

func (app *App) Handle(pattern string, handler func(context *Context)) {
	app.addRoute(pattern, "GET", handler)
}

func (app *App) Get(pattern string, handler func(context *Context)) {
	app.addRoute(pattern, "GET", handler)
}

func (app *App) Post(pattern string, handler func(context *Context)) {
	app.addRoute(pattern, "POST", handler)
}

func (app *App) Put(pattern string, handler func(context *Context)) {
	app.addRoute(pattern, "PUT", handler)
}

func (app *App) Delete(pattern string, handler func(context *Context)) {
	app.addRoute(pattern, "DELETE", handler)
}

func (app *App) Patch(pattern string, handler func(context *Context)) {
	app.addRoute(pattern, "PATCH", handler)
}
