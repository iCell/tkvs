package main

import (
	"errors"
	"github.com/iCell/tkvs/store"
	"github.com/iCell/tkvs/store/tkvs"
	"github.com/iCell/tkvs/store/util"
	"os"
	"strings"
)

var (
	ErrInvalidArguments = errors.New("invalid arguments")
	ErrUnSupportedCmd   = errors.New("unsupported command")
)

type Session struct {
	Kv store.IStore
}

func NewSession() *Session {
	store := tkvs.NewKvStore()
	return &Session{Kv: store}
}

func (_s *Session) Process(input string) {
	mappings := map[string]func(args []string){
		"GET":      _s.get,
		"SET":      _s.set,
		"COUNT":    _s.count,
		"DELETE":   _s.delete,
		"BEGIN":    _s.begin,
		"COMMIT":   _s.commit,
		"ROLLBACK": _s.rollback,
	}

	splits := strings.Split(input, " ")
	components := make([]string, 0, len(splits))
	for _, split := range splits {
		if len(split) == 0 {
			continue
		}
		components = append(components, split)
	}
	if len(components) == 0 {
		return
	}

	cmd, args := strings.ToUpper(components[0]), components[1:]
	if cmd == "EXIT" {
		os.Exit(0)
	}

	f, exist := mappings[cmd]
	if !exist {
		util.Error("unsupported command")
		return
	}
	f(args)
}

func (_s *Session) set(args []string) {
	if len(args) > 0 && len(args)%2 != 0 {
		util.Error("invalid key value pair")
		return
	}
	for i := 0; i < len(args); i += 2 {
		_s.Kv.Set(args[i], args[i+1])
	}
}

func (_s *Session) get(args []string) {
	if len(args) != 1 {
		util.Error("you should provide at least one key")
		return
	}
	v, exist := _s.Kv.Get(args[0])
	if !exist {
		util.Error("key not set")
		return
	}
	util.Output(v)
}

func (_s *Session) delete(args []string) {
	if len(args) < 1 {
		util.Error("you should provide at least one key")
		return
	}
	for _, arg := range args {
		_s.Kv.Delete(arg)
	}
}

func (_s *Session) count(args []string) {
	if len(args) != 1 {
		util.Error("you should provide one key")
		return
	}
	util.Output(_s.Kv.Count(args[0]))
}

func (_s *Session) begin(_ []string) {
	if err := _s.Kv.Begin(); err != nil {
		util.Error(err.Error())
	}
}

func (_s *Session) rollback(_ []string) {
	if err := _s.Kv.Rollback(); err != nil {
		util.Error(err.Error())
	}
}

func (_s *Session) commit(_ []string) {
	if err := _s.Kv.Commit(); err != nil {
		util.Error(err.Error())
	}
}
