# hanime
[![Go](https://img.shields.io/github/go-mod/go-version/IchBinLeoon/hanime?style=flat-square)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/IchBinLeoon/hanime?style=flat-square)](https://github.com/IchBinLeoon/hanime/releases)
[![Commit](https://img.shields.io/github/last-commit/IchBinLeoon/hanime?style=flat-square)](https://github.com/IchBinLeoon/hanime/commits/main)
[![License](https://img.shields.io/github/license/IchBinLeoon/hanime?style=flat-square)](https://github.com/IchBinLeoon/hanime/blob/main/LICENSE)

Command-line tool to download videos from hanime.tv

- [Requirements](#Requirements)
- [Installation](#Installation)
  - [Install via go get](#Install-via-go-get)
  - [Install from source](#Install-from-source)
  - [Install from release](#Install-from-release)
- [Usage](#Usage)
  - [Download a video](#Download-a-video)
  - [Specify the video quality](#Specify-the-video-quality)
  - [Specify a custom output path and name](#Specify-a-custom-output-path-and-name)
  - [Use a proxy](#Use-a-proxy)
  - [Display video info](#Display-video-info)
- [Contribute](#Contribute)
- [License](#License)

## Requirements
- [FFmpeg](https://www.ffmpeg.org/)

## Installation
### Install via go get
Make sure you have [Go](https://golang.org/) installed.
```
go get -u github.com/IchBinLeoon/hanime
```

### Install from source
Make sure you have [Go](https://golang.org/) installed.
```
git clone https://github.com/IchBinLeoon/hanime
cd hanime
go build
```
You should now be provided with an executable.

### Install from release
If you don't want to build the cli yourself, you can download an executable file [here](https://github.com/IchBinLeoon/hanime/releases).

## Usage
### Download a video
```
hanime get https://hanime.tv/videos/hentai/XXX
```

### Specify the video quality
The `-q` or `--quality` flag sets the video quality. Default is 1080.
```
hanime get https://hanime.tv/videos/hentai/XXX -q 720
```

### Specify a custom output path and name
The `-o` or `--output` flag sets a custom output path and the `-O` or `--Output` flag sets a custom output name.
```
hanime get https://hanime.tv/videos/hentai/XXX -o /home/ichbinleoon/XXX -O XXX.mp4
```

### Use a proxy
The `-p` or `--proxy` flag sets a proxy.
```
hanime get https://hanime.tv/videos/hentai/XXX -p XXX://host:port 
```

### Display video info
The `-i` or `--info` flag displays information about the video.
```
hanime get https://hanime.tv/videos/hentai/XXX -i
```

## Contribute
Contributions are welcome! Feel free to open issues or submit pull requests!

## License
MIT © [IchBinLeoon](https://github.com/IchBinLeoon/hanime/blob/main/LICENSE)
