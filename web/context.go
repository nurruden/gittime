package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req          *http.Request
	Resp         http.ResponseWriter
	PathParams   map[string]string
	queryValues  url.Values
	MatchedRoute string
}

func (c *Context) RespJSONOK(val any) error {
	return c.RespJSON(http.StatusOK, val)
}

func (c *Context) SetCookie(ck *http.Cookie) {
	http.SetCookie(c.Resp, ck)
}

func (c *Context) RespJSON(status int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Resp.WriteHeader(status)
	c.Resp.Header().Set("Content-Type", "application/json")
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
	n, err := c.Resp.Write(data)
	if n != len(data) {
		return errors.New("web: did not write all data")
	}
	return err
}

func (c *Context) BindJSON(val any) error {
	//if val == nil {
	//	return errors.New("web:input is nil")
	//}
	if c.Req.Body == nil {
		return errors.New("web: body is nil")
	}

	decoder := json.NewDecoder(c.Req.Body)
	//decoder.UseNumber()
	//decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}

func (c *Context) FormValue(key string) (string, error) {
	err := c.Req.ParseForm()
	if err != nil {
		return "", err
	}
	//vals, ok := c.Req.Form[key]
	//if !ok {
	//	return "", errors.New("web: key not found")
	//}
	return c.Req.FormValue(key), nil
}

func (c *Context) QueryValue(key string) (string, error) {
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}
	vals, ok := c.queryValues[key]
	if !ok || len(vals) == 0 {
		return "", errors.New("web: key not exist")
	}
	return vals[0], nil
	//return c.queryValues.Get(key), nil
}

func (c *Context) PathValue(key string) (string, error) {
	val, ok := c.PathParams[key]
	if !ok {
		return "", errors.New("web: key not exist")
	}
	return val, nil
}
func (c *Context) PathValueV1(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{
			err: errors.New("web: key not exist"),
		}
	}
	return StringValue{
		val: val,
	}
}

func (c *Context) QueryValueV1(key string) StringValue {
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}
	vals, ok := c.queryValues[key]
	if !ok || len(vals) == 0 {
		return StringValue{
			err: errors.New("web: key not exist"),
		}
	}
	return StringValue{
		val: vals[0],
	}
	//return c.queryValues.Get(key), nil
}

type StringValue struct {
	val string
	err error
}

func (s StringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}
