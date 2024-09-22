package main

import (
	"github.com/Meduzz/gloegg"
	. "github.com/Meduzz/gloegg-nats"
	"github.com/Meduzz/helper/nuts"
)

func main() {
	gloegg.AddMeta("asdf", "jklö")

	nc, _ := nuts.Connect()
	Setup(nc, "test", "test")

	logger := gloegg.CreateLogger("test")

	logger.Info("boom?", "boom", "no boom")
}
