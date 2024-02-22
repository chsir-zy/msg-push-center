package config

var JWT_KEY string

type Jwt struct {
	Key string `json:"key"`
}

func (j *Jwt) Get() string {
	return j.Key
}
