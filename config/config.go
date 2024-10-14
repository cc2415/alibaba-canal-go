package config

type AlibabaCanal struct {
	Address     string
	Port        int
	Username    string
	Password    string
	Destination string

	Database string
}
type Config struct {
	EsAddress    []string
	AlibabaCanal AlibabaCanal
}

var AppConfig Config

func init() {
	AppConfig.EsAddress = []string{"http://127.0.0.1:9200/", "http://127.0.0.1:9201/", "http://127.0.0.1:9202/"}
	AppConfig.AlibabaCanal.Address = "192.168.199.165"
	AppConfig.AlibabaCanal.Port = 11111
	AppConfig.AlibabaCanal.Username = "alibaba_canal"
	AppConfig.AlibabaCanal.Password = "abcd"
	AppConfig.AlibabaCanal.Destination = "destination"
	AppConfig.AlibabaCanal.Database = "database"
}
