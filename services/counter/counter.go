package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/dinever/golf"
)

const VERSION = "1.18"

func newCounter() *Counter {
	return &Counter{
		RequestCount:  0,
		UniqueClients: 0,
		Hostname:      os.Getenv("HOSTNAME"),
		Version:       VERSION,
		StartTimeNano: time.Now().UnixNano(),
		clientTokens:  make(map[string]bool),
		lock:          &sync.Mutex{},
	}
}

type Counter struct {
	RequestCount  int64
	UniqueClients int64
	Hostname      string
	StartTimeNano int64
	Version       string

	clientTokens map[string]bool
	lock         *sync.Mutex
}

func (me *Counter) addReqAndSerialize(token string) ([]byte, error) {
	me.lock.Lock()
	defer me.lock.Unlock()

	me.RequestCount++
	_, ok := me.clientTokens[token]
	if !ok {
		me.clientTokens[token] = true
		me.UniqueClients++
	}

	return json.Marshal(me)
}

func (me *Counter) Handle(ctx *golf.Context) {
	data, err := me.addReqAndSerialize(ctx.Param("token"))
	if err == nil {
		ctx.SetHeader("content-type", "application/json")
		ctx.Send(data)
	} else {
		log.Printf("ERROR in addReqAndSerialize: %v\n", err)
		ctx.Abort(500)
	}
}

func dumpEnv(ctx *golf.Context) {
	b := bytes.NewBufferString("")
	for _, e := range os.Environ() {
		b.WriteString(e)
		b.WriteString("\n")
	}
	ctx.SetHeader("content-type", "text/plain")
	ctx.Send(b.String())
}

func main() {
	counter := newCounter()

	app := golf.New()
	app.Get("/env", dumpEnv)
	app.Get("/:token", counter.Handle)
	app.Run(":9000")
}
