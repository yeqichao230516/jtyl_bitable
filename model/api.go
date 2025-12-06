package model

import "net/http"

type ApiReq struct {
	Body any
}

type ApiResp struct {
	StatusCode int         `json:"-"`
	Header     http.Header `json:"-"`
	RawBody    []byte      `json:"-"`
}
