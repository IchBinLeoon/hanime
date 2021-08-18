package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/IchBinLeoon/hanime/cmd/utils"
	"github.com/spf13/cobra"
)

const apiVideo = "https://hw.hanime.tv/api/v8/video?id="
const apiM3U8 = "https://weeb.hanime.tv/weeb-api-cache/api/v8/m3u8s/"

var tmpPath string

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
  hanime get <urls> [flags]

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
	Short: "Download videos by url",
	Long:  "Download one or more videos from hanime.tv by url",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utils.CatchInterrupt(&tmpPath)
		if err := get(args); err != nil {
			fmt.Println(err)
			cleanErr := utils.CleanUp(tmpPath)
			if cleanErr != nil {
				fmt.Println(cleanErr)
			}
			os.Exit(1)
		}
	},
}

func get(urls []string) error {
	client, err := utils.DefaultClient(proxyFlag)
	if err != nil {
		return err
	}

	var videos []utils.Video
	for i, url := range urls {
		video, err := getVideo(client, url)
		if err != nil {
			return err
		}
		index, err := getStreamIndex(video.VideosManifest.Servers[0].Streams)
		if err != nil {
			return err
		}
		video.StreamIndex = index
		path, err := getOutputPath(video.HentaiVideo.Slug, video.VideosManifest.Servers[0].Streams[video.StreamIndex].Height)
		if err != nil {
			return err
		}
		if len(urls) > 1 {
			if outputNameFlag != "" || utils.CheckIfMultipleInArray(urls, url) {
				path = fmt.Sprintf("%s-%d.mp4", path[:len(path)-4], i)
			}
		}
		video.OutputPath = path
		videos = append(videos, *video)
	}

	fmt.Print("\n")
	for _, video := range videos {
		if infoFlag {
			fmt.Printf("Name:\t\t%s\n", video.HentaiVideo.Name)
			fmt.Printf("Quality:\t%s\np", video.VideosManifest.Servers[0].Streams[video.StreamIndex].Height)
			fmt.Printf("Views:\t\t%d\n", video.HentaiVideo.Views)
			fmt.Printf("Interests:\t%d\n", video.HentaiVideo.Interests)
			fmt.Printf("Brand:\t\t%s\n", video.HentaiVideo.Brand)
			fmt.Printf("Likes:\t\t%d\n", video.HentaiVideo.Likes)
			fmt.Printf("Dislikes:\t%d\n", video.HentaiVideo.Dislikes)
			fmt.Printf("Downloads:\t%d\n", video.HentaiVideo.Downloads)
			fmt.Printf("Monthly Rank:\t%d\n", video.HentaiVideo.MonthlyRank)
			fmt.Printf("Created At:\t%s\n", time.Unix(video.HentaiVideo.CreatedAtUnix, 0))
			fmt.Printf("Released At:\t%s\n", time.Unix(video.HentaiVideo.ReleasedAtUnix, 0))
			fmt.Printf("Output:\t\t%s\n\n", video.OutputPath)
		} else {
			fmt.Printf("%s - %s\n", video.HentaiVideo.Name, video.VideosManifest.Servers[0].Streams[video.StreamIndex].Height)
			fmt.Printf("%s\n\n", video.OutputPath)
		}
	}
	var size int64
	for _, video := range videos {
		size += video.VideosManifest.Servers[0].Streams[video.StreamIndex].Size
	}
	fmt.Printf("Total Download Size: %d MB\n\n", size)

	c, err := utils.AskForConfirmation(":: Proceed with download?")
	if err != nil {
		return err
	}
	if !c {
		fmt.Println("\nCancelled")
		os.Exit(0)
	}

	downloader := utils.Downloader{Client: client}
	for _, video := range videos {
		err := download(&downloader, &video)
		if err != nil {
			return err
		}
	}

	fmt.Println("\nDownload completed!")

	return nil
}

func download(downloader *utils.Downloader, video *utils.Video) error {
	fmt.Printf("\n%s", video.HentaiVideo.Name)
	if utils.CheckIfPathExists(video.OutputPath) {
		fmt.Printf("\nwarning: file '%s' already exists, skipping\n", video.OutputPath)
		return nil
	}
	tmpPath = fmt.Sprintf("%s-%d", video.OutputPath[:len(video.OutputPath)-4], video.VideosManifest.Servers[0].Streams[video.StreamIndex].ID)
	err := downloader.Download(fmt.Sprintf("%s%d", apiM3U8, video.VideosManifest.Servers[0].Streams[video.StreamIndex].ID), tmpPath, video.OutputPath)
	if err != nil {
		return err
	}
	return nil
}

func getVideo(client *http.Client, url string) (*utils.Video, error) {
	video := &utils.Video{}

	slug, err := parseUrl(url)
	if err != nil {
		return video, err
	}

	headers := make(map[string]string)
	headers["User-Agent"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"
	headers["Origin"] = "https://hanime.tv"

	body, err := utils.Request("GET", client, fmt.Sprintf("%s%s", apiVideo, slug), headers, nil)
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

func getStreamIndex(streams []utils.Stream) (int, error) {
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

func getOutputPath(slug string, quality string) (string, error) {
	var outputName string
	if outputNameFlag != "" {
		re := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
		outputNameFlag = re.ReplaceAllString(outputNameFlag, "")
		if !strings.HasSuffix(outputNameFlag, ".mp4") {
			outputNameFlag += ".mp4"
		}
		outputName = outputNameFlag
	} else {
		outputName = fmt.Sprintf("%s-%s.mp4", slug, quality)
	}

	var outputPath string
	if outputPathFlag != "" {
		if !utils.CheckIfPathExists(outputPathFlag) {
			return "", fmt.Errorf("error: path '%s' does not exist", outputPathFlag)
		}
		outputPath = filepath.Join(outputPathFlag, outputName)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		outputPath = filepath.Join(wd, outputName)
	}

	return outputPath, nil
}
