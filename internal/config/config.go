package config

type ServerConfig struct {
	ConnectionString string `json:"connectionString"`
	DbType           string `json:"dbType"`
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
