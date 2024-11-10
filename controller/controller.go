package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type Request struct {
	IdMap map[string]bool
	mu    sync.Mutex
}

var req = &Request{
	IdMap: make(map[string]bool),
	mu:    sync.Mutex{},
}

func CountLogger() {
	log.Printf("Unique controller requests in the last minute: %d", getCount())
	// re-initialise it for the next minute after logging
	resetCount()
}

func getCount() int {
	req.mu.Lock()
	count := len(req.IdMap)
	req.mu.Unlock()
	return count
}

func resetCount() {
	req.mu.Lock()
	req.IdMap = make(map[string]bool)
	req.mu.Unlock()
}

func AddCount(id string) {
	req.mu.Lock()
	req.IdMap[id] = true
	req.mu.Unlock()
}

type Payload struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func SendGETRequest(endpoint string) []byte {
	url := fmt.Sprintf("%s?count=%d", endpoint, getCount())
	resp, err := http.Get(url)
	payLoad := Payload{}
	if err != nil {
		payLoad.StatusCode = http.StatusBadRequest
		response, err := json.Marshal(payLoad)
		if err != nil {
			fmt.Printf("Failed to Marshal: %v\n", payLoad)
			return []byte{}
		}
		return response
	}
	defer resp.Body.Close()

	payLoad.StatusCode = resp.StatusCode
	if buf, err := io.ReadAll(resp.Body); err == nil {
		payLoad.Message = string(buf)
	}
	response, err := json.Marshal(payLoad)
	if err != nil {
		fmt.Printf("Failed to Marshal: %v\n", payLoad)
		return []byte{}
	}
	return response
}

func SendPOSTRequest(url string) []byte {
	count := fmt.Sprintf("%d", getCount())
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(count)))
	payLoad := Payload{}

	if err != nil {
		payLoad.StatusCode = http.StatusBadRequest
		response, err := json.Marshal(payLoad)
		if err != nil {
			return []byte{}
		}
		return response
	}
	defer resp.Body.Close()

	payLoad.StatusCode = resp.StatusCode
	if buf, err := io.ReadAll(resp.Body); err == nil {
		payLoad.Message = string(buf)
	}
	response, err := json.Marshal(payLoad)
	if err != nil {
		return []byte{}
	}
	return response
}
