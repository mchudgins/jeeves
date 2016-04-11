package main

import (
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

type k8sClient struct {
	*unversioned.Client
}

var awsRegion = flag.String("region", "us-east-1", "AWS region")
var addr = flag.String("apiserver", "", "k8s server ip address (https://192.168.1.1)")
var user = flag.String("username", "", "apiserver username")
var pword = flag.String("password", "", "apiserver password")

func k8sClientFactory() *k8sClient {
	log.Printf("Host:  %v; UserName: %v; Password: %v\n", *addr, *user, *pword)
	if len(*addr) > 0 && len(*user) > 0 && len(*pword) > 0 {
		config := restclient.Config{
			Host:     *addr,
			Username: *user,
			Password: *pword,
			Insecure: true,
		}
		return &k8sClient{unversioned.NewOrDie(&config)}
	} else {
		kubernetesService := os.Getenv("KUBERNETES_SERVICE_HOST")
		if kubernetesService == "" {
			glog.Fatalf("Please specify the Kubernetes server with --server")
		}
		apiServer := fmt.Sprintf("https://%s:%s", kubernetesService, os.Getenv("KUBERNETES_SERVICE_PORT"))

		token, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
		if err != nil {
			glog.Fatalf("No service account token found")
		}

		config := restclient.Config{
			Host:        apiServer,
			BearerToken: string(token),
			Insecure:    true,
		}

		c, err := unversioned.New(&config)
		if err != nil {
			glog.Fatalf("Failed to make client: %v", err)
		}
		return &k8sClient{c}
	}
}

func main() {
	flag.Parse()
	fmt.Println("Hello, world.")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("processing htp request.")
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	out, err := exec.Command("go", "version").Output()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("output: %s\n", out)

	client := k8sClientFactory()
	log.Printf("%v\n", client.APIVersion())
	v, err := client.ServerVersion()
	log.Printf("server info: %v\n", v)
	p, err := client.Pods("mch-dev0").List(api.ListOptions{})
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("pods: %v\n", p)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
