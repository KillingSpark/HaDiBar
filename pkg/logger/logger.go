package logger

import (
	"errors"
	"os"
	"path"

	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/text"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
)

func PrepareFallbackLogger() {
	if viper.GetString("FallbackLogger") == "Stderr" {
		PrepareLoggerStderr()
	}else{
		PrepareLoggerStdout()
	}
}

func PrepareLoggerStdout() {
	setLogLevel()
	log.SetHandler(text.New(os.Stdout))
}

func PrepareLoggerStderr() {
	setLogLevel()
	log.SetHandler(text.New(os.Stderr))
}

func setLogLevel() {
	level := viper.GetString("LoggingLevel")
	switch level {
		case "DEBUG": log.SetLevel(log.DebugLevel)
		case "INFO": log.SetLevel(log.InfoLevel)
		case "WARN": log.SetLevel(log.WarnLevel)
		case "ERROR": log.SetLevel(log.ErrorLevel)
		case "FATAL": log.SetLevel(log.FatalLevel)
		default: log.SetLevel(log.WarnLevel)
	}
}

func PrepareLoggerFromViper() error {
	if viper.GetString("Logger") == "Lumberjack" {
		logdir := viper.GetString("LogDir")
		size := viper.GetInt("LogMaxSize")
		bkps := viper.GetInt("LogMaxBackups")
		age := viper.GetInt("LogMaxAge")
		cmprss := viper.GetBool("LogCompress")
	
		err := PrepareLogger(logdir, size, age, bkps, cmprss)
		if err != nil {
			PrepareFallbackLogger()
			return err
		}
	} else {
		PrepareFallbackLogger()
	}
	return nil
}

func PrepareLogger(logdir string, size, backups, age int, compress bool) error {
	if logdir == "" {
		return errors.New("No logdir given")
	}
	_, err := os.Stat(logdir) 
	if err != nil {
		os.MkdirAll(logdir, 0666)
	}

	if size == 0 {
		return errors.New("No maxsize given")
	}
	if age == 0 {
		return errors.New("No maxage given")
	}

	setLogLevel()
	log.SetHandler(logfmt.New(&lumberjack.Logger{
		Filename:   path.Join(logdir, "hadibar.log"),
		MaxSize:    size,
		MaxBackups: backups,
		MaxAge:     age,
		Compress:   compress,
	}))
	return nil
}
