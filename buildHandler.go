package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/mchudgins/jeeves/pkg/k8sClient"
)

type githubRepository struct {
	CloneURL string `json:"clone_url"`
}
type githubBuildRequest struct {
	Ref             string           `json:"ref"`
	Repo            githubRepository `json:"repository"`
	GithubPushEvent string           `json:"githubPushEvent"`
}

var buildRegex *regexp.Regexp

func init() {
	buildRegex = regexp.MustCompile("/builds/(?P<image>[a-zA-Z:]*)/(?P<buildName>[[:alpha:]]*)/*(?P<excess>.*){0,}")
}

func launchBuildPod(imageName string, buildName string, body []byte) (string, error) {

	log.Printf("build data: %s", body)

	evt := githubBuildRequest{}

	err := json.Unmarshal(body, &evt)
	if err != nil {
		log.Fatalf("unable to Unmarshal %s: %v", body, err)
		return "", err
	}

	return "", nil
}

func getBuildInfo(k8sClient *client.Client,
	dao *DaoBuilds,
	imageName string,
	buildName string,
	buildNum int) (string, error) {

	log.Printf("getBuildInfo")
	ns := k8sClient.CurrentNamespace

	b, err := dao.Fetch(ns, buildName, buildNum)
	if err != nil {
		return "", err
	}
	log.Printf("build:  %v", b)
	return "", nil
}

func NewBuildHandler(k8sClient *client.Client, dao *DaoBuilds) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("processing http request: %v\n", *r)

		url := r.URL.RequestURI()

		elements := buildRegex.FindStringSubmatch(url)
		log.Printf("elements: %q", elements)
		if len(elements) != 4 {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		imageName := elements[1]
		buildName := elements[2]
		buildNum := elements[3]

		var newBuildURL bool

		log.Printf("imageName: %s; buildName: %s", imageName, buildName)
		if len(buildNum) > 0 {
			log.Printf("buildNum: %s", buildNum)
		} else {
			newBuildURL = true
		}

		switch strings.ToUpper(r.Method) {
		case "POST":
			if !newBuildURL {
				http.Error(w, http.StatusText(400), 400)
				return
			}
			if r.Body != nil {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					log.Fatal(err)
				}
				defer r.Body.Close()

				if body != nil && len(body) > 0 {
					newBuild, err := launchBuildPod(imageName, buildName, body)
					if err != nil {
						log.Fatal(err)
						http.Error(w, http.StatusText(400), 400)
						return
					}
					log.Printf("New build:  %s", newBuild)
					return
				}
			}

			http.Error(w, http.StatusText(400), 400)

		case "GET":
			num, err := strconv.Atoi(buildNum)
			if err == nil {
				build, err := getBuildInfo(k8sClient, dao, imageName, buildName, num)
				if err != nil {
					switch err.(type) {
					case ErrorCode:
						if err.(ErrorCode).Code == 404 {
							http.Error(w, http.StatusText(404), 404)
						} else {
							log.Fatalf("While 'Fetch'ing %s/%s/%d from dynamodDB: %v",
								imageName, buildName, num, err)
							http.Error(w, http.StatusText(500), 500)
						}

					case error:
						http.Error(w, http.StatusText(500), 500)

					default:
						http.Error(w, http.StatusText(500), 500)
					}
					return
				}
				log.Printf("Build: %v", build)
				return
			}
			http.Error(w, http.StatusText(404), 404)
			return

		default:
			fmt.Fprintf(w, "%s : %q", r.Method, html.EscapeString(r.URL.Path))
		}
	}
}
