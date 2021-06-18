# hanime
[![Release](https://img.shields.io/github/v/release/IchBinLeoon/hanime?style=flat-square)](https://github.com/IchBinLeoon/hanime/releases)
[![Commit](https://img.shields.io/github/last-commit/IchBinLeoon/hanime?style=flat-square)](https://github.com/IchBinLeoon/hanime/commits/main)
[![License](https://img.shields.io/github/license/IchBinLeoon/hanime?style=flat-square)](https://github.com/IchBinLeoon/hanime/blob/main/LICENSE)

Command-line tool to download videos from hanime.tv

## Requirements
- Go
- FFmpeg

## Installation
```
go get -u github.com/IchBinLeoon/hanime
```



## Usage
```
hanime get https://hanime.tv/videos/hentai/XXX
```

### Specify output path and name
The `-o` or `--output` flag sets the output path and name
```
hanime get https://hanime.tv/videos/hentai/XXX -o /home/ichbinleoon/XXX.mp4
```

### Specify a proxy
The `-p` or `--proxy` flag sets a proxy
```
hanime get https://hanime.tv/videos/hentai/XXX -p XXX://host:port 
```

## Contribute
Contributions are welcome! Feel free to open issues or submit pull requests!

## License
MIT © [IchBinLeoon](https://github.com/IchBinLeoon/hanime/blob/main/LICENSE)
