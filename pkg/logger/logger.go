package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Log is the global logger.
var Log = logrus.New()

func init() {
	// Set the log output, level, formatter, etc.
	file, err := os.OpenFile("./logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Use logrus's standard logger to log this error, as our custom logger isn't set up yet.
		logrus.StandardLogger().Fatal("Could not open log file:", err)
	}

	// Set Logrus output to the file
	Log.Out = file // Use Log.Out to set the output for logrus logger
	Log.SetLevel(logrus.InfoLevel)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
