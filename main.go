package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/viper"
)

func main() {
	viper.SetEnvPrefix("MINARIA")
	viper.AutomaticEnv()

	server := NewServer()
	server.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	fmt.Println("got", sig)
	server.ShutDown()
}
