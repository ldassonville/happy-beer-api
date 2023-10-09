package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type JsonLogFormatter struct {
	Env string
}

func NewLogFormatter(env string) *JsonLogFormatter {
	return &JsonLogFormatter{
		Env: env,
	}
}

func (f *JsonLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	hostName, _ := os.Hostname()
	fmtTimestamp := entry.Time.Format(time.RFC3339)

	entry.Data["env"] = f.Env
	entry.Data["log_timestamp"] = fmtTimestamp
	entry.Data["level"] = entry.Level.String()
	entry.Data["message"] = entry.Message
	entry.Data["hostname"] = hostName

	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
