package logger

import (
	"os"

	"strings"

	logging "github.com/op/go-logging"
	"github.com/spf13/viper"
)

var (
	//Logger is the global logger
	Logger = logging.MustGetLogger("Debugging")
	format = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{color:reset} %{message}`)
)

//PrepareLogger puts a format and backend to the logger
func PrepareLogger() {
	formatted := logging.AddModuleLevel(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format))
	level := strings.ToUpper(viper.GetString("LoggingLevel"))

	switch level {
	case "DEBUG":
		formatted.SetLevel(logging.DEBUG, "debugging")
	case "CRITICAL":
		formatted.SetLevel(logging.CRITICAL, "critical")
	case "ERROR":
		formatted.SetLevel(logging.ERROR, "error")
	case "WARNING":
		formatted.SetLevel(logging.WARNING, "warning")
	default:
		formatted.SetLevel(logging.DEBUG, "debugging")
	}
	logging.SetBackend(formatted)
}
