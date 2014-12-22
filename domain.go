package clickcounteriframe

import "time"

type EntityInterface interface {
}

type Domain struct {
	EntityInterface
	Id              int64
	Name            string
	Redirect        string
	LandingPageJson []byte
	LandingPage     *LandingPage
	Created         *time.Time
	Updated         *time.Time
}

type LandingPage struct {
	Locale          string `json:"locale"`
	Title           string `json:"title"`
	About           string `json:"about"`
	MicroDonation   string `json:"microDonation"`
	LearnMore       string `json:"learnMore"`
	GetYourOwn      string `json:"getYourOwn"`
	TellYourFriends string `json:"tellYourFriends"`
	Tweet           string `json:"tweet"`
	TweetEncoded    string
	Imprint         string `json:"imprint"`
}
