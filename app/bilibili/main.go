package main

import (
	"fmt"
	"log"
	"maps"
	"os"
	"path"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/coghost/bilibili_cache_converter/bilibili"
	"github.com/coghost/bilibili_cache_converter/subtitles/bbot"
	"github.com/coghost/bilibili_cache_converter/subtitles/kedou"
	"github.com/coghost/bilibili_cache_converter/utils"
	"github.com/coghost/bilibili_cache_converter/versions"
	"github.com/coghost/pathlib"
	"github.com/coghost/sleep"
	"github.com/coghost/wee"
	"github.com/coghost/xpretty"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

const (
	byG     = "g"
	byGroup = "group"
	byV     = "v"
	byVideo = "video"
)

func main() {
	args := LoadEnvArgs(versions.Version)
	run(args)
}

func run(args *Args) {
	options := &bilibili.Options{
		InputDir:            args.InputDir,
		OutputDir:           args.OutputDir,
		ForceMerge:          args.Force,
		UseUploaderAsSubDir: args.UploaderAsSubDir,
	}

	if args.Scan {
		scanLocal(args)
		return
	}

	if args.GetSubtitle {
		downloadSubtitle(args)
		return
	}

	if args.Clean {
		scanAndClean(args)
		return
	}

	// by default will do convert video action
	convertVideos(args, options)
}

func convertVideos(args *Args, options *bilibili.Options) {
	videos := bilibili.SelectVideosByGroup(args.InputDir)
	if len(videos) == 0 {
		log.Printf("cannot get videos by group, end!")
		os.Exit(0)
	}

	var err error

	bcvc := bilibili.NewCacheVideoConverter(options, nil)

	switch args.By {
	case byG, byGroup:
		err = bcvc.ConvertByGroup(videos[0].GroupID)
	case byV, byVideo:
		video := bilibili.SelectVideo(videos)
		err = bcvc.ConvertByVideo(video.ItemID)
	default:
		log.Printf("(%s) not supported, check --help for usage", args.By)
		os.Exit(0)
	}

	if err != nil {
		log.Printf("convert failed: %v", err)
	}
}

func scanAndClean(args *Args) {
	videos := bilibili.SelectVideosByGroup(args.InputDir)
	if len(videos) == 0 {
		log.Printf("cannot get videos by group, end!")
		os.Exit(0)
	}

	log.Printf("running on group: %s", videos[0].GroupTitle)

	video := bilibili.SelectVideo(videos)

	log.Printf("running on video: %s:%s", video.Title, video.ItemID)
}

func scanLocal(args *Args) {
	videoGroups, err := bilibili.ScanForAllVideoGroups(args.InputDir)
	if err != nil {
		log.Printf("scan local groups failed: %v", err)
		os.Exit(-1)
	}

	leveledList := pterm.LeveledList{}

	titles := slices.Sorted(maps.Keys(videoGroups))

	for _, title := range titles {
		videos := videoGroups[title]

		grpMsg := fmt.Sprintf("%s(%s: %d)", title, videos[0].GroupID, len(videos))
		leveledList = append(leveledList, pterm.LeveledListItem{
			Level: 0,
			Text:  xpretty.Cyan(grpMsg),
		})

		if args.By == byG || args.By == byGroup {
			continue
		}

		for _, video := range videos {
			l2msg := fmt.Sprintf("[%d] %s", video.P, video.Title)
			leveledList = append(leveledList, pterm.LeveledListItem{
				Level: 1,
				Text:  l2msg,
			})

			l3msg := fmt.Sprintf("[URL] %s", video.URLWithP())
			leveledList = append(leveledList, pterm.LeveledListItem{
				Level: 2,
				Text:  fmt.Sprintf("[PTH] %s", path.Join(args.InputDir, video.ItemID)),
			}, pterm.LeveledListItem{
				Level: 2,
				Text:  l3msg,
			})
		}
	}

	root := putils.TreeFromLeveledList(leveledList)
	root.Text = xpretty.Yellow("Bilibili Cached Videos")

	_ = pterm.DefaultTree.WithRoot(root).Render()
}

func downloadSubtitle(args *Args) {
	videos := bilibili.SelectVideosByGroup(args.InputDir)
	if len(videos) == 0 {
		log.Printf("cannot get videos by group, end!")
		os.Exit(0)
	}

	sort.Slice(videos, func(i int, j int) bool {
		return videos[i].P < videos[j].P
	})

	for index, video := range videos {
		index += 1
		xpretty.CyanPrintf("[%-2d]: %s\n", index, strings.ReplaceAll(video.Title, "\n", " | "))
	}

	autoNext := true
	currIndex := -1
	downloaded := 0
	tried := 0

	var video *bilibili.VideoInfo

	for {
		newVideoRequired := requireNew(downloaded, &tried, video)
		if newVideoRequired {
			if autoNext && currIndex != -1 {
				currIndex += 1
				log.Printf("auto switch to next video: %d", currIndex+1)
			} else {
				// the index in listing videos starts with 1, so minus 1 here
				currIndex = utils.ScanfInt("Enter a new video") - 1
			}

			video = videos[currIndex]
		}

		xpretty.CyanPrintf("downling subtitle for: %s\n%s\n", video.Title, video.URLWithP())

		if n, err := downloadWithTime(args, video, true); err != nil {
			log.Printf("download failed: %v", err)
			sleep.PT4s()
		} else {
			downloaded = n
		}
	}
}

func requireNew(downloaded int, tried *int, video *bilibili.VideoInfo) bool {
	if video == nil || downloaded > 0 {
		return true
	}

	if *tried > 0 {
		*tried = 0

		return true
	} else {
		*tried += 1

		return false
	}
}

func downloadWithTime(args *Args, video *bilibili.VideoInfo, withAll bool) (int, error) {
	var count int
	var errDownload error

	if err := bbot.RunWithMinimumTime(
		func() error {
			count, errDownload = download(args, video, withAll)

			return errDownload
		},
		wee.PT60Sec*time.Second,
	); err != nil {
		return 0, err
	}

	return count, errDownload
}

func download(args *Args, video *bilibili.VideoInfo, withAll bool) (int, error) {
	subMgr := kedou.NewSubtitleManager()
	defer subMgr.CleanUp()

	// subMgr *kedou.SubtitleManger,
	outFs := pathlib.Path(args.OutputDir)
	grpTitle := utils.SanitizeFilename(video.GroupTitle)
	vidTitle := utils.SanitizeFilename(video.Title)

	cacheFs := outFs.Join(grpTitle, "cache", fmt.Sprintf("%s.raw.json", vidTitle))

	subInfo, err := subMgr.Scrape(cacheFs, video.URLWithP())
	if err != nil {
		return 0, err
	}

	if !withAll {
		for index, st := range subInfo.Subtitles {
			index += 1
			xpretty.CyanPrintf("[%-2d]: %s/%s\n", index, st.Lang, st.LangDesc)
		}

		choice := utils.ScanfInt("Select video for subtitle")
		st := subInfo.Subtitles[choice-1]
		vname := fmt.Sprintf("%s.%s.%d.srt", video.FilenameFromGroupAndVideo(), st.Lang, choice)

		return 0, outFs.Join(vname).WriteText(st.Content)
	}

	for index, st := range subInfo.Subtitles {
		vname := fmt.Sprintf("%s.%s.%d.srt", video.FilenameFromGroupAndVideo(), st.Lang, index)

		err := outFs.Join(vname).WriteText(st.Content)
		if err != nil {
			return 0, err
		}
	}

	num := len(subInfo.Subtitles)

	if num == 0 {
		xpretty.YellowPrintf("no subtitle downloaded: got raw: %s\n", subMgr.GetRawString())
	} else {
		xpretty.GreenPrintf("total %d subtitles added\n", num)
	}

	return num, nil
}
