package logger

import (
	"log"
)

const (
	debug = iota
	info
	warn
	error
)

var logLevels = map[string]int{
	"debug": 0,
	"info":  1,
	"warn":  2,
	"error": 3,
}

var logLevel = 1

func SetLogLevel(level string) {
	logLevel = logLevels[level]
}

func check(level int) bool {
	return level >= logLevel
}

func Debug(msg interface{}) {
	if msg==nil {
		return
	}
	if check(debug) {
		log.Println("\033[;35mDEBUG\033[0m", msg)
	}
}

func Info(msg interface{}) {
	if msg==nil {
		return
	}
	if check(info) {
		log.Println("\033[;32mINFO\033[0m", msg)
	}
}

func Warn(msg interface{}) {
	if msg==nil {
		return
	}
	if check(warn) {
		log.Println("\033[;33mWARN\033[0m", msg)
	}
}

func Error(msg interface{}) {
	if msg==nil {
		return
	}
	if check(error) {
		log.Println("\033[;31mERROR\033[0m", msg)
	}
}
