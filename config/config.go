package config

import (
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"log"
)

type config struct {
	// App directory where static files (presentations) will be placed.
	StaticDir string
	// App a port that will be used by a web-server (in URL part form)
	ServerPort string
	// Enable/disable Sentry usage
	UseSentry bool
	// Sentry setup
	SentryConfiguration sentry.ClientOptions
}

var Config config

func init() {
	// Obtain temporary directory
	dir, err := ioutil.TempDir("", "harbored")
	if err != nil {
		log.Fatal(err)
	}
	Config = config{
		StaticDir:  dir,
		ServerPort: ":8080",
		UseSentry:  false,
		//SentryConfiguration: sentry.ClientOptions{
		//  Dsn:              "",
		//  Environment:      "",
		//  Release:          "",
		//  Debug:            true,
		//  AttachStacktrace: true,
		//},
	}
}
