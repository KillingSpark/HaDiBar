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

func PrepareLoggerStdout() {
	log.SetLevel(log.DebugLevel)
	log.SetHandler(text.New(os.Stderr))
}

func PrepareLoggerFromViper() error {
	logdir := viper.GetString("LogDir")
	size := viper.GetInt("LogMaxSize")
	bkps := viper.GetInt("LogMaxBackups")
	age := viper.GetInt("LogMaxAge")
	cmprss := viper.GetBool("LogCompress")

	return PrepareLogger(logdir, size, age, bkps, cmprss)
}

func PrepareLogger(logdir string, size, backups, age int, compress bool) error {
	if logdir == "" {
		return errors.New("No logdir given")
	}
	if size == 0 {
		return errors.New("No maxsize given")
	}
	if age == 0 {
		return errors.New("No maxage given")
	}

	log.SetLevel(log.DebugLevel)
	log.SetHandler(logfmt.New(&lumberjack.Logger{
		Filename:   path.Join(logdir, "hadibar.log"),
		MaxSize:    size,
		MaxBackups: backups,
		MaxAge:     age,
		Compress:   compress,
	}))
	return nil
}
