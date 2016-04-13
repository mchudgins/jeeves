package main

/*
https://github.com/mchudgins/kubelite

https://dev.dstcorp.io:8443/swaggerapi/api/v1

curl -ik -H 'authorization: Bearer ****'  https://dev.dstcorp.io:8443/api/v1/namespaces/mch-dev0/pods

*/

import (
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/mchudgins/jeeves/pkg/k8sClient"
)

var awsRegion = flag.String("region", "us-east-1", "AWS region")
var addr = flag.String("apiserver", "", "k8s server ip address (https://192.168.1.1)")
var user = flag.String("username", "", "apiserver username")
var pword = flag.String("password", "", "apiserver password")

func main() {
	flag.Parse()
	fmt.Println("Hello, world.")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("processing htp request: %v\n", *r)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/builds", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("processing htp request: %v\n", *r)

		switch strings.ToUpper(r.Method) {
		case "POST":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)

			}
			defer r.Body.Close()

			log.Printf("read: %s\n\n", body)
		default:
			fmt.Fprintf(w, "%s : %q", r.Method, html.EscapeString(r.URL.Path))
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Ah, my health is just fine.  thanks.")
	})

	out, err := exec.Command("go", "version").Output()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("output: %s\n", out)

	c := client.NewClientOrDie()

	pods, err := c.PodList("mch-dev0")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	for _, p := range pods.Items {
		log.Printf("%s: %s %v\n\n", p.Name, p.Status, p)
	}

	pod, err := c.Pod("mch-dev0", "jumpbox")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Printf("pod:  %v\n", *pod)

	//	client := k8sClientFactory()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
