package yelp_dapp

type Business struct {
	BusinessName string
	Location     string
	Tag          string
	Id           string
}

func NewBusiness(businessName string, location string, tag string, id string) Business {
	return Business{
		BusinessName: businessName,
		Location:     location,
		Tag:          tag,
		Id:           id,
	}
}
