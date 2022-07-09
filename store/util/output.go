package util

import "fmt"

func Info(s any) {
	fmt.Printf("\033[1;34m%v\033[0m", s)
}

func Output(s any) {
	fmt.Printf("\033[1;36m%v\033[0m\n", s)
}

func Error(s any) {
	fmt.Printf("\033[1;31m%s\033[0m\n", s)
}
