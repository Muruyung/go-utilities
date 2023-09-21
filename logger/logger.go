package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Muruyung/go-utils/converter"

	"github.com/sirupsen/logrus"
)

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.WarnLevel, logrus.DebugLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}

type logger struct {
	*logrus.Logger
	enabled bool
	Env     string
	Name    string
	Date    string
}

type formatter struct {
	env     string
	isLocal bool
}

// Format custom formatter
func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	const formatColor = "\x1b[%dm"

	// output buffer
	b := &bytes.Buffer{}
	defer b.Reset()
	levelColor := getColorByLevel(entry.Level)

	_, _ = fmt.Fprintf(b, formatColor, levelColor)
	// Log Date Time
	now := time.Now().Format(time.RFC3339)
	b.WriteString("[")
	b.WriteString(now)
	b.WriteString("]")

	// Log level
	b.WriteString("[")
	level := strings.ToUpper(entry.Level.String())
	b.WriteString(level)
	b.WriteString("]")

	// Log direction
	// if entry.HasCaller() && f.env == "development" {
	// 	b.WriteString("[")
	// 	if f.isLocal {
	// 		_, _ = fmt.Fprintf(b, formatColor, levelColor)
	// 	}
	// 	fmt.Fprintf(
	// 		b,
	// 		"%s:%d",
	// 		entry.Caller.Function,
	// 		entry.Caller.Line,
	// 	)
	// 	if f.isLocal {
	// 		_, _ = fmt.Fprintf(b, formatColor, colorGray)
	// 	}
	// 	b.WriteString("]")
	// }

	// Log message
	if entry.Message != "" {
		b.WriteString("[")
		b.WriteString(entry.Message)
		b.WriteString("]")
	}

	keys := make([]string, 0, len(entry.Data))

	for key := range entry.Data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		json, _ := json.Marshal(entry.Data[key])
		b.WriteString("[")
		b.WriteString(key)
		_, _ = fmt.Fprintf(b, formatColor, colorGray)
		b.WriteString(":")
		b.WriteString(string(json))

		_, _ = fmt.Fprintf(b, formatColor, levelColor)
		b.WriteString("]")
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

// Logger instance of logger
var Logger logger

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FormatDate format date
func FormatDate(time time.Time) string {
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		time.Year(), int(time.Month()), time.Day(),
		time.Hour(), time.Minute(), time.Second())

}

// InitLogger create logger
func InitLogger(env string, name string) {
	var (
		l       = logrus.New()
		file    *os.File
		isLocal = env == "local"
		// isLocal = false
		today = converter.ConvertDateToStringCustom(time.Now(), converter.DateLayoutSimple)
	)

	if ok, _ := pathExists(".logs"); !ok && !isLocal {
		err := os.Mkdir(".logs", os.ModePerm)
		if err != nil {
			fmt.Printf("Failed make directory logs: %v\n", err)
		}
	}

	if env == "production" {
		l.SetLevel(logrus.InfoLevel)
	} else {
		env = "development"
		l.SetLevel(logrus.DebugLevel)
	}

	if !isLocal {
		filename := fmt.Sprintf("%s_%s_%s.log", name, env, today)
		// separator := string(os.PathSeparator)
		path := fmt.Sprintf(".logs/%s", filename)

		if _, err := os.Stat(path); err != nil {
			file, _ = os.Create(path)
		} else {
			file, _ = os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
		}
	}

	l.SetFormatter(&formatter{env, isLocal})

	if !isLocal {
		l.SetOutput(io.MultiWriter(file, os.Stdout))
	}

	l.SetReportCaller(true)
	Logger = logger{
		Logger:  l,
		enabled: false,
		Env:     env,
		Name:    name,
		Date:    today,
	}
}
