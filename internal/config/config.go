package config

type Config struct {
	CoingeckoKey   string `yaml:"coingeckoKey" required:"true"`
	StorageType    string `yaml:"storage_type" default:"ram"`
	AuthorizerType string `yaml:"authorizer_type" default:"internal"`
	HttpConfig     `yaml:"http-config"`
	PostgresConfig `yaml:"postgres-storage"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}
type HttpConfig struct {
	JwtKey  string `yaml:"jwt_key" required:"true"`
	Address string `yaml:"address" env-default:"localhost:8080"`
}
