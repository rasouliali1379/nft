package config

type Storage struct {
	Url         string  `yaml:"storage.url" required:"true"`
	Username    string  `yaml:"storage.username" required:"true"`
	Password    string  `yaml:"storage.password" required:"true"`
	SSL         bool    `yaml:"storage.ssl"`
	Buckets     Buckets `yaml:"storage.buckets" required:"true"`
	UrlExpInMin int     `yaml:"storage.urlExpInMin" required:"true"`
}

type Buckets struct {
	KYC string `yaml:"storage.buckets.kyc" required:"true"`
	NFT string `yaml:"storage.buckets.nft" required:"true"`
}
