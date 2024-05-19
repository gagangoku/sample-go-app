package main

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

var _initer string = func() string {
	var zLogger = GetLogger()
	zerolog.DefaultContextLogger = &zLogger
	return ""
}()

func GetLogger() zerolog.Logger {
	zerolog.CallerMarshalFunc = shortFileLog

	var base zerolog.Logger
	if *prettyLogFlag {
		base = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime})
	} else {
		base = zerolog.New(os.Stderr)
	}
	return base.
		With().Caller().Logger().
		With().Timestamp().Logger()
}

func shortFileLog(pc uintptr, file string, line int) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return file + ":" + strconv.Itoa(line)
}
