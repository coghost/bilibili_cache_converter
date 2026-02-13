package bilibili

import (
	"encoding/json"
	"fmt"

	"github.com/coghost/bilibili_cache_converter/utils"
	"github.com/coghost/pathlib"
	"github.com/spf13/cast"
)

type VideoInfo struct {
	Type    string `json:"type"`
	Codecid int    `json:"codecid"`

	// UID
	UID any `json:"uid"`

	// GroupID/ItemID may have both number and string, we've converted it to string when unmarshal in getVideoInfo
	GroupID    string `json:"-"`
	ItemID     string `json:"-"`
	GroupIDRaw any    `json:"groupId"`
	ItemIDRaw  any    `json:"itemId"`

	Aid            int     `json:"aid"`
	Cid            int     `json:"cid"`
	Bvid           string  `json:"bvid"`
	P              int     `json:"p"`
	TabP           int     `json:"tabP"`
	TabName        string  `json:"tabName"`
	Uname          string  `json:"uname"`
	Avatar         string  `json:"avatar"`
	CoverURL       string  `json:"coverUrl"`
	Title          string  `json:"title"`
	Duration       int     `json:"duration"`
	GroupTitle     string  `json:"groupTitle"`
	GroupCoverURL  string  `json:"groupCoverUrl"`
	Danmaku        int     `json:"danmaku"`
	View           int     `json:"view"`
	Pubdate        int     `json:"pubdate"`
	Vt             int     `json:"vt"`
	Status         string  `json:"status"`
	Active         bool    `json:"active"`
	Loaded         bool    `json:"loaded"`
	Qn             int     `json:"qn"`
	AllowHEVC      bool    `json:"allowHEVC"`
	CreateTime     int64   `json:"createTime"`
	CoverPath      string  `json:"coverPath"`
	GroupCoverPath string  `json:"groupCoverPath"`
	UpdateTime     int64   `json:"updateTime"`
	TotalSize      int     `json:"totalSize"`
	LoadedSize     int     `json:"loadedSize"`
	Progress       float64 `json:"progress"`
	Speed          int     `json:"speed"`
	CompletionTime int64   `json:"completionTime"`
	ReportedSize   int     `json:"reportedSize"`
}

// FilenameFromGroupAndVideo generates filename as `GroupTitle/Title.mp4`,
// and will replace common unsafe characters with underscores
func (v *VideoInfo) FilenameFromGroupAndVideo() string {
	grpTitle := utils.SanitizeFilename(v.GroupTitle)
	vidTitle := utils.SanitizeFilename(v.Title)
	outMP4 := fmt.Sprintf("%s/%s", grpTitle, vidTitle)

	return outMP4
}

// FilenameFromGroupAndVideo generates filename as `Uname/GroupTitle/Title.mp4`,
// and will replace common unsafe characters with underscores
func (v *VideoInfo) FilenameFromUnameGroupAndVideo() string {
	grpTitle := utils.SanitizeFilename(v.GroupTitle)
	vidTitle := utils.SanitizeFilename(v.Title)
	uname := utils.SanitizeFilename(v.Uname)

	outMP4 := fmt.Sprintf("%s/%s/%s", uname, grpTitle, vidTitle)

	return outMP4
}

func (v *VideoInfo) URLWithP() string {
	videoURL := fmt.Sprintf("https://www.bilibili.com/video/%s/?p=%d", v.Bvid, v.P)

	return videoURL
}

func ParseVideoInfo(file string) (*VideoInfo, error) {
	data, err := pathlib.Path(file).GetBytes()
	if err != nil {
		return nil, err
	}

	var video *VideoInfo

	if err := json.Unmarshal(data, &video); err != nil {
		return nil, err
	}

	video.GroupID = cast.ToString(video.GroupIDRaw)
	video.ItemID = cast.ToString(video.ItemIDRaw)

	return video, err
}
