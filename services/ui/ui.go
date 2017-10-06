package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/dinever/golf"
)

const DATE_TZ = "America/Los_Angeles"

func newUiState() *UiState {
	return &UiState{
		LastUpdated:     "",
		RequestCount:    0,
		ErrorCount:      0,
		LastError:       "",
		UiHostname:      os.Getenv("HOSTNAME"),
		RecentResponses: make([]CounterResponse, 0),
		lock:            &sync.Mutex{},
	}
}

type UiState struct {
	LastUpdated     string
	RequestCount    int64
	ErrorCount      int64
	LastError       string
	UiHostname      string
	RecentResponses []CounterResponse

	lock *sync.Mutex
}

type CounterResponse struct {
	Hostname       string
	RequestCount   int64
	Version        string
	ChildStartTime string
}

type Counter struct {
	RequestCount  int64
	UniqueClients int64
	Hostname      string
	StartTimeNano int64
	Version       string
}

func newHttpClient() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   time.Second * 5,
			KeepAlive: 0,
		}).Dial,
	}
	return &http.Client{
		Transport: transport,
	}
}

func requestCounter(loc *time.Location) (CounterResponse, error) {
	url := "http://demo_counter:9000/"
	resp, err := newHttpClient().Get(url)
	if err != nil {
		return CounterResponse{}, err
	}
	defer resp.Body.Close()

	var counter Counter
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&counter)
	if err != nil {
		return CounterResponse{}, err
	}

	startTime := time.Unix(0, counter.StartTimeNano)

	return CounterResponse{
		Hostname:       counter.Hostname,
		RequestCount:   counter.RequestCount,
		Version:        counter.Version,
		ChildStartTime: startTime.In(loc).Format("Mon Jan 2 15:04:05"),
	}, nil
}

func (me *UiState) CounterFrag(ctx *golf.Context) {
	loc, err := time.LoadLocation(DATE_TZ)
	if err != nil {
		log.Printf("ERROR in LoadLocation for tz %s - %v\n", DATE_TZ, err)
		return
	}

	resp, err := requestCounter(loc)

	me.lock.Lock()
	defer me.lock.Unlock()

	me.RequestCount++
	me.LastUpdated = time.Now().In(loc).Format("Mon Jan 2 15:04:05")
	if err == nil {
		if len(me.RecentResponses) < 10 {
			me.RecentResponses = append(me.RecentResponses, resp)
		} else {
			me.RecentResponses = append(me.RecentResponses[1:], resp)
		}
	} else {
		log.Printf("ERROR %v\n", err)
		me.ErrorCount++
		me.LastError = me.LastUpdated
	}

	data := map[string]interface{}{
		"date":  time.Now().In(loc).Format("Mon Jan 2 15:04:05"),
		"state": me,
	}
	ctx.Loader("templates").Render("counter.html", data)
}

func homeHandler(ctx *golf.Context) {
	data := map[string]interface{}{}
	ctx.Loader("templates").Render("home.html", data)
}

func clockFragHandler(ctx *golf.Context) {
	var clockData string
	url := "http://demo_clock:9000/"
	resp, err := newHttpClient().Get(url)
	if err != nil {
		clockData = fmt.Sprintf("ERROR from clock service: %v", err)
	} else {
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		clockData = string(bodyBytes)
	}

	ctx.SetHeader("content-type", "text/plain")
	ctx.Send(clockData)
}

func main() {
	// turn off keepalive to see round-robin requests
	rand.Seed(time.Now().UTC().UnixNano())
	uiState := newUiState()
	app := golf.New()
	app.View.SetTemplateLoader("templates", "templates/")
	app.Get("/fragment/counter", uiState.CounterFrag)
	app.Get("/fragment/clock", clockFragHandler)
	app.Get("/", homeHandler)
	app.Run(":9000")
}
