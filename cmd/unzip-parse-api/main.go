package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"bitbucket.com/bbd/unzip-parse/usecases"

	"bitbucket.com/bbd/unzip-parse/entities"
	apiExternal "bitbucket.com/bbd/unzip-parse/externals/api"
	mysqlExternal "bitbucket.com/bbd/unzip-parse/externals/mysql"
)

const (
	FlagConfig = "c"
)

type Config struct {
	Port            string              `json:"port"`
	MySQL           mysqlExternal.MySQL `json:"mysql"`
	ApiURL          string              `json:"api_URL"`
	ImagerURL       string              `json:"imager_URL"`
	UnzipSrcDir     string              `json:"unzip_src_dir"`
	UnzipDestDir    string              `json:"unzip_dest_dir"`
	PhotoCategories []string            `json:"photo_categories"`
}

func main() {
	log.SetOutput(os.Stdout)
	var configFilePath string
	flag.StringVar(&configFilePath, FlagConfig, "./config.json", "json config file path, example: /foo/bar/config.json")
	flag.Parse()
	conf, err := readConfig(configFilePath)
	usecases.PanicAndCheckError(err)

	start := time.Now()

	initConfig(conf)

	mysqlExternal.OpenDB(conf.MySQL.User, conf.MySQL.Password, conf.MySQL.DataSourceName, conf.MySQL.Database)
	defer entities.DB.Close()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "parser up since %s", start)
	})

	http.Handle("/unzip-parse", apiExternal.MustMethod(http.MethodGet, http.HandlerFunc(apiExternal.UnzipParseHandler)))

	srv := &http.Server{
		Addr:         conf.Port,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  7 * time.Second,
	}

	log.Println("connect to db: ", conf.MySQL.DataSourceName)
	if err := srv.ListenAndServe(); err != nil {
		log.Println("Failed to start: ", err)
	}
}

func readConfig(configFilePath string) (*Config, error) {
	confContent, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	conf := new(Config)
	err = json.Unmarshal(confContent, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func initConfig(conf *Config) {
	usecases.ApiURL = (*conf).ApiURL
	usecases.ImagerURL = (*conf).ImagerURL
	usecases.UnzipSrcDir = (*conf).UnzipSrcDir
	usecases.UnzipDestDir = (*conf).UnzipDestDir
	usecases.PhotoCategories = (*conf).PhotoCategories
}
