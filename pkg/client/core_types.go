package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/technosophos/kubelite/v1"
)

func (c *Client) PodList(namespace string) (*v1.PodList, error) {

	resp, err := c.Get("/api/v1/namespaces/" + namespace + "/pods")
	if err != nil {
		err := fmt.Errorf("Error on GET request for %s: %v\n", c.Host, err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf("resp:  %v\n", resp)
	body, err := ioutil.ReadAll(resp.Body)
	//	log.Printf("body:  %s\n", body)

	pods := v1.PodList{}
	err = json.Unmarshal(body, &pods)
	if err != nil {
		return nil, err
	}

	return &pods, nil
}

func (c *Client) Pod(namespace string, podName string) (*v1.Pod, error) {
	resp, err := c.Get("/api/v1/namespaces/" + namespace + "/pods/" + podName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf("resp: %v\n", resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pod := v1.Pod{}
	err = json.Unmarshal(body, &pod)
	if err != nil {
		return nil, err
	}

	return &pod, nil
}
