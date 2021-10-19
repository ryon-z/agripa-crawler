package crawler

import (
	"encoding/json"
	"fmt"
	"media_crawling/config"
	"media_crawling/models"
	"media_crawling/util"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

func getChannels(db *gorm.DB) []models.Channel {
	var channels []models.Channel

	tableName := models.Channel{}.TableName()
	sqlQuery := fmt.Sprintf("SELECT * FROM %s", tableName)
	db.Raw(sqlQuery).Scan(&channels)

	return channels
}

func requestYoutube(pageToken string, channelID string, order string, q string, usePeriod bool) YoutubeResponse {
	requestURL := "https://www.googleapis.com/youtube/v3/search"
	completedURL := fmt.Sprintf("%s?key=%s&part=snippet&maxResults=10&order=%s&type=video&regionCode=kr",
		requestURL, config.Secret["youtube:apikey"], order)

	if channelID != "" {
		completedURL += fmt.Sprintf("&channelId=%s", channelID)
	}

	if pageToken != "init" {
		completedURL += fmt.Sprintf("&pageToken=%s", pageToken)
	}

	if q != "" {
		completedURL += fmt.Sprintf("&q=%s", q)
	}

	if usePeriod {
		timeMonthsAgo := util.GetTimeForMonthsAgo(4)
		timeMonthsAgo = strings.ReplaceAll(timeMonthsAgo, "+09:00", "Z")
		completedURL += fmt.Sprintf("&publishedAfter=%s", timeMonthsAgo)
	}

	var headers map[string]string
	headers = make(map[string]string)

	result := Request(completedURL, headers)
	fmt.Println(result)
	videos := YoutubeResponse{}
	json.Unmarshal([]byte(result), &videos)

	return videos
}

func refineVideos(videos YoutubeResponse, query string) (map[string]interface{}, []models.Youtube) {
	var meta map[string]interface{}
	meta = make(map[string]interface{})
	meta["TotalResults"] = videos.PageInfo.TotalResults
	meta["ResultsPerPage"] = videos.PageInfo.ResultsPerPage
	meta["NextPageToken"] = videos.NextPageToken

	var refinedVideos []models.Youtube

	uselessWords := GetUselessWords()

	for _, video := range videos.Items {
		// video.ID.Kind가 youtube#video이 아닐 경우(ex. youtube#channel)
		if video.ID.VideoID == "" {
			continue
		}

		// video.Snippet.PublishedAt는 utc 시간이라 한국 시간으로 변환
		publishedAt := util.GetKoreanDateTime(video.Snippet.PublishedAt, "RFC3339", true, DatetimeFormat)
		title := util.ReplaceString(video.Snippet.Title, uselessWords)
		description := util.ReplaceString(video.Snippet.Description, uselessWords)
		channelTitle := util.ReplaceString(video.Snippet.ChannelTitle, uselessWords)
		refinedVideos = append(refinedVideos, models.Youtube{
			Query:        query,
			VideoID:      video.ID.VideoID,
			ChannelID:    video.Snippet.ChannelID,
			ChannelTitle: channelTitle,
			Title:        title,
			Description:  description,
			ThumbnailURL: video.Snippet.Thumbnails.Medium.URL,
			PublishedAt:  publishedAt,
		})
	}

	return meta, refinedVideos
}

// CrawlYoutube : youtube data API를 이용하여 youtube
func CrawlYoutube(firstCodeNum int, logFilePath string) {
	fpLog, logger := util.GetFileLogger(logFilePath)
	defer util.CloseFileLogger(fpLog)

	db := GetDB("operation")
	defer CloseDB(db)

	youtubeQueries := models.GetQueries(db, "youtube", firstCodeNum)

	logger.Println("유튜브 수집 정상시작")
	for _, youtubeQuery := range youtubeQueries {
		logger.Printf("검색어 \"%s\"에 대한 수집 시작\n", youtubeQuery)
		videos := requestYoutube("init", "", "viewCount", youtubeQuery, true)

		// 정제
		_, refinedVideos := refineVideos(videos, youtubeQuery)

		// DB에 업로드
		data := make([]interface{}, len(refinedVideos))
		for index, row := range refinedVideos {
			data[index] = row
		}
		ImportDataToDB(db, "youtube", "", data)
		logger.Printf("검색어 \"%s\"에 대한 수집 종료\n", youtubeQuery)

		time.Sleep(time.Second * 3)
	}

	logger.Println("유튜브 수집 정상종료")
}

// CrawlYoutubeByChannel : 채널 별 유튜브 수집
func CrawlYoutubeByChannel(logFilePath string) {
	fpLog, logger := util.GetFileLogger(logFilePath)
	defer util.CloseFileLogger(fpLog)

	db := GetDB("operation")
	defer CloseDB(db)

	channels := getChannels(db)
	commonQuery := "000공통000"

	logger.Println("유튜브 수집 정상시작")
	for _, channel := range channels {
		channelID := channel.ID
		logger.Printf("채널 \"%s\"에 대한 수집 시작\n", channelID)
		videos := requestYoutube("init", channelID, "date", "", false)

		// 정제
		_, refinedVideos := refineVideos(videos, commonQuery)

		// DB에 업로드
		data := make([]interface{}, len(refinedVideos))
		for index, row := range refinedVideos {
			data[index] = row
		}
		ImportDataToDB(db, "youtube", "", data)
		logger.Printf("채널 \"%s\"에 대한 수집 종료\n", channelID)

		time.Sleep(time.Second * 3)
	}

	logger.Println("유튜브 수집 정상종료")
}
