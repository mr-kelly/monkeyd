package main

import (
	"flag"
	"fmt"
	"github.com/mr-kelly/monkeyd/monkeyd"
    "os"
    "log"
)

func main() {
    log.SetPrefix("[monkeyd]")
	configFile := flag.String("config", "config.toml", "`config file path`")
	fmt.Println("=== Monkeyd v0.1")
    fmt.Println("=== Port forward tool")
    fmt.Println("")

    flag.Parse()

    if _, err := os.Stat(*configFile); err != nil {
        if os.IsNotExist(err) {
            log.Fatalf("Not found config file '%s', please use `-config` tell me the config file path", *configFile);
            return
        }
    }

    fmt.Printf("Begin read config...%s\n", *configFile)
    monkeyd := monkeyd.New(*configFile)
    monkeyd.Run("server")
}