package client

import (
	"crypto/tls"
	"log"
	"net/http"
)

type Client struct {
	*http.Client
	Host     string
	Username string
	Password string
}

func NewClientOrDie(host string, username string, password string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		Client:   &http.Client{Transport: tr},
		Host:     host,
		Username: username,
		Password: password,
	}
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	log.Printf("Get'ting:  %s\n", c.Host+url)
	req, err := http.NewRequest("GET", c.Host+url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	//	req.Header.Set("Authorization", "Bearer ")
	return c.Do(req)
}
