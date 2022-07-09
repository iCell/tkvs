package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	ErrInvalidArguments = errors.New("invalid arguments")
	ErrUnSupportedCmd   = errors.New("unsupported command")
)

type Session struct {
	Kv *KvStore
}

func NewSession() *Session {
	return &Session{Kv: NewKvStore()}
}

func (_s *Session) Process(input string) error {
	splits := strings.Split(input, " ")
	components := make([]string, 0, len(splits))
	for _, split := range splits {
		if len(split) == 0 {
			continue
		}
		components = append(components, split)
	}
	if len(components) == 0 {
		return nil
	}

	switch strings.ToUpper(components[0]) {
	case "SET":
		if len(components[1:])%2 != 0 {
			return ErrInvalidArguments
		}
		_s.set(components[1:])
	case "GET":
		if len(components[1:]) != 1 {
			return ErrInvalidArguments
		}
		_s.get(components[1:])
	case "DELETE":
		if len(components[1:]) < 1 {
			return ErrInvalidArguments
		}
		_s.delete(components[1:])
	case "COUNT":
		if len(components[1:]) < 1 {
			return ErrInvalidArguments
		}
		_s.count(components[1:])
	case "BEGIN":
		_s.begin()
	case "COMMIT":
		_s.commit()
	case "ROLLBACK":
		_s.rollback()
	case "EXIT":
		os.Exit(0)
	default:
		return ErrUnSupportedCmd
	}
	return nil
}

func (_s *Session) set(args []string) {
	for i := 0; i < len(args); i += 2 {
		_s.Kv.Set(args[i], args[i+1])
	}
}

func (_s *Session) get(args []string) {
	v, exist := _s.Kv.Get(args[0])
	if !exist {
		fmt.Println("key not set")
		return
	}
	fmt.Println(v)
}

func (_s *Session) delete(args []string) {
	for _, arg := range args {
		_s.Kv.Delete(arg)
	}
}

func (_s *Session) count(args []string) {
	results := make([]string, 0, len(args))
	for _, arg := range args {
		results = append(results, strconv.Itoa(_s.Kv.Count(arg)))
	}
	fmt.Println(strings.Join(results, " "))
}

func (_s *Session) begin() {
	if err := _s.Kv.Begin(); err != nil {
		fmt.Println(err)
	}
}

func (_s *Session) rollback() {
	if err := _s.Kv.Rollback(); err != nil {
		fmt.Println(err)
	}
}

func (_s *Session) commit() {
	if err := _s.Kv.Commit(); err != nil {
		fmt.Println(err)
	}
}
