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
	"path"
	"os"
	"errors"
	"strings"
)

type Context struct {
	app       *App
	Writer    http.ResponseWriter
	Request   *http.Request
	urlSlice  map[string]string
	query     url.Values
	bodyQuery url.Values
}
const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

func initContext(app *App, w http.ResponseWriter, r *http.Request) *Context {
	var err error
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		err = r.ParseMultipartForm(defaultMaxMemory)
		logger.Warn(err)
	}
	r.ParseForm()
	c := &Context{
		app:       app,
		Writer:    w,
		Request:   r,
		bodyQuery: r.PostForm,
	}
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

func (c *Context) GetBody(MaxMemory int64) []byte {
	safe := &io.LimitedReader{R: c.Request.Body, N: MaxMemory}
	requestBody, _ := ioutil.ReadAll(safe)
	c.Request.Body.Close()
	bf := bytes.NewBuffer(requestBody)
	c.Request.Body = ioutil.NopCloser(bf)
	return requestBody
}

func (c *Context) GetBodyUnsafe() []byte {
	requestBody, _ := ioutil.ReadAll(c.Request.Body)
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

var DEFAULT_VAILD_SUFFIX = []string{
	".doc", ".docx",
	".ppt", ".pptx",
	".xls", ".xlsx",
	".txt", ".pdf",
}

var VaildSuffix = DEFAULT_VAILD_SUFFIX

func (c *Context) SavePostFile(key string, dir string, vaildSuffix []string) (filename string, suffix string, err error) {
	file, handle, err := c.Request.FormFile(key)
	if err != nil {
		err = errors.New("上传文件失败")
		return
	}
	filename = handle.Filename
	suffix = path.Ext(filename)
	flag := false
	if vaildSuffix == nil {
		vaildSuffix = VaildSuffix
	}
	for _, s := range vaildSuffix {
		if suffix == s {
			flag = true
		}
	}
	if !flag {
		err = errors.New("不支持的文件后缀名")
		return
	}
	filepath := dir + filename
	os.MkdirAll(path.Dir(filepath), 0777)

	_, err = os.Stat(filepath)
	if !os.IsNotExist(err) {
		err = errors.New("该文件已存在")
		return
	}

	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	_, err = io.Copy(f, file)
	if err != nil {
		return
	}
	defer f.Close()
	defer file.Close()
	return
}
