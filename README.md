# bilibili_cache_converter

A lightweight bilibili cached m4s to mp4 converter. (bilibili 缓存 m4s 到 mp4 转换器)

## Usage

A lightweight bilibili cache converter (m4s to mp4).

### Basic Usage

```sh
bilibili_cache_converter [OPTIONS]
```

> Quick Start

```sh
bilibili_cache_converter -h
```

### Options:

- `-i, --input-dir <DIR>` (env: `BL_INPUT_DIR`)
  : Directory to the cached files.
- `--by <SCOPE>` (default: `group`)
  : Conversion scope: `g` (group) / `v` (video).
- `--scan <TYPE>` (default: `g`)
  : Scan and list available cache files with given type: `g` (group info) / `v` (video info).
- `--force`
  : Force merge even if output file already exists.
- `--clean`
  : Clean bilibili cache files
- `--subtitle`
  : Download subtitle from a third party website.
- `-o, --output-dir <DIR>` (env: `BL_OUTPUT_DIR`)
  : Directory to save converted files.
- `--uploader-as-subdir`
  : Use uploader name as a subdirectory of the output dir.
- `--dry-run`
  : Print parsed arguments and exit without converting.
- `--version`
  : Display version and exit.

### Examples:

1.  **Scan for available video groups in an input directory:**

    ```sh
    bilibili_cache_converter -i /path/to/bilibili/cache --scan g
    ```

2.  **Convert all videos in a group within an input directory to an output directory:**

    ```sh
    bilibili_cache_converter -i /path/to/bilibili/cache -o /path/to/output --by group
    ```

3.  **Run bilibili_cache_converter directly, no options required:**
    ```sh
    # if .env is found and input_dir/output_dir are added, just run it directly
    bilibili_cache_converter
    ```

### .env

We can add a .env file at the same dir with `bilibili_cache_converter` for global input/output dir.

```sh
BL_INPUT_DIR="~/Movies/bilibili/"
BL_OUTPUT_DIR="/tmp/bilibili"
BL_FFMPEG="/path/to/ffmpeg"
```

### FAQ

#### .1 input-dir: the root dir of bilibili cache root dir, check from your bilibili client for details.
#### .2 output-dir: where the converted mp4 file saved.
