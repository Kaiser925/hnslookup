package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// DNSResponse represents the response of dns query request.
// More information see here:
// https://developers.cloudflare.com/1.1.1.1/dns-over-https
type DNSResponse struct {
	Status   int         `json:"Status"`
	TC       bool        `json:"TC"`
	RD       bool        `json:"RD"`
	RA       bool        `json:"RA"`
	AD       bool        `json:"AD"`
	CD       bool        `json:"CD"`
	Question []*Question `json:"Question"`
	Answer   []*Answer   `json:"Answer"`
}

// Question represents a DNS question.
type Question struct {
	// The record name requested.
	// e.g: example.com
	Name string `url:"name" json:"name"`

	// The type of DNS record requested. Either a numeric value or text,
	// e.g: "AAAA"
	Type int `url:"type" json:"type"`
}

// Answer represents a DNS Answer.
type Answer struct {
	// The record owner.
	Name string `json:"name"`

	// The type of DNS record.
	Type int `json:"type"`

	// The number of seconds the answer can be stored in cache
	// before it is considered stale.
	TTL int `json:"TTL"`

	// The value of the DNS record for the given name and type.
	// The data will be in text for standardize recode types and
	// in hex for unknown types.
	Data string `json:"data"`
}

func MapTypeToInt(typeS string) int {
	if len(typeS) == 0 {
		return 0
	}
	return 0 // TODO: handle mapping
}

func MapTypeToString(typeI int) string {
	return "AAAA" // TODO:
}

func fmtAnswer(answer *Answer) string {
	return fmt.Sprintf("Name: %s\nType: %s\nTL : %d\nData: %s\n", answer.Name, MapTypeToString(answer.Type), answer.TTL, answer.Data)
}

func sendRequest(questions *Question) ([]byte, error) {
	val, err := query.Values(questions)
	if err != nil {
		return nil, err
	}

	targetUrl, err := url.Parse("https://cloudflare-dns.com/dns-query?" + val.Encode())
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	header.Add("accept", "application/dns-json")

	req := &http.Request{
		Method: "GET",
		URL:    targetUrl,
		Header: header,
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if !strings.Contains(resp.Status, "OK") {
		return nil, errors.New("HTTP request error: " + resp.Status + " " + targetUrl.String())
	}

	return ioutil.ReadAll(resp.Body)
}

func printAnswer(answers [][]*Answer) {
	for _, anss := range answers {
		for _, ans := range anss {
			fmt.Println(fmtAnswer(ans))
		}
	}
}

// HandleQuestions handles DNS questions.
func HandleQuestions(questions []*Question) error {
	answers := make([][]*Answer, len(questions))
	for i, question := range questions {
		data, err := sendRequest(question)
		if err != nil {
			return err
		}
		resp := new(DNSResponse)

		if err := json.Unmarshal(data, resp); err != nil {
			return err
		}
		answers[i] = resp.Answer
	}

	printAnswer(answers)
	return nil
}
