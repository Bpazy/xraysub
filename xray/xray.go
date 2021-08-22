package xray

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Bpazy/xraysub/util"
	"github.com/go-resty/resty/v2"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Author struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Uploader struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Asset struct {
	Url                string      `json:"url"`
	Id                 int         `json:"id"`
	NodeId             string      `json:"node_id"`
	Name               string      `json:"name"`
	Label              interface{} `json:"label"`
	Uploader           *Uploader   `json:"uploader"`
	ContentType        string      `json:"content_type"`
	State              string      `json:"state"`
	Size               int         `json:"size"`
	DownloadCount      int         `json:"download_count"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	BrowserDownloadUrl string      `json:"browser_download_url"`
}

type Reactions struct {
	Url        string `json:"url"`
	TotalCount int    `json:"total_count"`
	Field3     int    `json:"+1"`
	Field4     int    `json:"-1"`
	Laugh      int    `json:"laugh"`
	Hooray     int    `json:"hooray"`
	Confused   int    `json:"confused"`
	Heart      int    `json:"heart"`
	Rocket     int    `json:"rocket"`
	Eyes       int    `json:"eyes"`
}

type GithubLatestRelease struct {
	Url             string     `json:"url"`
	AssetsUrl       string     `json:"assets_url"`
	UploadUrl       string     `json:"upload_url"`
	HtmlUrl         string     `json:"html_url"`
	Id              int        `json:"id"`
	Author          *Author    `json:"author"`
	NodeId          string     `json:"node_id"`
	TagName         string     `json:"tag_name"`
	TargetCommitish string     `json:"target_commitish"`
	Name            string     `json:"name"`
	Draft           bool       `json:"draft"`
	Prerelease      bool       `json:"prerelease"`
	CreatedAt       time.Time  `json:"created_at"`
	PublishedAt     time.Time  `json:"published_at"`
	Assets          []*Asset   `json:"assets"`
	TarballUrl      string     `json:"tarball_url"`
	ZipballUrl      string     `json:"zipball_url"`
	Body            string     `json:"body"`
	Reactions       *Reactions `json:"reactions"`
}

func NewXrayDownloadCmdRun() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		downloadUrl, fileName, err := getDownloadUrl()
		if err != nil {
			util.CheckErr(err)
		}

		download(err, downloadUrl, fileName)
	}
}

func download(err error, downloadUrl string, fileName string) {
	client := resty.New()
	client.SetProxy("http://127.0.0.1:10809")
	client.SetTimeout(30 * time.Second)
	res, err := client.R().
		SetDoNotParseResponse(true).
		Get(downloadUrl)
	if err != nil {
		util.CheckErr(fmt.Errorf("download error: %w", err))
	}
	defer util.Closeq(res.RawResponse.Body)

	outFile, err := os.Create(filepath.Clean(fileName))
	if err != nil {
		util.CheckErr(fmt.Errorf("create file error: %w", err))
	}
	defer util.Closeq(outFile)

	bar := getDownloadProgressBar(res.RawResponse.ContentLength)
	// io.Copy reads maximum 32kb size, it is perfect for large file download too
	_, err = io.Copy(io.MultiWriter(outFile, bar), res.RawResponse.Body)
	if err != nil {
		util.CheckErr(fmt.Errorf("io copy error: %w", err))
	}

	fmt.Println()
	fmt.Printf("The xray-core is saved %s\n", fileName)
}

func getDownloadProgressBar(maxLength int64) *progressbar.ProgressBar {
	bar := progressbar.NewOptions64(
		maxLength,
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			_, _ = fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	_ = bar.RenderBlank()
	return bar
}

func getDownloadUrl() (string, string, error) {
	client := resty.New()
	client.SetTimeout(5 * time.Second)
	start := time.Now()
	fmt.Println("Requesting /XTLS/Xray-core latest release")
	log.Infof("querying /XTLS/Xray-core latest release")
	res, err := client.R().Get("https://api.github.com/repos/XTLS/Xray-core/releases/latest")
	util.CheckErr(err)
	log.Infof("querying /XTLS/Xray-core latest release cost %dms", time.Since(start).Milliseconds())

	r := &GithubLatestRelease{}
	if err = json.Unmarshal(res.Body(), r); err != nil {
		util.CheckErr(err)
	}

	t := fmt.Sprintf("Xray-%s-64.zip", runtime.GOOS)
	var downloadUrl *string
	for _, asset := range r.Assets {
		if asset.Name == t {
			downloadUrl = &asset.BrowserDownloadUrl
		}
	}
	if downloadUrl == nil {
		return "", t, errors.New("no xray-core for the current platform")
	}

	fmt.Printf("Got latest Xray-core: %s\n", *downloadUrl)
	log.Infof("got download url: %s", *downloadUrl)
	return *downloadUrl, t, nil
}
