package model

import "errors"

type Review struct {
	Id   int
	User User
	Text string
	StarCount int
}

func NewReview(id int, user User, text string, starCount int) (*Review, error) {
	if starCount < 1 || starCount > 5 {
		return nil, errors.New("star count must be between 1 and 5")
	}

	return &Review{
		Id: id,
		User: user,
		Text: text,
		StarCount: starCount,
	}, nil
}