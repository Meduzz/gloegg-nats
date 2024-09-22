package logs_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Meduzz/gloegg"
	"github.com/Meduzz/gloegg-nats/logs"
	"github.com/Meduzz/helper/nuts"
	"github.com/nats-io/nats.go"
)

const (
	LogLevel   = "level"
	LogMessage = "msg"
)

func TestLogs(t *testing.T) {
	nc, _ := nuts.Connect()
	fc := logs.NewLogFactory(nc, "test")
	gloegg.SetHandlerFactory(fc)
	feedback := make(chan map[string]any, 10)

	nc.Subscribe("logs.test.logs", func(msg *nats.Msg) {
		data := make(map[string]any)
		json.Unmarshal(msg.Data, &data)

		feedback <- data
	})

	subject := gloegg.CreateLogger("LoggerTest")

	t.Run("hello world!", func(t *testing.T) {
		subject.Info("Hello world!")
		event := <-feedback

		if event[LogLevel] != "INFO" {
			t.Errorf("log level was not INFO but %s", event[LogLevel])
		}

		if event[LogMessage] != "Hello world!" {
			t.Errorf("log message was not 'Hello world!' but %s", event[LogMessage])
		}
	})

	t.Run("the brief warning", func(t *testing.T) {
		subject.Warn("i have metadata", "x", "y")
		event := <-feedback

		if event[LogLevel] != "WARN" {
			t.Errorf("log level was not WARN but %s", event[LogLevel])
		}

		if event[LogMessage] != "i have metadata" {
			t.Errorf("log message was not 'i have metadata' but %s", event[LogMessage])
		}

		if event["x"] == nil {
			t.Error("there were no x metadata")
		}

		x := event["x"]
		xData, ok := x.(string)

		if !ok {
			t.Errorf("x was not string but %T", x)
		}

		if xData != "y" {
			t.Errorf("x was not 'y' but %s", xData)
		}
	})

	t.Run("there was an error", func(t *testing.T) {
		subject.Error("there's an error", "err", fmt.Errorf("im an error"))
		event := <-feedback

		if event[LogLevel] != "ERROR" {
			t.Errorf("log level was not ERROR but %s", event[LogLevel])
		}

		if event[LogMessage] != "there's an error" {
			t.Errorf("log message was not 'there's an error' but %s", event[LogMessage])
		}

		if event["err"] == nil {
			t.Error("there were no err metadata")
		}

		err := event["err"]
		errData, ok := err.(string)

		if !ok {
			t.Errorf("err was not string but %T", err)
		}

		if errData != "im an error" {
			t.Errorf("err was not 'im an error' but %s", errData)
		}
	})
}
