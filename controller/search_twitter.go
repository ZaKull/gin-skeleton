package controller

import (
	"encoding/json"
	"fmt"
	"github.com/cvcio/twitter"
	"github.com/gin-gonic/gin"
	"github.com/hyperjiang/gin-skeleton/config"
	"github.com/hyperjiang/gin-skeleton/model"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type SearchTwitterController struct{}

type SearchData struct {
	Text        string `form:"text" json:"text" xml:"text"  binding:"required"`
	Lang        string `form:"lang" json:"lang" xml:"lang"`
	SortOrder   string `form:"sort_order" json:"sort_order" xml:"sort_order"`
	TweetFields string `form:"tweet_fields" json:"tweet_fields" xml:"tweet_fields"`
	UserFields  string `form:"user_fields" json:"user_fields" xml:"user_fields"`
	Expansions  string `form:"expansions" json:"expansions" xml:"expansions"`
	PlaceFields string `form:"place_fields" json:"place_fields" xml:"place_fields"`
	PollFields  string `form:"poll_fields" json:"poll_fields" xml:"poll_fields"`
	MediaFields string `form:"media_fields" json:"media_fields" xml:"media_fields"`
}

// GetVersion version json
func (ctrl *SearchTwitterController) GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": config.Server.Version,
	})
}

func (ctrl *SearchTwitterController) Search(c *gin.Context) {
	var jsonData SearchData
	var search model.Search

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if jsonData.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text is empty"})
		return
	}

	search.Text = jsonData.Text

	search.GetFirstByText()

	if search.ID > 0 {
		c.JSON(http.StatusOK, gin.H{"data": search})
		return
	} else {
		search.CreateSearch()

		if jsonData.Lang == "" {
			jsonData.Lang = ""
		}

		api, err := twitter.NewTwitter(config.Twitter.ClientID, config.Twitter.ClientSecret)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		va := url.Values{}
		va.Add("query", jsonData.Text)
		if jsonData.SortOrder != "" {
			va.Add("sort_order", jsonData.SortOrder)
		}

		if jsonData.TweetFields != "" {
			va.Add("tweet.fields", jsonData.TweetFields)
		}

		if jsonData.UserFields != "" {
			va.Add("user.fields", jsonData.UserFields)
		}

		if jsonData.Expansions != "" {
			va.Add("expansions", jsonData.Expansions)
		}

		if jsonData.PlaceFields != "" {
			va.Add("place.fields", jsonData.PlaceFields)
		}

		if jsonData.PollFields != "" {
			va.Add("poll.fields", jsonData.PollFields)
		}

		if jsonData.MediaFields != "" {
			va.Add("media.fields", jsonData.MediaFields)
		}

		va.Add("max_results", "100")

		searchD, errs := api.GetTweetsSearchRecent(va)

		for {
			select {
			case r, ok := <-searchD:
				if !ok {
					searchD = nil
					break
				}
				b, err1 := json.Marshal(r.Data)

				if err1 != nil {
					panic(err1)
				}

				var data []*twitter.Tweet
				json.Unmarshal(b, &data)

				bu, err1u := json.Marshal(r.Includes.Users)

				if err1u != nil {
					panic(err1u)
				}

				var datau []*twitter.User
				json.Unmarshal(bu, &datau)

				twitterCheckAndInsert(search, data)
				twitterAuthorInsert(datau)
			case e, ok := <-errs:
				if !ok {
					errs = nil
					break
				}
				fmt.Sprintf("Twitter API Error: %v", e)
			}

			if searchD == nil && errs == nil {
				break
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"version": config.Server.Version,
		})
	}
}

