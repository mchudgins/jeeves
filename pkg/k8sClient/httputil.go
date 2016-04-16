package k8sClient

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/ghodss/yaml"
)

type Client struct {
	*http.Client
	Host             string
	Username         string
	Token            string
	CurrentNamespace string
}

type k8sClusterDetails struct {
	Insecure string `json:"insecure-skip-tls-verify"`
	Server   string `json:"server"`
}

type k8sCluster struct {
	Name    string            `json:"name"`
	Details k8sClusterDetails `json:"cluster"`
}

type k8sContextDetails struct {
	Cluster   string `json:"cluster"`
	Namespace string `json:"namespace"`
	User      string `json:"user"`
}

type k8sContext struct {
	Name    string            `json:"name"`
	Context k8sContextDetails `json:"context"`
}

type k8sUserDetails struct {
	Token string `json:"token"`
}

type k8sUser struct {
	Name string         `json:"name"`
	User k8sUserDetails `json:"user"`
}

type k8sConfig struct {
	Kind           string       `json:"kind"`
	ApiVersion     string       `json:"apiVersion"`
	Clusters       []k8sCluster `json:"clusters"`
	Contexts       []k8sContext `json:"contexts"`
	CurrentContext string       `json:"current-context"`
	Users          []k8sUser    `json:"users"`
}

func (c *k8sConfig) ActiveContext() (*k8sContext, error) {
	if len(c.CurrentContext) == 0 {
		return nil, fmt.Errorf("no active context")
	}

	return c.FindContext(c.CurrentContext)
}

func (c *k8sConfig) FindContext(name string) (*k8sContext, error) {
	for _, ctx := range c.Contexts {
		if name == ctx.Name {
			return &ctx, nil
		}
	}
	return nil, fmt.Errorf("unable to find context %s", name)
}

func (c *k8sConfig) FindCluster(name string) (*k8sCluster, error) {
	for _, cluster := range c.Clusters {
		if name == cluster.Name {
			return &cluster, nil
		}
	}
	return nil, fmt.Errorf("unable to find cluster %s", name)
}

func (c *k8sConfig) FindUser(name string) (*k8sUser, error) {
	for _, user := range c.Users {
		if name == user.Name {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("unable to find cluster %s", name)
}

func loadKubeConfig() *k8sConfig {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// load the .kube/config file
	log.Printf("home:  %s", user.HomeDir)
	cfg, err := os.Open(user.HomeDir + "/.kube/config")
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	info, err := cfg.Stat()
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}

	buf := make([]byte, info.Size(), info.Size())
	_, err = cfg.Read(buf)
	if err != nil {
		log.Fatal(err)
		os.Exit(4)
	}

	k8sConfig := k8sConfig{}
	err = yaml.Unmarshal(buf, &k8sConfig)
	if err != nil {
		log.Fatal(err)
		os.Exit(5)
	}

	return &k8sConfig
}

/*
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
*/

func NewClientOrDie() *Client {
	k8sConfig := loadKubeConfig()
	ctx, _ := k8sConfig.ActiveContext()
	cluster, _ := k8sConfig.FindCluster(ctx.Context.Cluster)
	user, _ := k8sConfig.FindUser(ctx.Context.User)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Client{
		Client:           &http.Client{Transport: tr},
		Host:             cluster.Details.Server,
		Username:         ctx.Context.User,
		Token:            user.User.Token,
		CurrentNamespace: ctx.Context.Namespace,
	}
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	log.Printf("GET'ting:  %s\n", c.Host+url)
	req, err := http.NewRequest("GET", c.Host+url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err = c.Do(req)

	if resp.StatusCode == 401 {
		log.Fatal("Unauthorized (are your login credentials current?)")
	}

	return resp, err
}

func (c *Client) Post(url string, body []byte) (resp *http.Response, err error) {
	log.Printf("POST'ting:  %s\n", c.Host+url)
	var reader io.Reader
	reader, err = gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequest("POST", c.Host+url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err = c.Do(req)

	if resp.StatusCode == 401 {
		log.Fatal("Unauthorized (are your login credentials current?)")
	}

	return resp, err
}
