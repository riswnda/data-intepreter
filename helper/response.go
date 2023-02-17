package helper

import (
	"encoding/json"
	"sync"
)

var (
	rsp    *Response
	oneRsp sync.Once
)

// Response presentation contract object
type Response struct {
	Code    int         `json:"-"`
	Status  string      `json:"status,omitempty"`
	Entity  string      `json:"entity,omitempty"`
	State   string      `json:"state,omitempty"`
	Message interface{} `json:"message,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	lang    string      `json:"-"`
	msgKey  string
}

// MetaData represent meta data response for multi data
type MetaData struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	TotalPage  int64 `json:"total_page"`
	TotalCount int64 `json:"total_count"`
}

// WithCode setter response var name
func (r *Response) WithCode(c int) *Response {
	r.Code = c
	return r
}

// WithStatus setter status response
func (r *Response) WithStatus(s string) *Response {
	r.Status = s
	return r
}

// WithEntity setter entity response
func (r *Response) WithEntity(e string) *Response {
	r.Entity = e
	return r
}

// WithState setter state response
func (r *Response) WithState(s string) *Response {
	r.State = s
	return r
}

// WithData setter data response
func (r *Response) WithData(v interface{}) *Response {
	r.Data = v
	return r
}

// WithError setter error messages
func (r *Response) WithError(v interface{}) *Response {
	r.Errors = v
	return r
}

func (r *Response) WithMsgKey(v string) *Response {
	r.msgKey = v
	return r
}

// WithMeta setter meta data response
func (r *Response) WithMeta(v interface{}) *Response {
	r.Meta = v
	return r
}

// WithLang setter language response
func (r *Response) WithLang(v string) *Response {
	r.lang = v
	return r
}

// WithMessage setter custom message response
func (r *Response) WithMessage(v interface{}) *Response {
	if v != nil {
		r.Message = v
	}

	return r
}

func (r *Response) Byte() []byte {
	b, _ := json.Marshal(r)
	return b
}

// NewResponse initialize response
func NewResponse() *Response {
	oneRsp.Do(func() {
		rsp = &Response{}
	})

	// clone response
	x := *rsp

	return &x
}
