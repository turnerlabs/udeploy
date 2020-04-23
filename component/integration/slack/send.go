package slack

import (
	"bytes"
	"net/http"
)

func Send(msg, webHookURL string) error {
	req, err := http.NewRequest("POST", webHookURL, bytes.NewBuffer([]byte(msg)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
