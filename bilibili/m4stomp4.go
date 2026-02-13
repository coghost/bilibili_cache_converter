package bilibili

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/coghost/bilibili_cache_converter/utils"
	"github.com/coghost/pathlib"
)

func ConvertVideo(options *Options) (string, error) {
	// inputDir, outputDir string, forceMerge bool, useUname bool
	inputFs := pathlib.Path(options.InputDir).ExpandUser()
	outputFs := pathlib.Path(options.OutputDir).ExpandUser()
	log.Printf("input/output: %s vs %s\n", inputFs, outputFs)

	pattern := _inputSuffix
	if !strings.HasPrefix(pattern, "*") {
		pattern = "*" + _inputSuffix
	}

	files, err := inputFs.ListFilesWithGlob(pattern)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", ErrNoM4S
	}

	videoInfo, err := ParseVideoInfo(inputFs.Join(_videoInfoFile).AbsPath())
	if err != nil {
		return "", err
	}

	outMP4 := videoInfo.FilenameFromGroupAndVideo() + _outputVideoDotMP4
	if options.UseUploaderAsSubDir {
		outMP4 = videoInfo.FilenameFromUnameGroupAndVideo()
	}

	outputMP4Fs := outputFs.Join(outMP4)
	if err := outputMP4Fs.MkParentDir(); err != nil {
		return "", err
	}

	if !options.ForceMerge && outputMP4Fs.Exists() {
		log.Printf("already converted, skip: %s", inputFs)
		return outputMP4Fs.AbsPath(), nil
	}

	m4sfiles := []string{}

	for _, file := range files {
		name := pathlib.Path(file).Name
		outFile := outputFs.Join(name).AbsPath()

		if _, err := copyWithout9zeroPrefix(file, outFile); err != nil {
			return "", err
		}

		m4sfiles = append(m4sfiles, outFile)
	}

	_, err = utils.ConvertWithFfmpeg(m4sfiles, outputMP4Fs.AbsPath())
	if err != nil {
		return "", err
	}

	for _, file := range m4sfiles {
		os.Remove(file)
	}

	return outMP4, err
}

func copyWithout9zeroPrefix(srcFile, dstFile string) (int64, error) {
	fin, err := os.Open(srcFile)
	if err != nil {
		return 0, nil
	}

	defer fin.Close()

	var header [_cachedM4SHeaderLen]byte
	if _, err := io.ReadFull(fin, header[:]); err != nil {
		return 0, err
	}

	if string(header[:]) != _cachedM4SHeader {
		return 0, fmt.Errorf("not 9 zero header found: %w", ErrNoCachePrefix)
	}

	fout, err := os.Create(dstFile)
	if err != nil {
		return 0, nil
	}
	defer fout.Close()

	// Offset is the number of bytes you want to exclude
	_, err = fin.Seek(_cachedM4SHeaderLen, io.SeekStart)
	if err != nil {
		return 0, nil
	}

	n, err := io.Copy(fout, fin)

	return n, err
}