func twitterCheckAndInsert(id model.Search, data []*twitter.Tweet) model.TwitterTweetMetrics {
	var TwitterDomains model.TwitterDomains
	var TwitterEntities model.TwitterEntities
	var TwitterDomainEntities model.TwitterDomainsEntities
	var TwitterTweetMetrics model.TwitterTweetMetrics
	var TwitterTweetTexts model.TwitterTweetTexts
	var TwitterTweets model.TwitterTweets
	var TwitterTweetTextNormalized model.TwitterTweetTextNormalized
	var TwitterTweetTextSentiments model.TwitterTweetTextsSentiments

	for _, v := range data {

		valTweetID, _ := strconv.ParseUint(v.ID, 10, 64)
		valAuthorID, _ := strconv.ParseUint(v.AuthorID, 10, 64)
		valConversationID, _ := strconv.ParseUint(v.ConversationID, 10, 64)
		valCreatedAt, _ := time.Parse("2006-01-02T15:04:05.000Z", v.CreatedAt)

		if v.ContextAnnotations != nil {
			for _, d := range v.ContextAnnotations {
				valDomain, _ := strconv.ParseUint(d.Domain.ID, 10, 64)
				if !checkForValueTwitterDomain(valDomain, TwitterDomains) {
					TwitterDomains = append(TwitterDomains, model.TwitterDomain{
						ID:          valDomain,
						Name:        d.Domain.Name,
						Description: d.Domain.Description,
					})
				}

				valEntity, _ := strconv.ParseUint(d.Entity.ID, 10, 64)
				if !checkForValueTwitterEntity(valEntity, TwitterEntities) {
					TwitterEntities = append(TwitterEntities, model.TwitterEntity{
						ID:          valEntity,
						Name:        d.Entity.Name,
						Description: d.Entity.Description,
					})
				}

				TwitterDomainEntities = append(TwitterDomainEntities, model.TwitterDomainEntity{
					TwitterDomainID: valDomain,
					TwitterEntityID: valEntity,
					TwitterTweetID:  valTweetID,
				})
			}
		}

		if v.PublicMetrics != nil {
			TwitterTweetMetrics = append(TwitterTweetMetrics, model.TwitterTweetMetric{
				TwitterTweetID: valTweetID,
				ReplyCount:     v.PublicMetrics.Replies,
				RetweetCount:   v.PublicMetrics.Retweets,
				LikeCount:      v.PublicMetrics.Likes,
				QuoteCount:     v.PublicMetrics.Quotes,
			})
		}

		TwitterTweetTexts = append(TwitterTweetTexts, model.TwitterTweetText{
			TwitterTweetID: valTweetID,
			Text:           v.Text,
		})

		NormalizeData := model.NormalizeData(v.Text)

		TwitterTweetTextNormalized = append(TwitterTweetTextNormalized, model.TwitterTweetTextNormalize{
			TwitterTweetID: valTweetID,
			NormalizeText:  NormalizeData,
			CreatedAt:      valCreatedAt,
		})

		AnalyzeSentiment, _ := model.AnalyzeSentiment(NormalizeData)

		TwitterTweetTextSentiments = append(TwitterTweetTextSentiments, model.TwitterTweetTextsSentiment{
			TwitterTweetID: valTweetID,
			Magnitude:      AnalyzeSentiment.Magnitude,
			Score:          AnalyzeSentiment.Score,
		})

		TwitterTweets = append(TwitterTweets, model.TwitterTweet{
			ID:             valTweetID,
			SearchID:       id.ID,
			TwitterUserID:  valAuthorID,
			ConversationID: valConversationID,
			Source:         v.Source,
			Lang:           v.Lang,
			CreatedAt:      valCreatedAt,
		})
	}

	if len(TwitterDomains) > 0 {
		TwitterDomains.CreateDomains()
	}

	if len(TwitterEntities) > 0 {
		TwitterEntities.CreateEntities()
	}

	if len(TwitterDomainEntities) > 0 {
		TwitterDomainEntities.CreateDomainsEntities()
	}

	if len(TwitterTweetMetrics) > 0 {
		TwitterTweetMetrics.CreateTweetMetrics()
	}

	if len(TwitterTweetTexts) > 0 {
		TwitterTweetTexts.CreateTweetTexts()
	}

	if len(TwitterTweetTextNormalized) > 0 {
		TwitterTweetTextNormalized.CreateSearchNormalized()
	}

	if len(TwitterTweetTextSentiments) > 0 {
		TwitterTweetTextSentiments.CreateTwitterTweetTextSentiments()
	}

	if len(TwitterTweets) > 0 {
		TwitterTweets.CreateTweets()
	}

	return TwitterTweetMetrics
}

func twitterAuthorInsert(data []*twitter.User) error {
	var TwitterUsers model.TwitterUsers
	var TwitterUserMetrics model.TwitterUserMetrics

	for _, v := range data {
		valUserID, _ := strconv.ParseUint(v.ID, 10, 64)
		valCreatedAt, _ := time.Parse("2006-01-02T15:04:05.000Z", v.CreatedAt)

		TwitterUserMetrics = append(TwitterUserMetrics, model.TwitterUserMetric{
			TwitterUserID:  valUserID,
			FollowersCount: v.PublicMetrics.Followers,
			FollowingCount: v.PublicMetrics.Following,
			TweetCount:     v.PublicMetrics.Tweets,
			ListedCount:    v.PublicMetrics.Listed,
		})

		TwitterUsers = append(TwitterUsers, model.TwitterUser{
			ID:          valUserID,
			Name:        v.Name,
			Username:    v.UserName,
			Location:    v.Location,
			Description: v.Description,
			Verified:    v.Verified,
			CreatedAt:   valCreatedAt,
		})
	}

	if len(TwitterUsers) > 0 {
		TwitterUsers.CreateUsers()
	}

	if len(TwitterUserMetrics) > 0 {
		TwitterUserMetrics.CreateUserMetrics()
	}

	return nil
}

func checkForValueTwitterDomain(idValue uint64, data []model.TwitterDomain) bool {
	for _, value := range data {
		if value.ID == idValue {
			return true
		}
	}
	return false
}

func checkForValueTwitterEntity(idValue uint64, data []model.TwitterEntity) bool {
	for _, value := range data {
		if value.ID == idValue {
			return true
		}
	}
	return false
}
