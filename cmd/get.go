package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/IchBinLeoon/hanime/types"
	"github.com/IchBinLeoon/hanime/utils"
	"github.com/grafov/m3u8"
	"github.com/spf13/cobra"
)

var tmpPath string
var outputPath string
var fileListPath string

var videoQualities = []string{
	"1080",
	"720",
	"480",
	"360",
}

var qualityFlag string
var outputPathFlag string
var outputNameFlag string
var proxyFlag string

var getUsage = `Usage:
  hanime get <url> [flags]

Flags:
  -h, --help      help for get
  -q, --quality	  video quality (default 1080)
  -o, --output    custom output path
  -O, --Output    custom output name
  -p, --proxy     proxy url
`

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.SetUsageTemplate(getUsage)
	getCmd.Flags().StringVarP(&qualityFlag, "quality", "q", "1080", "video quality")
	getCmd.Flags().StringVarP(&outputPathFlag, "output", "o", "", "custom output path")
	getCmd.Flags().StringVarP(&outputNameFlag, "Output", "O", "", "custom output name")
	getCmd.Flags().StringVarP(&proxyFlag, "proxy", "p", "", "proxy url")
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Download video by url",
	Long:  "Download a video from hanime.tv by url",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utils.CatchInterrupt(&tmpPath)
		if err := get(args[0]); err != nil {
			fmt.Println(err)
			cleanErr := utils.CleanUp(tmpPath)
			if cleanErr != nil {
				fmt.Println(cleanErr)
			}
			os.Exit(1)
		}
	},
}

func get(url string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	video, err := getVideo(client, url)
	if err != nil {
		return err
	}

	index, err := getStreamIndex(video.VideosManifest.Servers[0].Streams)
	if err != nil {
		return err
	}

	name := video.HentaiVideo.Name
	slug := video.HentaiVideo.Slug
	quality := video.VideosManifest.Servers[0].Streams[index].Height
	size := video.VideosManifest.Servers[0].Streams[index].Size
	id := video.VideosManifest.Servers[0].Streams[index].ID

	pathsErr := setPaths(slug, quality, id)
	if pathsErr != nil {
		return pathsErr
	}

	fmt.Println(fmt.Sprintf("\n%s - %s", name, quality))
	fmt.Println(fmt.Sprintf("\nTotal Download Size: %d MB", size))
	fmt.Println(fmt.Sprintf("\nOutput: %s\n", outputPath))

	c, err := utils.AskForConfirmation("» Proceed with download?")
	if err != nil {
		return err
	}
	if !c {
		fmt.Println("\nCancelled")
		os.Exit(0)
	}

	if utils.CheckIfPathExists(outputPath) {
		return fmt.Errorf("error: file '%s' already exists", outputPath)
	}

	fmt.Println("» Downloading media files...\n")
	dlErr := download(client, id)
	if dlErr != nil {
		return dlErr
	}

	fmt.Println("\n» Merging media files...")
	out, err := utils.MergeToMP4(fileListPath, outputPath)
	if err != nil {
		fmt.Println(string(out))
		return err
	}

	fmt.Println("» Cleaning up...")
	cleanErr := utils.CleanUp(tmpPath)
	if cleanErr != nil {
		fmt.Println(cleanErr)
	}

	fmt.Println("\nDownload completed!")

	return nil
}

func getClient() (*http.Client, error) {
	transport := &http.Transport{}

	if proxyFlag != "" {
		proxy, err := url.Parse(proxyFlag)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxy)
	}

	client := &http.Client{Transport: transport}

	return client, nil
}

