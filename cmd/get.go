package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/IchBinLeoon/hanime/types"
	"github.com/IchBinLeoon/hanime/utils"
	"github.com/spf13/cobra"
)

var tmpPath string
var outputPath string

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
var infoFlag bool

var getUsage = `Usage:
  hanime get <url> [flags]

Flags:
  -h, --help      help for get
  -q, --quality	  video quality (default 1080)
  -o, --output    custom output path
  -O, --Output    custom output name
  -p, --proxy     proxy url
  -i, --info      display video info
`

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.SetUsageTemplate(getUsage)
	getCmd.Flags().StringVarP(&qualityFlag, "quality", "q", "1080", "video quality")
	getCmd.Flags().StringVarP(&outputPathFlag, "output", "o", "", "custom output path")
	getCmd.Flags().StringVarP(&outputNameFlag, "Output", "O", "", "custom output name")
	getCmd.Flags().StringVarP(&proxyFlag, "proxy", "p", "", "proxy url")
	getCmd.Flags().BoolVarP(&infoFlag, "info", "i", false, "display video info")
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

	if infoFlag {
		fmt.Println(fmt.Sprintf("\nName:\t\t%s", name))
		fmt.Println(fmt.Sprintf("Quality:\t%sp", quality))
		fmt.Println(fmt.Sprintf("Views:\t\t%d", video.HentaiVideo.Views))
		fmt.Println(fmt.Sprintf("Interests:\t%d", video.HentaiVideo.Interests))
		fmt.Println(fmt.Sprintf("Brand:\t\t%s", video.HentaiVideo.Brand))
		fmt.Println(fmt.Sprintf("Likes:\t\t%d", video.HentaiVideo.Likes))
		fmt.Println(fmt.Sprintf("Dislikes:\t%d", video.HentaiVideo.Dislikes))
		fmt.Println(fmt.Sprintf("Downloads:\t%d", video.HentaiVideo.Downloads))
		fmt.Println(fmt.Sprintf("Monthly Rank:\t%d", video.HentaiVideo.MonthlyRank))
		fmt.Println(fmt.Sprintf("Created At:\t%s", time.Unix(video.HentaiVideo.CreatedAtUnix, 0)))
		fmt.Println(fmt.Sprintf("Released At:\t%s", time.Unix(video.HentaiVideo.ReleasedAtUnix, 0)))
	} else {
		fmt.Println(fmt.Sprintf("\n%s - %s", name, quality))
	}
	fmt.Println(fmt.Sprintf("\nTotal Download Size: %d MB", size))
	fmt.Println(fmt.Sprintf("\nOutput: %s\n", outputPath))

	c, err := utils.AskForConfirmation(":: Proceed with download?")
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

	downloader := utils.Downloader{Client: client}
	fmt.Printf("\n%s", name)
	dlErr := downloader.Download(fmt.Sprintf("https://weeb.hanime.tv/weeb-api-cache/api/v8/m3u8s/%d", id), tmpPath, outputPath)
	if dlErr != nil {
		return err
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

	return nil
}
