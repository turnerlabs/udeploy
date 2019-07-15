package github

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"encoding/json"
)

// Tag ...
type Tag struct {
	Name   string
	Commit TagCommit
}

// TagCommit ...
type TagCommit struct {
	SHA string
}

// GetTags ...
func GetTags(org, repo, url, accessToken string) (map[string]Tag, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repos/%s/%s/tags", url, org, repo), nil)
	if err != nil {
		return map[string]Tag{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", accessToken))
	resp, err := client.Do(req)
	if err != nil {
		return map[string]Tag{}, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return map[string]Tag{}, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return map[string]Tag{}, err
	}

	tags := []Tag{}
	if err := json.Unmarshal(body, &tags); err != nil {
		return map[string]Tag{}, err
	}

	tagMap := map[string]Tag{}
	for _, tag := range tags {
		tagMap[tag.Name] = tag
	}

	return tagMap, nil
}

// Commit ...
type Commit struct {
	SHA     string
	Parents []CommitParent
	Commit  CommitCommit
}

// CommitParent ...
type CommitParent struct {
	SHA string
}

// CommitCommit ...
type CommitCommit struct {
	Message string
}

// GetCommit ...
func GetCommit(org, repo, SHA, url, accessToken string) (Commit, error) {
	client := &http.Client{}

	commit := Commit{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repos/%s/%s/commits/%s", url, org, repo, SHA), nil)
	if err != nil {
		return commit, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", accessToken))
	resp, err := client.Do(req)
	if err != nil {
		return commit, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return commit, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return commit, err
	}

	err = json.Unmarshal(body, &commit)

	return commit, err
}
