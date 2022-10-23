package model

import (
	language "cloud.google.com/go/language/apiv1"
	"context"
	"fmt"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"regexp"
	"strings"
	"time"
)

type TwitterTweetTextsSentiment struct {
	ID             uint64  `gorm:"primaryKey" json:"id"`
	TwitterTweetID uint64  `gorm:"index" json:"twitter_tweet_id"`
	Magnitude      float32 `json:"magnitude"`
	Score          float32 `json:"score"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type TwitterTweetTextsSentiments []TwitterTweetTextsSentiment

func (TwitterTweetTextsSentiment) TableName() string {
	return "twitter_tweet_texts_sentiments"
}

func (TwitterTweetTextsSentiments) TableName() string {
	return "twitter_tweet_texts_sentiments"
}

func (u *TwitterTweetTextsSentiment) CreateTwitterTweetTextSentiment() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterTweetTextsSentiments) CreateTwitterTweetTextSentiments() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func AnalyzeSentiment(text string) (*languagepb.Sentiment, error) {
	ctx := context.Background()

	// Creates a client.
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	// Detects the sentiment of the text.
	sentiment, err := client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})

	if err != nil {
		log.Fatalf("Failed to analyze text: %v", err)
	}

	fmt.Printf("Class sentiment: %v\n", sentiment.DocumentSentiment)

	return sentiment.DocumentSentiment, nil
}

func NormalizeData(text string) string {
	text = text + " "

	// replace to lower case
	text = strings.ToLower(text)

	// replace accent
	text = strings.ReplaceAll(text, "á", "a")
	text = strings.ReplaceAll(text, "é", "e")
	text = strings.ReplaceAll(text, "í", "i")
	text = strings.ReplaceAll(text, "ó", "o")
	text = strings.ReplaceAll(text, "ú", "u")
	text = strings.ReplaceAll(text, "Á", "a")
	text = strings.ReplaceAll(text, "É", "é")
	text = strings.ReplaceAll(text, "Í", "í")
	text = strings.ReplaceAll(text, "Ó", "ó")
	text = strings.ReplaceAll(text, "Ú", "u")
	text = strings.ReplaceAll(text, "ü", "u")

	//replace #name
	text = strings.ReplaceAll(text, "(", "")
	text = strings.ReplaceAll(text, ")", "")
	text = strings.ReplaceAll(text, ";", "")

	//replace url
	urls := regexp.MustCompile(`https?:\/\/.*?[\s+]`)
	for _, url := range urls.FindAllString(text, -1) {
		fmt.Println(url)
		text = strings.ReplaceAll(text, url, "")
	}

	// replace infinitive space, newline and tab
	space := regexp.MustCompile(`\s+`)
	text = space.ReplaceAllString(text, " ")

	return text
}

func NormalizeDataFull(text string) string {
	text = text + " "

	// replace to lower case
	text = strings.ToLower(text)

	// replace accent
	text = strings.ReplaceAll(text, "á", "a")
	text = strings.ReplaceAll(text, "é", "e")
	text = strings.ReplaceAll(text, "í", "i")
	text = strings.ReplaceAll(text, "ó", "o")
	text = strings.ReplaceAll(text, "ú", "u")
	text = strings.ReplaceAll(text, "Á", "a")
	text = strings.ReplaceAll(text, "É", "é")
	text = strings.ReplaceAll(text, "Í", "í")
	text = strings.ReplaceAll(text, "Ó", "ó")
	text = strings.ReplaceAll(text, "Ú", "ú")
	text = strings.ReplaceAll(text, "ü", "u")
	text = strings.ReplaceAll(text, "ñ", "n")
	text = strings.ReplaceAll(text, "Ñ", "n")

	//replace url
	urls := regexp.MustCompile(`https?:\/\/.*?[\s+]`)
	for _, url := range urls.FindAllString(text, -1) {
		text = strings.ReplaceAll(text, url, " ")
	}

	// replace username twitter
	usernames := regexp.MustCompile(`^@?(\w){1,15}$`)
	for _, username := range usernames.FindAllString(text, -1) {
		text = strings.ReplaceAll(text, username, " ")
	}

	// replace hashtag
	hashtags := regexp.MustCompile(`^#?(\w){1,15}$`)
	for _, hashtag := range hashtags.FindAllString(text, -1) {
		text = strings.ReplaceAll(text, hashtag, " ")
	}

	//replace #name
	text = strings.ReplaceAll(text, "(", " ")
	text = strings.ReplaceAll(text, ")", " ")
	text = strings.ReplaceAll(text, ";", " ")
	text = strings.ReplaceAll(text, ":", " ")
	text = strings.ReplaceAll(text, "!", " ")
	text = strings.ReplaceAll(text, "?", " ")
	text = strings.ReplaceAll(text, "¿", " ")
	text = strings.ReplaceAll(text, "¡", " ")
	text = strings.ReplaceAll(text, "°", " ")
	text = strings.ReplaceAll(text, "º", " ")
	text = strings.ReplaceAll(text, "ª", " ")
	text = strings.ReplaceAll(text, "´", " ")
	text = strings.ReplaceAll(text, "`", " ")
	text = strings.ReplaceAll(text, "¨", " ")
	text = strings.ReplaceAll(text, ",", " ")
	text = strings.ReplaceAll(text, ".", " ")

	// replace infinitive space, newline and tab
	space := regexp.MustCompile(`\s+`)
	text = space.ReplaceAllString(text, " ")

	return text
}
