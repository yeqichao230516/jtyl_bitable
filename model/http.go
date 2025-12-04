package model

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type Req struct {
	Data any `json:"data"`
}
