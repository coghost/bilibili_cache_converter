package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/coghost/pathlib"
	"github.com/coghost/xpretty"
	"github.com/joho/godotenv"
)

type Args struct {
	InputDir string `arg:"-i,--input-dir,env:BL_INPUT_DIR" help:"Directory to bilibili cache root dir"`
	// OutputDir where to save the converted files
	OutputDir string `arg:"-o,--output-dir,env:BL_OUTPUT_DIR" help:"Directory to save converted files"`
	Ffmpeg    string `arg:"--ffmpeg-bin,env:BL_FFMPEG" help:"Path to ffmpeg binary"`

	// Actions
	By string `arg:"--by" default:"group" help:"Conversion scope: g(group) /v(video)"`

	// Scan do scan
	Scan bool `arg:"--scan" default:"false" help:"Scan and list(instead of converting) available cache files with the type from '--by'"`
	// Force merge, in case you want to overwrite existed one.
	Force bool `arg:"--force" default:"false" help:"Force merge even if output file already exists"`

	// Clean up cleans downloaded cache file.
	Clean bool `arg:"--clean" default:"false" help:"Clean cache(Warn: cached files will be delete forever)"`

	// GetSubtitle will try to get subtitle from Internet.
	GetSubtitle bool `arg:"--subtitle" default:"false" help:"Download subtitle(may not working)"`
	// use uploader name as subdir or not
	UploaderAsSubDir bool `arg:"--uploader-as-subdir" default:"false" help:"Use uploader name as a subdirectory of the output dir"`

	InitEnv bool `arg:"--init" help:"Init the running env(.env) file"`
	DryRun  bool `arg:"--dry-run" help:"Print arguments and exit without converting"`
	Version bool `arg:"--version" help:"Display version and exit"`
}

func (Args) Description() string {
	filename := "\033[32;4;2m" + filepath.Base(os.Args[0]) + "\033[0m"

	return filename + " is a lightweight bilibili cache converter (m4s to mp4).\n"
}

func LoadEnvArgs(versions ...string) *Args {
	_ = godotenv.Load()

	args := &Args{}
	arg.MustParse(args)

	if args.InitEnv {
		initRunningEnv()
		os.Exit(0)
	}

	xpretty.GreenPrintf("%s\n Input dir: %s\nOutput dir: %s\n%s\n", strings.Repeat("-", 32), args.InputDir, args.OutputDir, strings.Repeat("-", 32))

	showVersionAndExit(args.Version, versions...)
	if args.DryRun {
		dryRunAndExit(args)
	}

	if err := args.Validate(); err != nil {
		xpretty.PrintToStderr("Invalid argument: %v\n", err)
		os.Exit(1)
	}

	return args
}

func (args *Args) Validate() error {
	if args.InputDir == "" || args.OutputDir == "" {
		xpretty.PrintToStderr("InputDir/OutputDir is missing:\n - Check -h/--help for usage;\n - Or use --init-env to add a sample .env config.\n")
		dryRunAndExit(args)
	}

	return nil
}

func dryRunAndExit(args *Args) {
	_ = xpretty.PrettyStruct(args)

	os.Exit(0)
}

func showVersionAndExit(show bool, vers ...string) {
	if !show {
		return
	}

	xpretty.PrintToStdout("Current Version: %s\n", vers[0])
	os.Exit(0)
}

func initRunningEnv() {
	raw := `BL_INPUT_DIR=~/Movies/bilibili/
BL_OUTPUT_DIR=/tmp/bilibili
BL_FFMPEG=ffmpeg`

	dotEnv := pathlib.Path(".env")
	if dotEnv.Exists() {
		log.Printf(".env file is already existed.\n")
		log.Printf("Please manually add following data.\n\n%s\n", xpretty.Cyanf(raw))
		return
	}

	if err := pathlib.Path(".env").WriteText(raw); err != nil {
		log.Printf("cannot add .env file, please add it manually with following data.\n%s\n", raw)
		return
	}

	log.Printf(".env file created.")
}
