// internal/logger/logger.go
package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Log adalah instance global dari logrus logger
var Log *logrus.Logger

func init() {
	Log = logrus.New()

	// Set output ke stdout (bisa juga ke file atau sistem logging eksternal)
	Log.SetOutput(os.Stdout)

	// Set formatter (JSONFormatter untuk structured logging, atau TextFormatter untuk development)
	logLevelStr := strings.ToLower(os.Getenv("LOG_LEVEL")) // Ambil dari env, default ke info
	formatter := os.Getenv("LOG_FORMATTER")                // Ambil dari env, default ke text

	if formatter == "json" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00", // ISO8601
			// PrettyPrint: true, // Bisa diaktifkan untuk development jika format JSON sulit dibaca
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			// ForceColors: true, // Bisa diaktifkan jika terminal mendukung
			// DisableColors: false,
		})
	}

	// Set level log
	level, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		Log.SetLevel(logrus.InfoLevel) // Default ke Info jika parsing gagal
		Log.Warnf("Invalid LOG_LEVEL '%s', defaulting to 'info'", logLevelStr)
	} else {
		Log.SetLevel(level)
	}

	// Tambahkan field default (opsional)
	// Log = Log.WithFields(logrus.Fields{
	// 	"service": "auth-gox-vue-example",
	// }).Logger // Jika ingin selalu ada field 'service'

	Log.Info("Logger initialized")
	Log.Infof("Log level set to: %s", Log.GetLevel().String())
	Log.Infof("Log formatter set to: %s", formatter)
}
