package logs

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/Meduzz/gloegg/common"
	"github.com/Meduzz/gloegg/logging"
	"github.com/nats-io/nats.go"
)

type (
	Factory struct {
		nc    *nats.Conn
		topic string
	}
)

func NewLogFactory(nc *nats.Conn, project string) logging.HandlerFactory {
	topic := fmt.Sprintf("logs.%s.logs", project)
	return &Factory{nc, topic}
}

func (f *Factory) Spawn(level slog.Leveler, tags common.Tags) slog.Handler {
	r, w := io.Pipe()

	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
	})

	go f.listen(r)

	return handler.WithAttrs(tags.ToSlog())
}

func (f *Factory) listen(reader io.Reader) {
	decoder := json.NewDecoder(reader)

	for decoder.More() {
		record := json.RawMessage{}
		err := decoder.Decode(&record)

		if err == nil {
			msg := nats.NewMsg(f.topic)
			msg.Data = record
			f.nc.PublishMsg(msg)
		}

		// TODO safe to ignore this error?
	}
}
