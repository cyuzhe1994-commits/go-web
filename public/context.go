package public

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Params     map[string]string
	StatusCode int
}

// 辅助方法：快速返回 JSON
func (c *Context) JSON(code int, obj interface{}) {
	c.StatusCode = code
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	json.NewEncoder(c.Writer).Encode(obj)
}

func (c *Context) Echo(code int, msg string) {
	c.StatusCode = code
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(msg))
}

func (c *Context) String(code int, msg string) {
	c.StatusCode = code
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(msg))
}

func (c *Context) HTML(code int, html string) {
	c.StatusCode = code
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) Query(key string, value ...string) string {
	v := c.Request.URL.Query().Get(key)
	if v == "" && len(value) > 0 {
		v = value[0]
	}
	return v
}

func (c *Context) Param(key string, value ...string) string {
	if c.Params != nil {
		v, ok := c.Params[key]
		if ok {
			return v
		}
	}
	if len(value) > 0 {
		return value[0]
	}
	return ""
}

func (c *Context) PostForm(key string, value ...string) string {
	v := c.Request.FormValue(key)
	if v == "" && len(value) > 0 {
		v = value[0]
	}
	return v
}

func (c *Context) BindJSON(obj interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	return decoder.Decode(obj)
}
