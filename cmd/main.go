package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/opensvn/auth-client"
)

func main() {
	c := &client.Client{}
	err := c.Init()
	if err != nil {
		return
	}

	// Connect to the broker
	err = c.Connect()
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	// Messages will be handled through the callback so we really just need to wait until a shutdown is requested
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	fmt.Println("signal caught - exiting")

	// We could cancel the context at this point but will call Disconnect instead (this waits for autopaho to shutdown)
	err = c.Disconnect()
	if err != nil {
		log.Printf("%s\n", err)
	}

	fmt.Println("shutdown complete")
}
