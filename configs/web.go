package configs

type WebConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`

	TokenSecret     string `json:"token_secret"`
	TokenExpiration int    `json:"token_expiration"`

	CookiesField string `json:"cookies_field"`
}
