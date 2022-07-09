package main

import (
	"bufio"
	"github.com/iCell/tkvs/store/util"
	"os"
	"os/signal"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			os.Exit(0)
		}
	}()

	session := NewSession()
	scanner := bufio.NewScanner(os.Stdin)
	util.Info("the transactional key-value Kvs started...\n")
	for {
		util.Info("> ")
		scanner.Scan()
		input := scanner.Text()
		if err := session.Process(input); err != nil {
			util.Error(err.Error())
		}
	}
}
