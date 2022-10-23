package model

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type TwitterTweet struct {
	ID             uint64 `gorm:"primaryKey" json:"id"`
	SearchID       uint64 `gorm:"index" json:"search_id"`
	TwitterUserID  uint64 `gorm:"index" json:"twitter_user_id"`
	ConversationID uint64 `gorm:"index" json:"conversation_id"`
	Source         string `json:"source"`
	Lang           string `json:"lang"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type TwitterTweets []TwitterTweet

func (TwitterTweet) TableName() string {
	return "twitter_tweets"
}

func (TwitterTweets) TableName() string {
	return "twitter_tweets"
}

type TwitterTweetText struct {
	gorm.Model
	TwitterTweetID uint64 `json:"twitter_user_id"`
	Text           string `json:"text"`
}

type TwitterTweetTexts []TwitterTweetText

func (TwitterTweetText) TableName() string {
	return "twitter_tweet_text"
}

func (TwitterTweetTexts) TableName() string {
	return "twitter_tweet_text"
}

type TwitterDomain struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type TwitterDomains []TwitterDomain

func (TwitterDomain) TableName() string {
	return "twitter_domains"
}

func (TwitterDomains) TableName() string {
	return "twitter_domains"
}

type TwitterEntity struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type TwitterEntities []TwitterEntity

func (TwitterEntity) TableName() string {
	return "twitter_entities"
}

func (TwitterEntities) TableName() string {
	return "twitter_entities"
}

type TwitterDomainEntity struct {
	ID              uint64 `gorm:"primaryKey" json:"id"`
	TwitterTweetID  uint64 `json:"twitter_tweet_id"`
	TwitterDomainID uint64 `json:"twitter_tweet_id"`
	TwitterEntityID uint64 `json:"twitter_tweet_id"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type TwitterDomainsEntities []TwitterDomainEntity

func (TwitterDomainEntity) TableName() string {
	return "twitter_domains_entities"
}

func (TwitterDomainsEntities) TableName() string {
	return "twitter_domains_entities"
}

type TwitterTweetMetric struct {
	gorm.Model
	TwitterTweetID uint64 `json:"twitter_tweet_id"`
	LikeCount      int    `json:"like_count"`
	QuoteCount     int    `json:"quote_count"`
	ReplyCount     int    `json:"reply_count"`
	RetweetCount   int    `json:"retweet_count"`
}

type TwitterTweetMetrics []TwitterTweetMetric

func (TwitterTweetMetric) TableName() string {
	return "twitter_tweet_metrics"
}

func (TwitterTweetMetrics) TableName() string {
	return "twitter_tweet_metrics"
}

type TwitterUser struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	Name        string
	Username    string
	Location    string
	Description string
	Verified    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type TwitterUsers []TwitterUser

func (TwitterUser) TableName() string {
	return "twitter_users"
}

func (TwitterUsers) TableName() string {
	return "twitter_users"
}

type TwitterUserMetric struct {
	gorm.Model
	TwitterUserID  uint64 `json:"twitter_user_id"`
	FollowersCount int    `json:"followers_count"`
	FollowingCount int    `json:"following_count"`
	TweetCount     int    `json:"tweet_count"`
	ListedCount    int    `json:"listed_count"`
}

type TwitterUserMetrics []TwitterUserMetric

func (TwitterUserMetric) TableName() string {
	return "twitter_users_metrics"
}

func (TwitterUserMetrics) TableName() string {
	return "twitter_users_metrics"
}

type TwitterTweetTextNormalize struct {
	ID             uint64 `gorm:"primaryKey" json:"id"`
	TwitterTweetID uint64 `json:"twitter_tweet_id"`
	NormalizeText  string `json:"normalize_text"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type TwitterTweetTextNormalized []TwitterTweetTextNormalize

func (TwitterTweetTextNormalize) TableName() string {
	return "twitter_tweet_text_nomalized"
}

func (TwitterTweetTextNormalized) TableName() string {
	return "twitter_tweet_text_nomalized"
}

func (u *TwitterTweet) GetFirstByID(id string) error {
	db := DB().Where("id=?", id).First(u)

	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	} else if db.Error != nil {
		return db.Error
	}
	return nil
}

// Create a new user
func (u *TwitterTweet) Create() error {
	db := DB().Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterDomain) CreateDomain() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterDomains) CreateDomains() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterEntity) CreateEntity() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterEntities) CreateEntities() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterDomainEntity) CreateDomainEntity() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterDomainsEntities) CreateDomainsEntities() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterTweetMetric) CreateTweetMetric() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterTweetMetrics) CreateTweetMetrics() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterTweet) CreateTweet() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterTweets) CreateTweets() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterTweetText) CreateTweetText() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterTweetTexts) CreateTweetTexts() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterUser) CreateUser() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterUsers) CreateUsers() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterUserMetric) CreateUserMetric() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterUserMetrics) CreateUserMetrics() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterTweetTextNormalize) CreateSearchNormalize() error {
	u.NormalizeText = NormalizeData(u.NormalizeText)
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *TwitterTweetTextNormalized) CreateSearchNormalized() error {
	for i, v := range *u {
		(*u)[i].NormalizeText = NormalizeData(v.NormalizeText)
	}
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}
