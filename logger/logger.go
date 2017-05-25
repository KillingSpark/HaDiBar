package logger

import (
	"os"

	logging "github.com/op/go-logging"
)

var (
	//Logger is the global logger
	Logger = logging.MustGetLogger("Debugging")
	format = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{color:reset} %{message}`)
)

//PrepareLogger puts a format and backend to the logger
func PrepareLogger() {
	formatted := logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format)
	logging.SetBackend(formatted)
}
