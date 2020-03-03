package util

import (
	"github.com/cihub/seelog"
	"fmt"
)

var Logger seelog.LoggerInterface

func init() {
	// Disable logger by default.
	DisableLog()
	loadAppConfig()

}

// DisableLog disables all library log output.
func DisableLog() {
	Logger = seelog.Disabled
}

// UseLogger uses a specified seelog.LoggerInterface to output library log.
// Use this func if you are using Seelog logging system in your app.
func UseLogger(newLogger seelog.LoggerInterface) {
	Logger = newLogger
}

// Call this before app shutdown
func FlushLog() {
	Logger.Flush()
}

//
func loadAppConfig() {
	Logger, err := seelog.LoggerFromConfigAsFile("conf/seelog.xml")
	if err != nil {
		fmt.Println(err)
		return
	}
	UseLogger(Logger)
}

