package model

type SuccessResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

type ErrorResp struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Error any    `json:"error,omitempty"`
}
