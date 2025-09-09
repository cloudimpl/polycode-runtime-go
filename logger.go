package runtime

import (
	"encoding/json"
	"github.com/cloudimpl/byte-os/sdk"
	"time"
)

type LogLevel string

const (
	DebugLevel LogLevel = "DEBUG"
	InfoLevel  LogLevel = "INFO"
	WarnLevel  LogLevel = "WARN"
	ErrorLevel LogLevel = "ERROR"
)

type LogMsg struct {
	Level     LogLevel               `json:"level"`
	Section   string                 `json:"section"`
	Tags      map[string]interface{} `json:"tags"`
	Timestamp int64                  `json:"timestamp"`
	Message   string                 `json:"message"`
}

type LogEntry struct {
	msg *LogMsg
}

func (entry *LogEntry) Str(key string, val string) sdk.LogEntry {
	entry.msg.Tags[key] = val
	return entry
}

func (entry *LogEntry) Int64(key string, val int64) sdk.LogEntry {
	entry.msg.Tags[key] = val
	return entry
}

func (entry *LogEntry) Float64(key string, val float64) sdk.LogEntry {
	entry.msg.Tags[key] = val
	return entry
}

func (entry *LogEntry) Bool(key string, val bool) sdk.LogEntry {
	entry.msg.Tags[key] = val
	return entry
}

func (entry *LogEntry) Done() {
	entry.msg.Timestamp = time.Now().UnixMicro()
	logJson, err := json.Marshal(entry.msg)
	if err == nil {
		println(string(logJson))
	}
}

func (entry *LogEntry) Msg(msg string) {
	entry.msg.Message = msg
	entry.Done()
}

type JsonLogger struct {
	section string
}

func (logger *JsonLogger) log(level LogLevel) *LogEntry {
	return &LogEntry{
		msg: &LogMsg{
			Level:   level,
			Section: logger.section,
			Tags:    make(map[string]interface{}),
		},
	}
}

func (logger *JsonLogger) Debug() sdk.LogEntry {
	return logger.log(DebugLevel)
}

func (logger *JsonLogger) Info() sdk.LogEntry {
	return logger.log(InfoLevel)
}

func (logger *JsonLogger) Warn() sdk.LogEntry {
	return logger.log(WarnLevel)
}

func (logger *JsonLogger) Error() sdk.LogEntry {
	return logger.log(ErrorLevel)
}

func CreateLogger(section string) sdk.Logger {
	return &JsonLogger{
		section: section,
	}
}
