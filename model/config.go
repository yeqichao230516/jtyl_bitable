package model

type Config struct {
	Addr  ServerAddr
	App   FeiShuApp
	Token string
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
