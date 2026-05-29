package jwt

type JwtConfig struct {
	AccessSecret         string `env:"ACCESS_SECRET"`
	RefreshSecret        string `env:"REFRESH_SECRET"`
	Issuer               string `env:"ISSUER"`
	ExpireTimeoutAccess  string `env:"EXPIRE_TIMEOUT_ACCESS"`
	ExpireTimeoutRefresh string `env:"EXPIRE_TIMEOUT_REFRESH"`
}
