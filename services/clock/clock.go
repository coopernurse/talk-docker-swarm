package main

import (
	"os"
	"time"

	"github.com/dinever/golf"
)

const VERSION = "1"

func clock(ctx *golf.Context) {
	ctx.SetHeader("content-type", "text/plain")
	ctx.Send(time.Now().Format("15:04:05 MST") + "\n" +
		os.Getenv("HOSTNAME") + "\n" +
		VERSION)
}

func main() {
	app := golf.New()
	app.Get("/", clock)
	app.Run(":9000")
}