func getVideo(client *http.Client, url string) (*types.Video, error) {
	video := &types.Video{}

	slug, err := parseUrl(url)
	if err != nil {
		return video, err
	}

	headers := make(map[string]string)
	headers["User-Agent"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"
	headers["Origin"] = "https://hanime.tv"

	body, err := utils.Request("GET", client, fmt.Sprintf("https://hw.hanime.tv/api/v8/video?id=%s", slug), headers, nil)
	if err != nil {
		return video, err
	}

	data, err := ioutil.ReadAll(body)
	if err := json.Unmarshal(data, &video); err != nil {
		return video, err
	}

	if err := body.Close(); err != nil {
		return video, err
	}

	return video, nil
}

func parseUrl(url string) (string, error) {
	re := regexp.MustCompile(`^https://hanime.tv/videos/hentai/(.*)$`)
	m := re.FindStringSubmatch(url)
	if len(m) > 1 && m[1] != "" {
		return m[1], nil
	}
	return "", fmt.Errorf("error: url '%s' is invalid", url)
}

func getStreamIndex(streams []types.Stream) (int, error) {
	if !utils.CheckIfInArray(videoQualities, qualityFlag) {
		return 0, fmt.Errorf("error: quality '%s' is invalid, possible values: 1080, 720, 480, 360", qualityFlag)
	}
	for k, v := range streams {
		if v.Height == qualityFlag {
			return k, nil
		}
	}
	return 0, nil
}

func setPaths(slug string, quality string, id int64) error {
	var outputName string
	if outputNameFlag != "" {
		if !strings.HasSuffix(outputNameFlag, ".mp4") {
			outputNameFlag += ".mp4"
		}
		outputName = outputNameFlag
	} else {
		outputName = fmt.Sprintf("%s-%s.mp4", slug, quality)
	}

	if outputPathFlag != "" {
		if !utils.CheckIfPathExists(outputPathFlag) {
			return fmt.Errorf("error: path '%s' does not exist", outputPathFlag)
		}
		outputPath = filepath.Join(outputPathFlag, outputName)
		tmpPath = filepath.Join(outputPathFlag, fmt.Sprintf("%s-%s-%d-tmp", slug, quality, id))
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		outputPath = filepath.Join(wd, outputName)
		tmpPath = filepath.Join(wd, fmt.Sprintf("%s-%s-%d-tmp", slug, quality, id))
	}

	fileListPath = filepath.Join(tmpPath, "filelist.txt")

	return nil
}

func download(client *http.Client, id int64) error {
	if utils.CheckIfPathExists(tmpPath) {
		fmt.Println(fmt.Errorf("error: cannot create temporary folder, path '%s' already exists", tmpPath))
		os.Exit(1)
	}

	err := utils.MakeDirectoryIfNotExists(tmpPath)
	if err != nil {
		return err
	}

	data, err := getM3U8(client, id)
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

		var segments []string
		for _, v := range mediapl.Segments {
			if v != nil {
				segments = append(segments, v.URI)
			}
		}

		err := createFileList(len(segments))
		if err != nil {
			return err
		}

		key, err := getKey(client, mediapl.Key.URI)
		if err != nil {
			return err
		}

		var iv []byte
		if mediapl.Key.IV == "" {
			iv = key
		} else {
			iv = []byte(mediapl.Key.IV)
		}

		var bar utils.Bar
		bar.New(0, int64(len(segments)), "█")

		wg := sync.WaitGroup{}
		for k, v := range segments {
			wg.Add(1)
			go func(index int, url string) {
				err := downloadTS(client, url, filepath.Join(tmpPath, fmt.Sprintf("%s.ts", strconv.Itoa(index))), key, iv)
				if err != nil {
					fmt.Printf("\n\n%s\n", err)
					cleanErr := utils.CleanUp(tmpPath)
					if cleanErr != nil {
						fmt.Println(cleanErr)
					}
					os.Exit(1)
				}
				wg.Done()
				bar.Next()
			}(k, v)
		}
		wg.Wait()

		bar.Finish()

		return nil
	}

	return nil
}

func getM3U8(client *http.Client, id int64) ([]byte, error) {
	body, err := utils.Request("GET", client, fmt.Sprintf("https://weeb.hanime.tv/weeb-api-cache/api/v8/m3u8s/%d", id), nil, nil)
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

func createFileList(count int) error {
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
	body, err := utils.Request("GET", client, url, nil, nil)
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

func downloadTS(client *http.Client, url string, path string, key []byte, iv[]byte) error {
	body, err := utils.Request("GET", client, url, nil, nil)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	decrypted, err := utils.Decrypt(data, key, iv)
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
