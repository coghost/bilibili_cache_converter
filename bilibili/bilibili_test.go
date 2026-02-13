package bilibili

import (
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/coghost/bilibili_cache_converter/fixtures/testutil"
	"github.com/coghost/bilibili_cache_converter/utils"
	"github.com/coghost/pathlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_testOutDir   = path.Join(os.TempDir(), "bilibili")
	_testInputDir = path.Join(testutil.GetProjectRoot(), "fixtures")
)

func mustFileStat(fs *pathlib.FsPath) fs.FileInfo {
	st, _ := fs.Stat()
	return st
}

func TestWithDetailSteps(t *testing.T) {
	assert := assert.New(t)

	const (
		videoID = "26349405204"
		bvid    = "BV1JcCUYSEEL"
		itemID  = "26349405204"

		uname      = "乐乐乐雨_"
		groupTitle = "【星露谷物语】复古小卧室"
		// The group name of the current video happens to be the same as the video name
		videoTitle = groupTitle
	)

	inputFs := pathlib.Path(_testInputDir).Join(videoID)
	outputFs := pathlib.Path(_testOutDir)

	err := outputFs.Mkdirs()
	require.NoError(t, err, "mkdir for output")

	assert.Equal(videoID, inputFs.Stem, "video id")
	assert.True(inputFs.IsDir(), "should be dir")

	files, err := inputFs.ListFilesWithGlob("*.m4s")
	require.NoError(t, err, "list m4s files")
	assert.Len(files, 2, "total 2 m4s files")

	videoInfo, err := ParseVideoInfo(inputFs.Join(_videoInfoFile).AbsPath())
	require.NoError(t, err, "parse video info")

	assert.Equal(bvid, videoInfo.Bvid, "bvid")
	assert.Equal(itemID, videoInfo.ItemID, "item id")

	m4sfiles := []string{}

	for _, file := range files {
		inFs := pathlib.Path(file)
		outFs := outputFs.Join(inFs.Name)

		n, err := copyWithout9zeroPrefix(file, outFs.AbsPath())
		assert.Greater(n, int64(0), "copy to new file")
		require.NoError(t, err, "copy")

		assert.Equal(int64(_cachedM4SHeaderLen), mustFileStat(inFs).Size()-mustFileStat(outFs).Size(), "copied file less than orig")

		m4sfiles = append(m4sfiles, outFs.AbsPath())
	}

	tests := []struct {
		name     string
		nameFn   func() string
		wantName string
	}{
		{
			name:     "use group and title",
			nameFn:   videoInfo.FilenameFromGroupAndVideo,
			wantName: path.Join(_testOutDir, groupTitle, videoTitle+_outputVideoDotMP4),
		},
		{
			name:     "use uname group and title",
			nameFn:   videoInfo.FilenameFromUnameGroupAndVideo,
			wantName: path.Join(_testOutDir, uname, groupTitle, videoTitle+_outputVideoDotMP4),
		},
	}
	for _, tt := range tests {
		outMP4 := tt.nameFn() + _outputVideoDotMP4
		outputMP4Fs := outputFs.Join(outMP4)

		assert.Equal(tt.wantName, outputMP4Fs.AbsPath(), tt.name+" output filename")

		err = outputMP4Fs.MkParentDir()
		require.NoError(t, err, tt.name+" mkdir for mp4 file")

		res, err := utils.ConvertWithFfmpeg(m4sfiles, outputMP4Fs.AbsPath())
		assert.NotEmpty(t, res, tt.name+" convert with ffmpeg")
		require.NoError(t, err, tt.name+" merge m4s to mp4")
	}
}
