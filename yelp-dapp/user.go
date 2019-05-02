package yelp_dapp

type User struct {
	PrivateKey string
	PublicKey  string
	UserName   string
	UserId     string
}

func NewUser() User {
	return User{}
}
