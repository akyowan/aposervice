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

type Storage struct {
	Bucket    string
	EndPoint  string
	AccessID  string
	AccessKey string
	CDNHost   string
}
