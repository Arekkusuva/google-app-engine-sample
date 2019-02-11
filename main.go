package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/go-redis/redis"
	"github.com/kelseyhightower/envconfig"
)

type EnvConf struct {
	RedisHost string `envconfig:"REDIS_HOST"`
	RedisPort uint64 `envconfig:"REDIS_PORT"`
	DbHost string `envconfig:"DB_HOST"`
	DbPort uint64 `envconfig:"DB_PORT"`
	DbName string `envconfig:"DB_NAME"`
	DbUser string `envconfig:"DB_USER"`
	DbPassword string `envconfig:"DB_PASSWORD"`
}

type User struct {
	Id uint `sql:",type:bigserial"`
	CreatedAt time.Time `sql:",notnull,default:CURRENT_TIMESTAMP"`

	Name string `sql:",notnull,type:varchar(120)"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func main() {
	conf := new(EnvConf)
	err := envconfig.Process("", conf)
	// TODO: Use validator
	if err != nil {
		log.Fatalln("parse env configs,", err)
	} else if conf.RedisHost == "" {
		log.Fatalln("REDIS_HOST is not defined")
	} else if conf.DbHost == "" {
		log.Fatalln("DB_HOST is not defined")
	} else if conf.DbName == "" {
		log.Fatalln("DB_NAME is not defined")
	} else if conf.DbUser == "" {
		log.Fatalln("DB_USER is not defined")
	} else if conf.DbPassword == "" {
		log.Fatalln("DB_PASSWORD is not defined")
	}

	if conf.RedisPort == 0 {
		conf.RedisPort = 6379
	}
	if conf.DbPort == 0 {
		conf.DbPort = 5432
	}

	// Connect to Memorystore Redis instance
	redisAddr := conf.RedisHost + ":" + strconv.FormatUint(conf.RedisPort, 10)
	log.Println("connect to redis on", redisAddr)
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		DB:       0,
	})
	pingResp, err := client.Ping().Result()
	if err != nil {
		log.Fatalln("ping failed,", err)
	} else if pingResp != "PONG" {
		log.Fatalln("ping failed")
	}

	// Connect to database
	dbAddr := conf.DbHost + ":" + strconv.FormatUint(conf.DbPort, 10)
	log.Println("connect to database on", dbAddr)
	db := pg.Connect(&pg.Options{
		Addr:     dbAddr,
		User:     conf.DbUser,
		Database: conf.DbName,
		Password: conf.DbPassword,
	})

	// Create users table if not exists
	createOptions := orm.CreateTableOptions{
		FKConstraints: true,
		IfNotExists:   true,
	}
	err = db.CreateTable(&User{}, &createOptions)
	if err != nil {
		log.Println("create users table,", err)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
