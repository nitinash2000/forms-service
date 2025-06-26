package config

type Config struct {
	AppName    string `json:"app_name"`
	ServerPort string `json:"server_port"`
	DbUrl      string `json:"db_url"`
}
