/*
Package kedou will download subtitles from www.kedou.life/caption/subtitle/bilibili
*/
package kedou

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/coghost/bilibili_cache_converter/subtitles"
	"github.com/coghost/bilibili_cache_converter/subtitles/bbot"
	"github.com/coghost/pathlib"
	"github.com/coghost/wee"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
)

type Subtitle struct {
	Lang     string `json:"lang,omitempty"`
	LangDesc string `json:"langDesc,omitempty"`
	Content  string `json:"content,omitempty"`
}

type SubtitleInfo struct {
	Vid       string `json:"vid,omitempty"`
	Host      string `json:"host,omitempty"`
	HostAlias string `json:"hostAlias,omitempty"`
	Title     string `json:"title,omitempty"`
	Status    string `json:"status,omitempty"`

	Subtitles []Subtitle `json:"subtitleItemVoList,omitempty"`
}

const (
	url     = `https://www.kedou.life/caption/subtitle/bilibili`
	input   = `input.el-input__inner`
	submit  = `button.el-button`
	items   = `a@@@下载`
	noitems = `div.shouldnotexisted`

	script = `() => { return JSON.stringify(window.__NUXT__.pinia['captionStore']['subtitleExtractInfo']) }`
)

var ErrCode500 = errors.New("code 500 found, server issue")

type SubtitleManger struct {
	bot *wee.Bot

	raw string
}

func NewSubtitleManager() *SubtitleManger {
	return &SubtitleManger{}
}

func (m *SubtitleManger) CleanUp() {
	if m.bot != nil {
		m.bot.Cleanup()
	}
}

func (m *SubtitleManger) GetRawString() string {
	return m.raw
}

func (m *SubtitleManger) Scrape(cacheFs *pathlib.FsPath, videoURL string) (*SubtitleInfo, error) {
	var subInfo *SubtitleInfo
	// reset raw
	m.raw = ""

	if cacheFs.Exists() {
		err := cacheFs.GetJSON(&subInfo)
		if err != nil {
			return nil, err
		}

		if len(subInfo.Subtitles) > 0 {
			return subInfo, err
		}

		log.Printf("no subtitles found in cache, search online")
	}

	subInfo, err := m.scrape(videoURL)
	if err != nil {
		return nil, err
	}

	if m.GetRawString() == "" {
		return nil, fmt.Errorf("raw is empty: %w", subtitles.ErrNoSubtitlesFound)
	}

	if err := cacheFs.WriteText(m.GetRawString()); err != nil {
		return nil, fmt.Errorf("cannot save raw json: %w", err)
	}

	return subInfo, nil
}

func (m *SubtitleManger) scrape(videoURL string) (*SubtitleInfo, error) {
	if m.bot == nil {
		m.bot = bbot.NewBot()

		wee.BindBotLanucher(m.bot)
		m.bot.DisableImages()
	}

	m.bot.MustOpen(url)
	m.bot.MustInput(input, videoURL)
	m.bot.MustClick(submit)

	timeout := cast.ToInt(os.Getenv("BL_TIMEOUT"))
	if timeout == 0 {
		timeout = wee.PT20Sec
	}

	elems, _ := m.bot.AnyElem([]string{items, noitems}, wee.WithTimeout(float64(timeout)))
	log.Printf("found total %d elems", len(elems))

	if len(elems) != 0 {
		err := retry.Do(func() error {
			res, err := m.bot.Eval(script)
			if err != nil {
				return err
			}

			m.raw = res.Value.String()

			return nil
		},
			retry.LastErrorOnly(false),
			retry.Attempts(10), //nolint:mnd
			retry.Delay(time.Second*1),
		)
		if err != nil {
			return nil, err
		}
	}

	code := gjson.Get(m.raw, "code")
	if code.String() != "" {
		// handle raw with {"code":500,"message":"未知错误，请确保链接正确！","data":null}
		if code.String() == "500" {
			return nil, ErrCode500
		}
	} else {
		log.Printf("try title: %s", gjson.Get(m.raw, "title"))
	}

	var sub SubtitleInfo

	if err := json.Unmarshal([]byte(m.raw), &sub); err != nil {
		return nil, err
	}

	return &sub, nil
}
