package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// DestructionLevel of type (0 .. 2)
type DestructionLevel int

const (
	// Version of DRAX
	Version string = "0.3.0"
	// DraxPort is the port DRAX is listening on
	DraxPort int = 7777
	// DefaultNumTargets is the number of tasks to kill
	DefaultNumTargets int = 2
)

const (
	// DlBasic means destroy random tasks
	DlBasic DestructionLevel = iota
	// DlAdvanced means destroy random apps
	DlAdvanced
	// DlAll means destroy random apps and services
	DlAll
)

var (
	mux                *http.ServeMux
	marathonURL        string
	destructionLevel   = DestructionLevel(DlBasic)
	numTargets         = int(DefaultNumTargets)
	overallTasksKilled uint64
)

func init() {
	mux = http.NewServeMux()

	// per default, use the cluster-internal, non-auth endpoint:
	marathonURL = "http://marathon.mesos:8080"
	if murl := os.Getenv("MARATHON_URL"); murl != "" {
		marathonURL = murl
	}
	log.WithFields(log.Fields{"main": "init"}).Info("Using Marathon at  ", marathonURL)

	if dl := os.Getenv("DESTRUCTION_LEVEL"); dl != "" {
		l, _ := strconv.Atoi(dl)
		destructionLevel = DestructionLevel(l)
	}
	log.WithFields(log.Fields{"main": "init"}).Info("On destruction level ", destructionLevel)

	if nt := os.Getenv("NUM_TARGETS"); nt != "" {
		n, _ := strconv.Atoi(nt)
		numTargets = n
	}
	log.WithFields(log.Fields{"main": "init"}).Info("I will destroy ", numTargets, " tasks on a rampage")

	if ll := os.Getenv("LOG_LEVEL"); ll != "" {
		switch strings.ToUpper(ll) {
		case "DEBUG":
			log.SetLevel(log.DebugLevel)
		case "INFO":
			log.SetLevel(log.InfoLevel)
		default:
			log.SetLevel(log.ErrorLevel)
		}
	}
}

func main() {
	log.Info("This is DRAX in version ", Version, " listening on port ", DraxPort)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"handle": "/health"}).Info("health check")
		fmt.Fprint(w, "I am Groot")
	})
	mux.Handle("/stats", new(NounStats))
	mux.Handle("/rampage", new(NounRampage))
	p := strconv.Itoa(DraxPort)
	log.Fatal(http.ListenAndServe(":"+p, mux))
}
