package model

type Config struct {
	Addr     ServerAddr
	App      FeiShuApp
	Token    string
	Approval FeiShuApproval
	Event    FeiShuEvent
}

// 服务端口
type ServerAddr struct {
	Host string
	Port string
}

// 飞书应用
type FeiShuApp struct {
	Id     string
	Secret string
}

// 飞书事件
type FeiShuEvent struct {
	EncryptKey        string
	VerificationToken string
}

// 飞书审批
type FeiShuApproval struct {
	BlgsCode string
}
