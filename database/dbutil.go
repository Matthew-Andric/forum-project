package database

import (
	"database/sql"
	"fmt"
	"forum/util"
	"io/ioutil"

	"github.com/gorilla/securecookie"
	"gopkg.in/yaml.v2"
)

var (
	config YAMLFile
	db     *sql.DB
)

func StartDB() {
	err := ParseConfig()
	if err != nil {
		panic(err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("DB Opened")

	util.S = *securecookie.New([]byte(config.Secret.Hash), []byte(config.Secret.Block))

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("DB pinged")
}

func CloseDB() {
	db.Close()
}

func ParseConfig() error {
	cfg, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(cfg, &config)
	return err
}
