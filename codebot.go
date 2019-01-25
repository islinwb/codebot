package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/golang/glog"
)

const ConfigFilePath = "config.json"

func main() {
	var codebotServer CodebotServer
	err := ReadConfig(&codebotServer)
	if err != nil {
		glog.Fatalf("fail to read config: %v", err)
	}

	// health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			glog.Errorf("fail to write: %v", err)
		}
	})

	http.HandleFunc("/", codebotServer.Handler)

	err = http.ListenAndServe(codebotServer.Addr, nil)
	if err != nil {
		glog.Fatalf("fail to start server: %v", err)
	}

	// monitor the config file
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		glog.Fatalf("fail to create new watcher: %v", err)
	}
	defer watch.Close()
	err = watch.Add(ConfigFilePath)
	if err != nil {
		glog.Fatalf("fail to add file to watch: %v", err)
	}
	go func() {
		for {
			select {
			case ev := <-watch.Events:
				{
					if ev.Op&fsnotify.Write == fsnotify.Write {
						glog.Infof("The config file was changed with Write. Read the config again...")
						err := ReadConfig(&codebotServer)
						if err != nil {
							glog.Fatalf("fail to read config: %v", err)
						}
					}
				}
			}
		}
	}()
}

type CodebotServer struct {
	Addr       string           `json:"addr"`
	Token      string           `json:"token"`
	JenkinsMap []JenkinsMapping `json:"jenkins_map"`
}

type JenkinsMapping struct {
	TargetBranch string `json:"target_branch"`
	Label        string `json:"label"`
	JenkinsAddr  string `json:"jenkins_addr"`
}

const (
	PushHookEvent         = "Push Hook"
	TagPushHookEvent      = "Tag Push Hook"
	IssueHookEvent        = "Issue Hook"
	NoteHookEvent         = "Note Hook"
	MergeRequestHookEvent = "Merge Request Hook"
)

func (s *CodebotServer) Handler(w http.ResponseWriter, r *http.Request) {
	eventType := r.Header.Get("X-Gitlab-Event")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Errorf("fail to read request body: %v", err)
	}
	defer r.Body.Close()

	switch eventType {
	case PushHookEvent:
		glog.Infof("It's a Push event")
	case NoteHookEvent:
		glog.Infof("It's a Note event")
		s.handleNote(b)
	case MergeRequestHookEvent:
		glog.Infof("It's a Merge Request event")
		s.handleMergeRequest(b)
	case IssueHookEvent:
		glog.Infof("It's a Issue event")
	default:
		glog.Infof("It's not supported yet: %s", eventType)
	}
}

func ReadConfig(server *CodebotServer) error {
	b, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		glog.Errorf("fail to read config file: %v", err)
		return err
	}
	err = json.Unmarshal(b, &server)
	if err != nil {
		glog.Errorf("fail to unmarshal: %v", err)
	}
	return err
}
