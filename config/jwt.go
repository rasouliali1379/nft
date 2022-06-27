package config

// JWT struct
type JWT struct {
	HMACSecret  string `yaml:"jwt.hmac_secret" json:"hmac_secret"`
	RSASecret   string `yaml:"jwt.rsa_secret" json:"rsa_secret"`
	ECDSASecret string `yaml:"jwt.ecdsa_secret" json:"ecdsa_secret"`
}
