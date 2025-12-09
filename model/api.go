package model

type SuccessResp struct {
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

type ErrorResp struct {
	Msg   string `json:"msg"`
	Error any    `json:"error,omitempty"`
}
