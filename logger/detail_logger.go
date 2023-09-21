package logger

import (
	"context"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

type key string

const (
	// ActivityID generate uuid
	ActivityID key = "activityID"
)

// APILogger api logger
func APILogger(ctx context.Context, method, path string, body interface{}) {
	data := map[string]interface{}{
		"activityID": ctx.Value(ActivityID),
		"activity":   path,
		"data":       body,
	}
	if Logger.Name != "test" {
		Logger.WithFields(
			logrus.Fields{
				"method": method,
				"path":   path,
				"body":   data,
			},
		).Info("HTTP")
	}
}

// DetailLoggerInfo detail logger info
func DetailLoggerInfo(ctx context.Context, command, details string, logData interface{}) {
	data := map[string]interface{}{
		"activityID": ctx.Value(ActivityID),
		"activity":   details,
		"data":       logData,
	}
	if Logger.Name != "test" {
		Logger.WithFields(
			logrus.Fields{
				"command": command,
				"details": data,
			},
		).Info("Function")
	}
}

// DetailLoggerError detail logger error
func DetailLoggerError(ctx context.Context, command, details string, err ...interface{}) {
	reflectData := reflect.ValueOf(err)
	dataString := fmt.Sprintf("%v", reflectData)
	data := map[string]interface{}{
		"activityID": ctx.Value(ActivityID),
		"activity":   details,
		"data":       dataString,
	}
	if Logger.Name != "test" {
		Logger.WithFields(
			logrus.Fields{
				"command": command,
				"details": data,
			},
		).Error("Function")
	}
}

// DetailLoggerWarn detail logger warning
func DetailLoggerWarn(ctx context.Context, command, details string, warn ...interface{}) {
	reflectData := reflect.ValueOf(warn)
	dataString := fmt.Sprintf("%v", reflectData.Interface())
	data := map[string]interface{}{
		"activityID": ctx.Value(ActivityID),
		"activity":   details,
		"data":       dataString,
	}
	if Logger.Name != "test" {
		Logger.WithFields(
			logrus.Fields{
				"command": command,
				"details": data,
			},
		).Warn("Function")
	}
}
