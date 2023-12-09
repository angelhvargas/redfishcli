package logger

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

// Log is the global logger.
var Log = logrus.New()

func init() {
	// Set the log output, level, formatter, etc.
	file, err := os.OpenFile("./logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Could not open log file:", err)
	}

	// Set Logrus output to the file
	log.SetOutput(file)
	Log.SetLevel(logrus.InfoLevel)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
