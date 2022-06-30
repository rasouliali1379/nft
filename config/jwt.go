package config

// JWT struct
type JWT struct {
	HMACSecret  string `yaml:"jwt.hmacSecret" json:"hmac_secret"`
	RSASecret   string `yaml:"jwt.rsaSecret" json:"rsa_secret"`
	ECDSASecret string `yaml:"jwt.ecdsaSecret" json:"ecdsa_secret"`
}
