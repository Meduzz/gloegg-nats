package flags_test

import (
	"encoding/json"
	"testing"

	"github.com/Meduzz/gloegg-nats/api"
	"github.com/Meduzz/gloegg-nats/flags"
	"github.com/Meduzz/gloegg/toggles"
	"github.com/Meduzz/helper/nuts"
	"github.com/nats-io/nats.go"
)

func TestFlags(t *testing.T) {
	nc, _ := nuts.Connect()

	feedback := make(chan *api.FlagEvent, 10)

	// listen to request for flags
	nc.Subscribe("logs.test.request", func(msg *nats.Msg) {
		reply := msg.Reply
		flag := &api.FlagEvent{}
		flag.Kind = toggles.KindString
		flag.Name = "request"
		flag.Value = "test"

		bs, _ := json.Marshal(flag)
		nc.Publish(reply, bs)
	})

	// listen to local flag updates
	nc.Subscribe("logs.test.local", func(msg *nats.Msg) {
		flag := &api.FlagEvent{}
		json.Unmarshal(msg.Data, flag)

		feedback <- flag
	})

	handler := flags.NewFlagHandler(nc, "test")
	handler.Setup()

	t.Run("verify requested flags", func(t *testing.T) {
		<-feedback // await the ping pong to settle
		subject := toggles.GetStringToggle("request")

		if !subject.Equals("test") {
			t.Errorf("request toggle value was not 'test' but %s", subject.Value())
		}
	})

	t.Run("update local flag", func(t *testing.T) {
		toggles.SetStringToggle("update", "test")

		result := <-feedback

		if result.Kind != toggles.KindString && result.Name != "update" {
			t.Errorf("kind and name was incorrect, was kind:%s name:%s", result.Kind, result.Name)
		}

		if result.Value != "test" {
			t.Errorf("value was not 'test' but %s", result.Value)
		}
	})

	t.Run("update remote flag", func(t *testing.T) {
		flag := &api.FlagEvent{}
		flag.Kind = toggles.KindString
		flag.Name = "remote"
		flag.Value = "test"

		bs, _ := json.Marshal(flag)
		nc.Publish("logs.test.remote", bs)

		<-feedback // wait for the ping pong to settle
		subject := toggles.GetStringToggle("remote")

		if !subject.Equals("test") {
			t.Errorf("remote toggle value was not 'test' but %s", subject.Value())
		}
	})
}
