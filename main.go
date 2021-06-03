package main

import (
	"context"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"github.com/getsentry/sentry-go"
	"github.com/leaanthony/mewn"
	"harbored/config"
	"harbored/screens/setup-files"
	"harbored/services/presentations"
	"harbored/webserver"
	"os"
	"time"
)

func main() {
	// Sentry setup
	if config.Config.UseSentry {
		err := sentry.Init(config.Config.SentryConfiguration)
		if err != nil {
			fmt.Printf("Sentry initialization failed: %v\n", err)
		}
		defer sentry.Flush(2 * time.Second)
		defer func() {
			err := recover()
			if err != nil {
				switch v := err.(type) {
				case string:
					sentry.CaptureMessage(v)
				case error:
					sentry.CaptureException(v)
				}
				panic(err)
			}
		}()
	}

	// Start web-server
	go webserver.Server.Start(config.Config.ServerPort)

	// UI setup
	a := app.New()
	a.Settings().SetTheme(HarboredTheme())
	ico := mewn.MustBytes("./Icon.png")
	a.SetIcon(&Font{
		name:    "Icon.png",
		content: ico,
	})

	// Display main window
	w := a.NewWindow("Harbored")
	w.Resize(fyne.NewSize(1024, 768))
	view := setup_files.NewFilesView(&a, &w)
	view.Render()
	w.CenterOnScreen()
	w.SetMaster()
	w.ShowAndRun()

	// Cleanup
	presentations.End()
	if config.Config.StaticDir != "" {
		os.RemoveAll(config.Config.StaticDir)
	}
	webserver.Server.Shutdown(context.Background())
}
