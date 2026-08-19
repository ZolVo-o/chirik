package realtime

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type Publisher interface {
	Publish(Event)
}

type HTTPPublisher struct {
	url    string
	secret string
	client *http.Client
}

func NewHTTPPublisher(url, secret string) *HTTPPublisher {
	return &HTTPPublisher{
		url:    strings.TrimRight(url, "/"),
		secret: secret,
		client: &http.Client{Timeout: 2 * time.Second},
	}
}

func (p *HTTPPublisher) Publish(event Event) {
	if p.url == "" || p.secret == "" {
		return
	}
	body, err := json.Marshal(event)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, p.url+"/publish", bytes.NewReader(body))
	if err != nil {
		return
	}
	go p.send(req)
}

func (p *HTTPPublisher) send(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-Secret", p.secret)
	resp, err := p.client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
}
