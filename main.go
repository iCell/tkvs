package main

import (
	"bufio"
	"fmt"
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
	fmt.Println("the transactional key-value store started...")
	for scanner.Scan() {
		input := scanner.Text()
		if err := session.Process(input); err != nil {
			fmt.Println(err)
		}
	}
}
