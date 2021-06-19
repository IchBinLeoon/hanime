# hanime
[![Release](https://img.shields.io/github/v/release/IchBinLeoon/hanime?style=flat-square)](https://github.com/IchBinLeoon/hanime/releases)
[![Commit](https://img.shields.io/github/last-commit/IchBinLeoon/hanime?style=flat-square)](https://github.com/IchBinLeoon/hanime/commits/main)
[![License](https://img.shields.io/github/license/IchBinLeoon/hanime?style=flat-square)](https://github.com/IchBinLeoon/hanime/blob/main/LICENSE)

Command-line tool to download videos from hanime.tv

- [Installation](#Installation)
  - [Requirements](#Requirements)
  - [Install via `go get`](#Install-via-go-get)
  - [Install from source](#Install-from-source)
- [Usage](#Usage)
  - [Download a video](#Download-a-video)
  - [Specify the video quality](#Specify-the-video-quality)
  - [Specify a custom output path and name](#Specify-a-custom-output-path-and-name)
  - [Use a proxy](#Use-a-proxy)
  - [Display video info](#Display-video-info)
- [Contribute](#Contribute)
- [License](#License)

## Installation
### Requirements
- [Go](https://golang.org/)
- [FFmpeg](https://www.ffmpeg.org/)

### Install via `go get`
```
go get -u github.com/IchBinLeoon/hanime
```

### Install from source
```
git clone https://github.com/IchBinLeoon/hanime
cd hanime
go build
```



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
