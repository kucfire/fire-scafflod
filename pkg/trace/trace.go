package trace

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"sync"
)

type T interface {
	i()
	ID() string
	WithRequest(req *Request) *Trace
	WithResponse(resp *Response) *Trace
	// AppendDialog(dialog *Dialog) *Trace
	// AppendSQL(sql *SQL) *Trace
	AppendRedis(redis *Redis) *Trace
}

type Trace struct {
	mutex      sync.Mutex
	Identifier string    `json:"trace_id"` // 链路ID
	Request    *Request  `json:"request"`  // 请求信息
	Response   *Response `json:"response"` // 返回信息
	// ThirdPartyRequests []*Dialog `json:"third_party_requests"` // 调用第三方接口的信息
	// Debugs             []*Debug  `json:"debug"`                // 调试信息
	// SQLs               []*SQL    `json:"sqls"`                 // 执行的SQL信息
	Redis       []*Redis `json:"redis"`        // 执行的redis信息
	Success     bool     `json:"success"`      // 请求结果 true of false
	CostSeconds float64  `json:"cost_seconds"` // 执行市场(单位秒)
}

type Request struct {
	TTL        string      `json:"ttl"`         // 请求超时时间
	Method     string      `json:"method"`      // 请求方式
	DecodedURL string      `json:"decoded_url"` // 请求地址
	Header     interface{} `json:"header"`      // 请求 Header 信息
	Body       interface{} `json:"body"`        // 请求 Body 信息
}

type Response struct {
	Header          interface{} `json:"header"`                      // Header 信息
	Body            interface{} `json:"body"`                        // Body 信息
	BusinessCode    int         `json:"business_code,omitempty"`     // 业务码
	BusinessCodeMsg string      `json:"business_code_mgs,omitempty"` // 提示信息
	HttpCode        int         `json:"http_code"`                   // HTTP状态码

}

func New(id string) *Trace {
	if id == "" {
		buf := make([]byte, 0)
		io.ReadFull(rand.Reader, buf)

		id = hex.EncodeToString(buf)
	}

	return &Trace{
		Identifier: id,
	}
}

func (t *Trace) i()

func (t *Trace) ID() string {
	return t.Identifier
}

func (t *Trace) WithRequest(req *Request) *Trace {
	t.Request = req

	return t
}

func (t *Trace) WithResponse(resp *Response) *Trace {
	t.Response = resp

	return t
}

func (t *Trace) AppendRedis(redis *Redis) *Trace {
	if redis == nil {
		return t
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.Redis = append(t.Redis, redis)

	return t
}
