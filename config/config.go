package config

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// ConfigFile is the default config file
var ConfigFile = "./config.yml"

// GlobalConfig is the global config
type GlobalConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Twitter  TwitterConfig  `yaml:"twitter"`
}

// ServerConfig is the server config
type ServerConfig struct {
	Addr               string
	Mode               string
	Version            string
	StaticDir          string `yaml:"static_dir"`
	Cache              string `yaml:"cache"`
	ViewDir            string `yaml:"view_dir"`
	LogDir             string `yaml:"log_dir"`
	UploadDir          string `yaml:"upload_dir"`
	MaxMultipartMemory int64  `yaml:"max_multipart_memory"`
	KeySecure          string `yaml:"key_secure"`
}

// DatabaseConfig is the database config
type DatabaseConfig struct {
	Dialect      string
	DSN          string `yaml:"datasource"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

// DatabaseConfig is the database config
type TwitterConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	ApiKey       string `yaml:"api_key"`
	ApiSecret    string `yaml:"api_secret"`
	BearerToken  string `yaml:"bearer_token"`
}

// global configs
var (
	Global   GlobalConfig
	Server   ServerConfig
	Database DatabaseConfig
	Twitter  TwitterConfig
)

// Load config from file
func Load(file string) (GlobalConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("%v", err)
		return Global, err
	}

	err = yaml.Unmarshal(data, &Global)
	if err != nil {
		log.Printf("%v", err)
		return Global, err
	}

	Server = Global.Server
	Database = Global.Database
	Twitter = Global.Twitter

	// set log dir flag for glog
	flag.CommandLine.Set("log_dir", Server.LogDir)

	return Global, nil
}

// loads configs
func init() {
	if os.Getenv("config") != "" {
		ConfigFile = os.Getenv("config")
	}
	Load(ConfigFile)
}
