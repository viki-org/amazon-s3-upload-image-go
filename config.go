package s3

import (
  "encoding/json"
  "github.com/viki-org/glog"
  "io/ioutil"
  "os"
)

type Config struct {
	AccessKey string
	Secret    string
	Bucket    string
}

var (
	config *Config
)

func GetConfig() *Config {
	return config
}

func init() {
	loadConfig()
}

func loadConfig() {
  configPath := GetPath("config.json")
	if _, err := os.Stat(configPath); err != nil {
		glog.Fatalf("Cannot find config file: err= %v", err)
	}

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		glog.Infof("Cannot read config file; err= %v", err)
	}

	config := make(map[string]*Config)
	err = json.Unmarshal(file, &config)
	if err != nil {
		glog.Infof("Cannot unmarshal config file; err= %v", err)
	}
}

func GetPath(filename string) string {
	rootPath := os.Getenv("USERS_ROOT")
	if len(rootPath) != 0 {
		return rootPath + "/" + filename
	}
	return filename
}
