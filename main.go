package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/kelseyhightower/envconfig"
)

type EnvConf struct {
	RedisHost string `envconfig:"REDIS_HOST"`
	RedisPort uint64 `envconfig:"REDIS_PORT"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func main() {
	envConf := new(EnvConf)
	err := envconfig.Process("", envConf)
	if err != nil {
		log.Fatalln("parse env configs,", err)
	} else if envConf.RedisHost == "" {
		log.Fatalln("REDIS_HOST is not defined")
	} else if envConf.RedisPort == 0 {
		log.Fatalln("REDIS_PORT is not defined")
	}

	redisAddr := envConf.RedisHost + ":" + strconv.FormatUint(envConf.RedisPort, 10)
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

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
