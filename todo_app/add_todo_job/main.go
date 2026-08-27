package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type newTodoRequest struct {
	Text string `json:"text"`
}

func main() {
	todoBackendUrl := os.Getenv("TODO_BACKEND_URL")
	if todoBackendUrl == "" {
		panic("TODO_BACKEND_URL is not set")
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", "https://en.wikipedia.org/wiki/Special:Random", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "add_todo_job/1.0 (helsinki_kube_course todo app cronjob)")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	if location == "" {
		panic("no Location header in response")
	}
	locationUrl, err := url.Parse(location)
	if err != nil {
		panic(err)
	}
	articleUrl := req.URL.ResolveReference(locationUrl).String()

	body, err := json.Marshal(newTodoRequest{Text: fmt.Sprintf("Read %s", articleUrl)})
	if err != nil {
		panic(err)
	}

	postResp, err := http.Post(todoBackendUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	defer postResp.Body.Close()

	fmt.Printf("created todo: %s (status %d)\n", articleUrl, postResp.StatusCode)
}
