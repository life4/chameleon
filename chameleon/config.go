package chameleon

type CRepo struct {
	URL     string `toml:"url"`
	URLPath string `toml:"url_path"`
}

type Config struct {
	CachePath Path    `toml:"cache_path"`
	Repos     []CRepo `toml:"repo"`
}
