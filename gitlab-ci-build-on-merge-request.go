package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

type requestBody struct {
	ObjectKind string `json:"object_kind"` // merge_request
	Project    struct {
		Name string `json:"name"`
	} `json:"project"`
	ObjectAttributes struct {
		SourceBranch    string `json:"source_branch"`
		SourceProjectId int    `json:"source_project_id"`
		State           string `json:"state"` // merged, opened or closed
		LastCommit      struct {
			Id string `json:"id"`
		} `json:"last_commit"`
		WorkInProgress bool `json:"work_in_progress"`
	} `json:"object_attributes"`
}

type trigger struct {
	Token string `json:"token"`
}

func printUsageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg+"\n\n")
	}
	flag.Usage()
	os.Exit(1)
}

// todo: make sure private token does not leak through response/log
func main() {
	var baseURL = flag.String("url", "", "URL (e.g. http://gitlab.com)")
	var privateToken = flag.String("private_token", "", "Authorization Token (e.g. XXxXXx0xxxXXXxXxXxxX)")
	var port = flag.Int("port", 8080, "Port")
	flag.Parse()
	if *baseURL == "" {
		printUsageAndExit("Error: --url is required")
	} else if *privateToken == "" {
		printUsageAndExit("Error: --private_token is required")
	}
	http.HandleFunc("/hook", func(w http.ResponseWriter, r *http.Request) {
		var requestBody = &requestBody{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			log.Printf("WARN: Failed to deserialize request body (%s)", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		requestBodyAsByteArray, _ := json.Marshal(requestBody)
		log.Printf("INFO: Received %s", string(requestBodyAsByteArray))
		// do not trigger build if merge request is WIP or merged/closed
		if requestBody.ObjectKind != "merge_request" || requestBody.ObjectAttributes.State != "opened" ||
			requestBody.ObjectAttributes.WorkInProgress {
			return
		}
		// do not trigger if build for commit was already triggered
		buildsUrl := fmt.Sprintf(
			"%s/api/v3/projects/%d/repository/commits/%s/builds?private_token=%s",
			*baseURL,
			requestBody.ObjectAttributes.SourceProjectId,
			requestBody.ObjectAttributes.LastCommit.Id,
			*privateToken)
		buildsRes, err := http.Get(buildsUrl)
		if err != nil {
			log.Printf("WARN: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer buildsRes.Body.Close()
		var builds []map[string]interface{}
		if err := json.NewDecoder(buildsRes.Body).Decode(&builds); err != nil {
			log.Printf("WARN: Failed to deserialize response of GET %s (%s)", buildsUrl, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(builds) > 0 {
			return
		}
		trigger, err := resolveTrigger(*baseURL, *privateToken, requestBody.ObjectAttributes.SourceProjectId)
		if err != nil {
			log.Printf("WARN: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		triggerUrl := fmt.Sprintf(
			"%s/api/v3/projects/%d/trigger/builds?ref=%s&token=%s",
			*baseURL,
			requestBody.ObjectAttributes.SourceProjectId,
			requestBody.ObjectAttributes.SourceBranch,
			trigger.Token)
		triggerRes, err := http.PostForm(triggerUrl, url.Values{})
		if err != nil {
			log.Printf("WARN: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer triggerRes.Body.Close()
		if triggerRes.StatusCode != 201 {
			log.Printf("WARN: POST %s resulted in %d", triggerUrl, triggerRes.StatusCode)
			http.Error(w, fmt.Sprint("POST %s resulted in %d", triggerUrl, triggerRes.StatusCode),
				http.StatusInternalServerError)
			return
		}
		log.Printf("INFO: Triggered build of %s#%s", requestBody.Project.Name,
			requestBody.ObjectAttributes.SourceBranch)
	})
	log.Printf(fmt.Sprintf("INFO: Listening on port %d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func resolveTrigger(baseURL string, privateToken string, projectId int) (*trigger, error) {
	fullURL := fmt.Sprintf(
		"%s/api/v3/projects/%d/triggers?private_token=%s",
		baseURL,
		projectId,
		privateToken)
	res, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var triggers []trigger
	if err := json.NewDecoder(res.Body).Decode(&triggers); err != nil {
		return nil, fmt.Errorf("Failed to deserialize response of GET %s (%s)", fullURL, err.Error())
	}
	if len(triggers) == 0 {
		res, err := http.PostForm(fullURL, url.Values{})
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode != 201 {
			return nil, fmt.Errorf("POST %s resulted in %d", fullURL, res.StatusCode)
		}
		if err := json.NewDecoder(res.Body).Decode(&triggers); err != nil {
			return nil, fmt.Errorf("Failed to deserialize response of POST %s (%s)", fullURL, err.Error())
		}
	}
	return &triggers[0], nil
}
