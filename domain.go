package clickcounteriframe

import (
	"database/sql"
	"time"
)

type EntityInterface interface {
}

type Domain struct {
	EntityInterface
	Id              int64
	Name            string
	Redirect        sql.NullString
	LandingPageJson []byte
	LandingPage     *LandingPage
	Created         *time.Time
	Updated         *time.Time
}

type LandingPage struct {
	DefaultLocale string `json:"defaultLocale"`
	Strings       map[string]*LandingPageText
}

type LandingPageText struct {
	Locale          string `json:"-"`
	Title           string `json:"title"`
	About           string `json:"about"`
	LearnMore       string `json:"learnMore"`
	GetYourOwn      string `json:"getYourOwn"`
	TellYourFriends string `json:"tellYourFriends"`
	Tweet           string `json:"tweet"`
	TweetEncoded    string `json:"-"`
	Imprint         string `json:"imprint"`
}
