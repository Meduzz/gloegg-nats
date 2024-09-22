package flags

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Meduzz/gloegg"
	"github.com/Meduzz/gloegg-nats/api"
	"github.com/Meduzz/gloegg/toggles"
	"github.com/nats-io/nats.go"
)

type (
	FlagHandler struct {
		nc         *nats.Conn
		project    string
		localTopic string
		logger     *slog.Logger
	}
)

func NewFlagHandler(nc *nats.Conn, project string) *FlagHandler {
	logger := gloegg.CreateLogger("FlagSync")
	localTopic := fmt.Sprintf("logs.%s.local", project)

	return &FlagHandler{nc, project, localTopic, logger}
}

func (f *FlagHandler) Setup() {
	// listen to local updates
	toggles.Subscribe(f.localFlagUpdated)

	// listen to remote updates
	remoteTopic := fmt.Sprintf("logs.%s.remote", f.project)
	f.nc.Subscribe(remoteTopic, f.remoteFlagEventHandler)

	// request existing flags from server
	msg := nats.NewMsg(fmt.Sprintf("logs.%s.request", f.project))
	reply := fmt.Sprintf("logs.%s.remote", f.project)
	msg.Reply = reply
	f.nc.PublishMsg(msg)
}

func (f *FlagHandler) localFlagUpdated(flag *toggles.UpdatedToggle) {
	event := &api.FlagEvent{}
	event.Name = flag.Name
	event.Kind = flag.Kind
	event.Value = flag.Value

	bs, err := json.Marshal(event)

	if err == nil {
		f.nc.Publish(f.localTopic, bs)
	}
}

func (f *FlagHandler) remoteFlagUpdated(flag *api.FlagEvent) {
	switch flag.Kind {
	case toggles.KindString:
		value, ok := flag.Value.(string)

		if ok {
			toggles.SetStringToggle(flag.Name, value)
		}
	case toggles.KindInt:
		value, ok := flag.Value.(int)

		if ok {
			toggles.SetIntToggle(flag.Name, value)
		}
	case toggles.KindInt64:
		value, ok := flag.Value.(int64)

		if ok {
			toggles.SetInt64Toggle(flag.Name, value)
		}
	case toggles.KindBool:
		value, ok := flag.Value.(bool)

		if ok {
			toggles.SetBoolToggle(flag.Name, value)
		}
	case toggles.KindObject:
		value, ok := flag.Value.(map[string]any)

		if ok {
			toggles.SetObjectToggle(flag.Name, value)
		}
	case toggles.KindFloat32:
		value, ok := flag.Value.(float32)

		if ok {
			toggles.SetFloat32Toggle(flag.Name, value)
		}
	case toggles.KindFloat64:
		value, ok := flag.Value.(float64)

		if ok {
			toggles.SetFloat64Toggle(flag.Name, value)
		}
	default:
	}
}

func (f *FlagHandler) remoteFlagEventHandler(msg *nats.Msg) {
	event := &api.FlagEvent{}
	err := json.Unmarshal(msg.Data, event)

	if err != nil {
		f.logger.Error("parsing json threw error", "error", err.Error())
	} else {
		f.remoteFlagUpdated(event)
	}
}
