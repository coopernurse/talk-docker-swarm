package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/dinever/golf"
)

const DATE_TZ = "America/Los_Angeles"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getOrSetToken(ctx *golf.Context) string {
	token, _ := ctx.Cookie("demotoken")
	if token == "" {
		token = randSeq(20)
		ctx.SetCookie("demotoken", token, 0)
	}
	return token
}

func newChildState() *ChildState {
	return &ChildState{
		LastUpdated:     "",
		RequestCount:    0,
		ErrorCount:      0,
		LastError:       "",
		RecentResponses: make([]ChildResponse, 0),
		lock:            &sync.Mutex{},
	}
}

type ChildState struct {
	LastUpdated     string
	RequestCount    int64
	ErrorCount      int64
	LastError       string
	RecentResponses []ChildResponse

	lock *sync.Mutex
}

type ChildResponse struct {
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

func requestChild(token string, loc *time.Location) (ChildResponse, error) {
	url := "http://demo_counter:9000/" + token
	resp, err := newHttpClient().Get(url)
	if err != nil {
		return ChildResponse{}, err
	}
	defer resp.Body.Close()

	var counter Counter
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&counter)
	if err != nil {
		return ChildResponse{}, err
	}

	startTime := time.Unix(0, counter.StartTimeNano)

	return ChildResponse{
		Hostname:       counter.Hostname,
		RequestCount:   counter.RequestCount,
		Version:        counter.Version,
		ChildStartTime: startTime.In(loc).Format("Mon Jan 2 15:04:05"),
	}, nil
}

func (me *ChildState) HandleFrag(ctx *golf.Context) {
	loc, err := time.LoadLocation(DATE_TZ)
	if err != nil {
		log.Printf("ERROR in LoadLocation for tz %s - %v\n", DATE_TZ, err)
		return
	}

	resp, err := requestChild(getOrSetToken(ctx), loc)

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
	childState := newChildState()
	app := golf.New()
	app.View.SetTemplateLoader("templates", "templates/")
	app.Get("/fragment/counter", childState.HandleFrag)
	app.Get("/fragment/clock", clockFragHandler)
	app.Get("/", homeHandler)
	app.Run(":9000")
}
