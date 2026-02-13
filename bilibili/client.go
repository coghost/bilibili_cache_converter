package bilibili

import (
	"fmt"
	"io/fs"
	"log"
	"path"

	"github.com/coghost/pathlib"
)

type Options struct {
	InputDir   string
	OutputDir  string
	ForceMerge bool

	UseUploaderAsSubDir bool
}

type converter func(*Options) (string, error)

type CacheVideoConverter struct {
	options *Options
	convert converter
}

func NewCacheVideoConverter(options *Options, convert converter) *CacheVideoConverter {
	if convert == nil {
		convert = ConvertVideo
	}

	return &CacheVideoConverter{
		options: options,
		convert: convert,
	}
}

func (c *CacheVideoConverter) ConvertByVideo(videoID string) error {
	options := c.options
	inputFolder := path.Join(options.InputDir, videoID)

	options.InputDir = inputFolder

	name, err := c.convert(options)
	if err != nil {
		log.Printf("cannot convert %s, %v", inputFolder, err)
	} else {
		log.Printf("converted: %s", name)
	}

	return err
}

func (c *CacheVideoConverter) ConvertByGroup(groupID string) error {
	options := c.options

	inputFs := pathlib.Path(options.InputDir)

	// if videoInfo.json found, means this is video folder
	if inputFs.Join(_videoInfoFile).Exists() {
		return ErrNotGroupFolder
	}

	log.Printf("scan all videos for %s with group: %s", options.InputDir, groupID)

	errWalk := inputFs.Walk(
		func(path string, info fs.FileInfo, _ error) error {
			if path == "." {
				return nil
			}

			if !info.IsDir() {
				return nil
			}

			subDir := inputFs.Join(path)
			videoFs := subDir.Join(_videoInfoFile)

			if !videoFs.Exists() {
				return nil
			}

			videoInfo, err := ParseVideoInfo(videoFs.AbsPath())
			if err != nil {
				log.Printf("cannot get videoInfo for %s", videoFs)
				return err
			}

			if videoInfo.GroupID != groupID {
				// log.Printf("not wanted groupid [wanted/got]: %s != %s", groupID, videoInfo.GroupID)
				return nil
			}

			log.Printf("converting %s...", subDir)

			options.InputDir = subDir.AbsPath()

			name, err := c.convert(options)
			if err != nil {
				log.Printf("cannot convert %s, %v", subDir.AbsPath(), err)
			} else {
				log.Printf("converted: %s", name)
			}

			return err
		})

	return errWalk
}

func ScanForAllVideoGroups(input string) (map[string][]*VideoInfo, error) {
	inputFs := pathlib.Path(input)
	if !inputFs.Exists() {
		return nil, fmt.Errorf("%w: %s", ErrDirNotFound, inputFs)
	}

	videoGroups := make(map[string][]*VideoInfo)

	errWalk := inputFs.Walk(func(path string, info fs.FileInfo, _ error) error {
		if path == "." {
			return nil
		}

		if !info.IsDir() {
			return nil
		}

		subDir := inputFs.Join(path)
		videoFs := subDir.Join(_videoInfoFile)

		if !videoFs.Exists() {
			return nil
		}

		videoInfo, err := ParseVideoInfo(videoFs.AbsPath())
		if err != nil {
			return err
		}

		groupTitle := videoInfo.GroupTitle
		if groupTitle == "" {
			groupTitle = videoInfo.GroupID
		}

		vg := videoGroups[groupTitle]
		if len(vg) == 0 {
			videoGroups[groupTitle] = []*VideoInfo{}
		}

		videoGroups[groupTitle] = append(videoGroups[groupTitle], videoInfo)

		return nil
	})

	return videoGroups, errWalk
}
