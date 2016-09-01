package main

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	// error, but not fatal
	logger.Error("Some error, but continue.")

	// just a warning, nothing to worry about ... yet
	logger.Warn("Warning....crash to ground imminent.")

	// this line WILL NOT appear because the default log.Level does not
	// log anything that is debug level
	// for debugging purpose...
	logger.Debug("Just debugging information.")

	// to enable debug
	logger.Level = logrus.DebugLevel

	// this LINE WILL APPEAR
	logger.Debug("Just debugging information.")

	// add fields to debug information
	logger.WithFields(logrus.Fields{
		"variable": "value",
		"username": "adam",
	}).Debug("Even useful debugging information.")

	// notice that debug as no color coding for fields...but not with Info!
	logger.WithFields(logrus.Fields{
		"variable": "value",
		"username": "adam",
	}).Info("Even useful debugging information.")

	// just for your information
	logger.Info("Oh, just FYI!")

	// Executes os.Exit(1) function after this line
	logger.Fatal("Abort!")

	// Executess panic() function after this line
	logger.Panic("Panic!")
}
