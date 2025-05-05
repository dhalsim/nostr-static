package types

type Layout struct {
	Color string `yaml:"color"`
	Logo  string `yaml:"logo"`
	Title string `yaml:"title"`
}

type Features struct {
	Comments     bool   `yaml:"comments"`
	NostrLinks   string `yaml:"nostr_links"`
	TagDiscovery bool   `yaml:"tag_discovery"`
}

type TagDiscoverySettings struct {
	FetchCountPerTag     int `yaml:"fetch_count_per_tag"`
	PopularArticlesCount int `yaml:"popular_articles_count"`
}

type Settings struct {
	TagDiscovery TagDiscoverySettings `yaml:"tag_discovery"`
}

type Config struct {
	Relays   []string `yaml:"relays"`
	Articles []string `yaml:"articles"`
	Profiles []string `yaml:"profiles"`
	BlogURL  string   `yaml:"blog_url"`
	Layout   Layout   `yaml:"layout"`
	Features Features `yaml:"features"`
	Settings Settings `yaml:"settings"`
}
