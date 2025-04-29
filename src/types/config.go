package types

type Layout struct {
	Color string `yaml:"color"`
	Logo  string `yaml:"logo"`
	Title string `yaml:"title"`
}

type Features struct {
	Comments bool `yaml:"comments"`
}

type Config struct {
	Relays   []string `yaml:"relays"`
	Articles []string `yaml:"articles"`
	Profiles []string `yaml:"profiles"`
	BlogURL  string   `yaml:"blog_url"`
	Layout   Layout   `yaml:"layout"`
	Features Features `yaml:"features"`
}
