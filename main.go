package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/go-nats"
)

type EnvConf struct {
	BusHost string `envconfig:"BUS_HOST"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func main() {
	envConf := new(EnvConf)
	err := envconfig.Process("", envConf)
	if err != nil {
		log.Fatalln("parse env configs,", err)
	} else if envConf.BusHost == "" {
		log.Fatalln("bus host is not defined")
	}

	log.Println("connect to bus on", envConf.BusHost)
	_, err = nats.Connect(envConf.BusHost)
	if err != nil {
		log.Fatalln("nats connect,", err)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
