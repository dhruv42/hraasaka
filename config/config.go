package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Port     int
	DbUrl    string
	DbName   string
	CacheUrl string
}

func GetConfig() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed reading current dir: [%v]", err)
		return nil, err
	}
	// fmt.Println(wd)
	cf, err := ioutil.ReadFile(wd + "/config/config.json")
	if err != nil {
		log.Fatalf("Failed reading config.json: [%v]", err)
		return nil, err
	}
	var data *Config
	err = json.Unmarshal(cf, &data)
	if err != nil {
		log.Fatalf("Failed reading config.json: [%v]", err)
		return nil, err
	}
	return data, nil
}
