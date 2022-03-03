package log

import (
	"testing"
	"time"
)

func TestLogInstance(t *testing.T) {
	newLog := NewLogger()

	logConf := LogConfig{
		Level: "trace",
		FW: ConfFileWirter{
			On:              true,
			LogPath:         "./log_test.log",
			RotateLogPath:   "./log_test.log",
			WfLogPath:       "./log_test.wf.log",
			RotateWfLogPath: "./log_test.wf.log",
		},
		CW: ConfConsoleWritet{
			On:    true,
			color: true,
		},
	}

	SetupLogInstanceWithConf(logConf, newLog)

	newLog.Info("test message")

	newLog.Close()

	time.Sleep(time.Second)
}
