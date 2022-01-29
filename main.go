package main

import (
	"fmt"
	"github.com/r2binx/heimboard-wakeup-go/config"
	"github.com/r2binx/heimboard-wakeup-go/router"
	"github.com/r2binx/heimboard-wakeup-go/util"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load(".config")
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}
	conf, _ := config.NewConfig(cfg)
	wake := util.NewWakeup(conf)

	go func() {
		for {
			wake.CheckSchedule()
			time.Sleep(30 * time.Second)
		}
	}()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins(conf.Origins)
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Println("Listening on port:", conf.Port)
	rtr := router.New(conf, wake)
	if err := http.ListenAndServe(fmt.Sprint(":", conf.Port), handlers.CORS(originsOk, headersOk, methodsOk)(rtr)); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
