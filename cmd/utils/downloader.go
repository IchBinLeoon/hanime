package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/grafov/m3u8"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type Downloader struct {
	Client *http.Client
}

func NewDownloader(client *http.Client) *Downloader {
	return &Downloader{
		Client: client,
	}
}

func (downloader *Downloader) Download(m3u8Url string, tmpPath string, outputPath string) error {
	fmt.Print("\n» Creating temporary folder\n")
	err := createTmpFolder(tmpPath)
	if err != nil {
		return err
	}

	fmt.Print("» Parsing M3U8 playlist\n")
	data, err := getM3U8(downloader.Client, m3u8Url)
	if err != nil {
		return err
	}
	p, listType, err := m3u8.DecodeFrom(bytes.NewBuffer(data), true)
	if err != nil {
		return err
	}

	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)

		fmt.Print("» Downloading media files\n")

		var segments []string
		for _, v := range mediapl.Segments {
			if v != nil {
				segments = append(segments, v.URI)
			}
		}

		fileListPath := filepath.Join(tmpPath, "filelist.txt")
		err := createFileList(fileListPath, tmpPath, len(segments))
		if err != nil {
			return err
		}

		key, err := getKey(downloader.Client, mediapl.Key.URI)
		if err != nil {
			return err
		}

		var iv []byte
		if mediapl.Key.IV == "" {
			iv = key
		} else {
			iv = []byte(mediapl.Key.IV)
		}

		bar := NewProgressBar(0, int64(len(segments)), "█")
		wg := sync.WaitGroup{}
		for k, v := range segments {
			wg.Add(1)
			go func(index int, url string) {
				err := downloadTS(downloader.Client, url, filepath.Join(tmpPath, fmt.Sprintf("%s.ts", strconv.Itoa(index))), key, iv)
				if err != nil {
					fmt.Printf("\n\n%s\n", err)
					cleanErr := CleanUp(tmpPath)
					if cleanErr != nil {
						fmt.Println(cleanErr)
					}
					os.Exit(1)
				}
				wg.Done()
				bar.Add(1)
			}(k, v)
		}
		wg.Wait()
		bar.Finish()

		fmt.Print("» Merging media files\n")
		out, err := MergeToMP4(fileListPath, outputPath)
		if err != nil {
			fmt.Println(string(out))
			return err
		}
		fmt.Printf("» Files merged to %s\n", filepath.Base(outputPath))

		fmt.Print("» Cleaning up\n")
		cleanErr := CleanUp(tmpPath)
		if cleanErr != nil {
			return cleanErr
		}

		return nil
	}

	return nil
}

func createTmpFolder(path string) error {
	if CheckIfPathExists(path) {
		fmt.Print(fmt.Errorf("error: cannot create temporary folder, path '%s' already exists", path))
		os.Exit(1)
	}

	err := MakeDirectoryIfNotExists(path)
	if err != nil {
		return err
	}

	return nil
}

func getM3U8(client *http.Client, url string) ([]byte, error) {
	body, err := Request("GET", client, url, nil, nil)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	if err := body.Close(); err != nil {
		return nil, err
	}

	return data, nil
}

func createFileList(fileListPath string, tmpPath string, count int) error {
	f, err := os.Create(fileListPath)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	for i := 0; i < count; i++ {
		_, err := fmt.Fprintln(w, fmt.Sprintf("file '%s.ts'", filepath.Join(tmpPath, strconv.Itoa(i))))
		if err != nil {
			return err
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func getKey(client *http.Client, url string) ([]byte, error) {
	body, err := Request("GET", client, url, nil, nil)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	if err := body.Close(); err != nil {
		return nil, err
	}

	return data, nil
}

func downloadTS(client *http.Client, url string, path string, key []byte, iv []byte) error {
	body, err := Request("GET", client, url, nil, nil)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	decrypted, err := Decrypt(data, key, iv)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err := f.Write(decrypted); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	if err := body.Close(); err != nil {
		return err
	}

	return nil
}
