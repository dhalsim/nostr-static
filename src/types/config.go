package types

type Layout struct {
	Color string `yaml:"color"`
	Logo  string `yaml:"logo"`
}

type Config struct {
	Relays     []string `yaml:"relays"`
	ArticleIDs []string `yaml:"article_ids"`
	Layout     Layout   `yaml:"layout"`
}
