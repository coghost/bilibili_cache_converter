/*
Package bbot bot that use chrome to crawl
*/
package bbot

import (
	"log"
	"os"
	"time"

	"github.com/coghost/wee"
	"github.com/spf13/cast"
)

func NewBot() *wee.Bot {
	proxyFolder := ""

	if folder := os.Getenv("BL_PROXY_FOLDER"); folder != "" {
		log.Printf("use local proxy extension dir: %s", folder)
		proxyFolder = folder
	}

	var headless bool

	val, found := os.LookupEnv("BL_HEADLESS")
	if found {
		headless = cast.ToBool(val)
	}

	options := []wee.BotOption{
		wee.Headless(headless),
		wee.Devtools(true),
		wee.WithBounds(wee.NewBoundsHD()),
	}

	options = append(
		options,
		wee.WithBrowserOptions([]wee.BrowserOptionFunc{
			wee.BrowserExtensions(proxyFolder),
		}))

	bot := wee.NewBotWithOptionsOnly(options...)

	return bot
}

type TimedTask func() error

func RunWithMinimumTime(task TimedTask, minDur time.Duration) error {
	start := time.Now()

	err := task()

	elapsed := time.Since(start)
	if elapsed < minDur {
		time.Sleep(minDur - elapsed)
	}

	return err
}
