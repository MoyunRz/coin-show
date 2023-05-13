package config

var Cfg *Config

type Config struct {
	Vpn vpnConfig `toml:"vpn"`
}

type vpnConfig struct {
	Url string `toml:"url"`
}

func InitConfig() {

	Cfg = &Config{
		Vpn: vpnConfig{
			Url: "http://127.0.0.1:7890",
		},
	}

	// if _, err := toml.DecodeFile("config.toml", Cfg); err != nil {
	// 	log.Println("解析出错")
	// }
}
