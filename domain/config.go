package domain

type DBConf struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type ServerConf struct {
	AppKey                string
	ExternalListenAddress string
	InternalListenAddress string
}
