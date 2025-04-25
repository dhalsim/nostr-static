package types

type Layout struct {
	Color string `yaml:"color"`
	Logo  string `yaml:"logo"`
	Title string `yaml:"title"`
}

type Config struct {
	Relays   []string `yaml:"relays"`
	Articles []string `yaml:"articles"`
	Layout   Layout   `yaml:"layout"`
}
