package main

import (
	"flag"
	"fmt"

	"github.com/cjc7373/bitcoin_go/internal/utils"
)

func main() {
	configPath := flag.String("config", "config.yaml", "config path")
	flag.Parse()
	conf := utils.ParseConfig(*configPath)
	fmt.Println(conf)
}
