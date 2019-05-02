package yelp_dapp

import (
	"encoding/json"
	"log"
)

const RatingLowerBound float32 = 0
const RatingUpperBound float32 = 5

type Review struct {
	UserId     string  `json:"user_id"`
	BusinessId string  `json:"business_id"`
	Rating     float32 `json:"rating"`
	Comment    string  `json:"comment"`
}

func NewReview(userId string, businessId string, rating float32, comment string) Review {
	if rating < RatingLowerBound {
		rating = 0
	} else if rating > RatingUpperBound {
		rating = 5
	}
	return Review{
		UserId:     userId,
		BusinessId: businessId,
		Rating:     rating,
		Comment:    comment,
	}
}

func (review *Review) ToJson() string {
	res, err := json.Marshal(review)
	if err != nil {
		log.Fatal("cannot convert review to json string", err)
	}
	return string(res)
}
