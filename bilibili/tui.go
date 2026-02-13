package bilibili

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/coghost/bilibili_cache_converter/utils"
	"github.com/coghost/xpretty"
	"github.com/thoas/go-funk"
)

func SelectVideo(groupVideoList []*VideoInfo) *VideoInfo {
	var video *VideoInfo

	for {
		video = selectVideo(groupVideoList)
		if video != nil {
			break
		}
	}

	return video
}

func selectVideo(groupVideoList []*VideoInfo) *VideoInfo {
	for index, video := range groupVideoList {
		index += 1
		xpretty.CyanPrintf("[%-2d]: %s\n", index, strings.ReplaceAll(video.Title, "\n", " | "))
	}

	choice := utils.ScanfInt("Select the video")
	if choice < 0 || choice > len(groupVideoList) {
		return nil
	}

	return groupVideoList[choice-1]
}

func SelectVideosByGroup(inputDir string) []*VideoInfo {
	videoGroups, err := ScanForAllVideoGroups(inputDir)
	if err != nil {
		log.Printf("cannot scan local groups (dir:%s): %v\n", inputDir, err)
		os.Exit(-1)
	}

	if len(videoGroups) == 0 {
		xpretty.YellowPrintf("no cached files found from %s\n", inputDir)
		os.Exit(0)
	}

	var titleSelected string

	for {
		titleSelected = selectGroupByTitle(videoGroups)
		if titleSelected != "" {
			break
		}
	}

	videoGroup := videoGroups[titleSelected]

	return videoGroup
}

func selectGroupByTitle(videoGroups map[string][]*VideoInfo) string {
	titleArr := []string{}

	titledGrps, _ := funk.Keys(videoGroups).([]string)
	sort.Strings(titledGrps)

	for index, title := range titledGrps {
		titleArr = append(titleArr, title)
		videos := videoGroups[title]

		xpretty.CyanPrintf("[%-2d]: [videos:%-2d] %s\n", index+1, len(videos), strings.ReplaceAll(title, "\n", " | "))
	}

	choice := utils.ScanfInt("Select the Group")
	if choice > len(titleArr) || choice < 0 {
		return ""
	}

	titleSelected := titleArr[choice-1]

	return titleSelected
}
