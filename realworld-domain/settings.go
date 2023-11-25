package domain

type Settings struct {
	Security Security
}

type Security struct {
	TokenSecret string
}
