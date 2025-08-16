package config

type ServerConfig struct {
	Server   string `json:"server"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	DbType   string `json:"dbType"`
}

type FileConfig struct {
	File       string `json:"file"`
	Out        string `json:"out"`
	Procs      string `json:"procs"`
	NumWorkers int    `json:"workers"`
}

type Config struct {
	Server *ServerConfig
	File   *FileConfig
}
