package koala

import (
	"net/http"
	"strconv"
	"net/url"
	"github.com/MeteorKL/koala/logger"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"bytes"
)

type Context struct {
	app       *App
	Writer    http.ResponseWriter
	Request   *http.Request
	urlSlice  map[string]string
	query     url.Values
	bodyQuery url.Values
}

func initContext(app *App, w http.ResponseWriter, r *http.Request) *Context {
	r.ParseForm()
	c := &Context{
		app:       app,
		Writer:    w,
		Request:   r,
		bodyQuery: r.PostForm,
	}
	var err error
	if c.query, err = url.ParseQuery(r.URL.RawQuery); err != nil {
		logger.Info(err)
	}
	return c
}

func (c *Context) Render(file string, data interface{}) {
	t, err := template.New(file).ParseFiles(c.app.renderPath + file)
	logger.Error(err)
	t.Execute(c.Writer, data)
}

func (c *Context) GetQueryInt(key string, def int) (int, error) {
	return strconv.Atoi(c.query.Get(key))
}

func (c *Context) GetQueryIntOrDefault(key string, def int) int {
	ret, err := strconv.Atoi(c.query.Get(key))
	if err != nil {
		ret = def
	}
	return ret
}

func (c *Context) GetQueryString(key string) string {
	return c.query.Get(key)
}

func (c *Context) GetQueryStringOrDefault(key string, def string) string {
	if c.query == nil {
		return def
	}
	return c.query[key][0]
}

func (c *Context) GetBodyAsJson(MaxMemory int64) []byte {
	safe := &io.LimitedReader{R: c.Request.Body, N: MaxMemory}
	requestBody, _ := ioutil.ReadAll(safe)
	c.Request.Body.Close()
	bf := bytes.NewBuffer(requestBody)
	c.Request.Body = ioutil.NopCloser(bf)
	return requestBody
}

func (c *Context) GetBodyQueryInt(key string, def int) (int, error) {
	return strconv.Atoi(c.bodyQuery.Get(key))
}

func (c *Context) GetBodyQueryIntOrDefault(key string, def int) int {
	ret, err := strconv.Atoi(c.bodyQuery.Get(key))
	if err != nil {
		ret = def
	}
	return ret
}

func (c *Context) GetBodyQueryString(key string) string {
	return c.bodyQuery.Get(key)
}

func (c *Context) GetBodyQueryStringOrDefault(key string, def string) string {
	if c.bodyQuery == nil {
		return def
	}
	return c.bodyQuery[key][0]
}

func (c *Context) WriteJSON(data interface{}) {
	ret, err := json.Marshal(data)
	logger.Warn(err)
	c.Writer.Write(ret)
}

func (c *Context) Relocation(URL string) {
	// w.Header().Set("Location", URL)
	t, err := template.New("x").Parse("<script>window.location.href='" + URL + "';</script>")
	logger.Warn(err)
	t.Execute(c.Writer, nil)
}

func (c *Context) Back() {
	t, err := template.New("x").Parse("<script>history.go(-1);</script>")
	logger.Warn(err)
	t.Execute(c.Writer, nil)
}
