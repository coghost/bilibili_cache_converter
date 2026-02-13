package utils

import "os"

const _ffmpeg = "ffmpeg"

func ConvertWithFfmpeg(inputFiles []string, output string, ffmpegBins ...string) (string, error) {
	args := []string{}
	for _, file := range inputFiles {
		args = append(args, "-i", file)
	}

	fixedArgs := []string{
		"-c:v", "copy",
		"-c:a", "copy",
		"-strict", "experimental",
		"-hide_banner",
		"-stats",
	}

	outputArgs := []string{
		"-y",
		output,
	}

	args = append(args, fixedArgs...)
	args = append(args, outputArgs...)

	bin := os.Getenv("BL_FFMPEG")
	if len(ffmpegBins) != 0 {
		bin = ffmpegBins[0]
	}

	if bin == "" {
		bin = _ffmpeg
	}

	return RunCommand(bin, args)
}
