package bilibili

import "errors"

const (
	_videoInfoFile = "videoInfo.json"
	_inputSuffix   = "m4s"

	_outputVideoDotMP4 = ".mp4"
	// _outputDotSrt      = ".srt"
)

const (
	// bilibili cached header
	_cachedM4SHeaderLen = 9
	// bilibili encryption 9 zero
	_cachedM4SHeader = "000000000"
)

var (
	ErrNoM4S         = errors.New("no m4s file found")
	ErrNoCachePrefix = errors.New("no prefix of 9 zero")
	ErrUserCanceled  = errors.New("user canceled the operation")
	ErrDirNotFound   = errors.New("directly not found")

	ErrNotGroupFolder = errors.New("not a group folder, video folder found")
)
