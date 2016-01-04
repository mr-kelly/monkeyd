package main

import (
	"flag"
	"fmt"
	"github.com/mr-kelly/monkeyd/monkeyd"
	"github.com/op/go-logging"
	"os"
)

var log = logging.MustGetLogger("monkeyd:main")

func main() {
	configFile := flag.String("config", "config.toml", "`config file path`")
	log.Info("=== Monkeyd v0.1")
	log.Info("=== Port forward tool")
	log.Info("")

	flag.Parse()

	if _, err := os.Stat(*configFile); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Not foune config file '%s', please use `-config` tell me the config file path", *configFile)
			return
		}
	}

	fmt.Printf("Begin read config...%s\n", *configFile)
	newMonkeyd, err := monkeyd.New(*configFile)
	if err != nil {
		fmt.Printf("Error on new monkeyd %s", err.Error())
		return
	}
	newMonkeyd.Run("server")
}
