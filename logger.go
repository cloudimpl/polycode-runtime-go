package runtime

import (
	"encoding/json"
	"github.com/cloudimpl/polycode-sdk-go"
	"time"
)

const (
	DebugLevel polycode.LogLevel = "DEBUG"
	InfoLevel  polycode.LogLevel = "INFO"
	WarnLevel  polycode.LogLevel = "WARN"
	ErrorLevel polycode.LogLevel = "ERROR"
)

type LogMsg struct {
	Level     polycode.LogLevel      `json:"level"`
	Section   string                 `json:"section"`
	Tags      map[string]interface{} `json:"tags"`
	Timestamp int64                  `json:"timestamp"`
	Message   string                 `json:"message"`
}

type LogEntry struct {
	msg *LogMsg
}

func (entry *LogEntry) Str(key string, val string) polycode.LogEntry {
	entry.msg.Tags[key] = val
	return entry
}

func (entry *LogEntry) Int64(key string, val int64) polycode.LogEntry {
	entry.msg.Tags[key] = val
	return entry
}

func (entry *LogEntry) Float64(key string, val float64) polycode.LogEntry {
	entry.msg.Tags[key] = val
	return entry
}

func (entry *LogEntry) Bool(key string, val bool) polycode.LogEntry {
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

func (logger JsonLogger) Debug() polycode.LogEntry {
	return logger.Log(DebugLevel)
}

func (logger JsonLogger) Info() polycode.LogEntry {
	return logger.Log(InfoLevel)
}

func (logger JsonLogger) Warn() polycode.LogEntry {
	return logger.Log(WarnLevel)
}

func (logger JsonLogger) Error() polycode.LogEntry {
	return logger.Log(ErrorLevel)
}

func (logger JsonLogger) Log(level polycode.LogLevel) polycode.LogEntry {
	return &LogEntry{
		msg: &LogMsg{
			Level:   level,
			Section: logger.section,
			Tags:    make(map[string]interface{}),
		},
	}
}
