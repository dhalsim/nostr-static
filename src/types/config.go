package types

// Default weights for scoring
const (
	DefaultReactionWeight    = 1.0
	DefaultRepostWeight      = 2.0
	DefaultReplyWeight       = 1.5
	DefaultReportPenalty     = 5.0
	DefaultZapAmountWeight   = 0.1
	DefaultZapperCountWeight = 0.5
	DefaultZapAvgWeight      = 2.0
	DefaultZapTotalWeight    = 3.0
	DefaultMaxSats           = 10000000.0 // 0.1 BTC in sats
)

type Layout struct {
	Color      string `yaml:"color"`
	Logo       string `yaml:"logo"`
	FaviconDir string `yaml:"favicondir"`
	Title      string `yaml:"title"`
}

type Features struct {
	Comments     bool   `yaml:"comments"`
	NostrLinks   string `yaml:"nostr_links"`
	TagDiscovery bool   `yaml:"tag_discovery"`
}

type ScoringWeights struct {
	ReactionWeight    *float64 `yaml:"reaction_weight,omitempty" json:"reaction_weight,omitempty"`
	RepostWeight      *float64 `yaml:"repost_weight,omitempty" json:"repost_weight,omitempty"`
	ReplyWeight       *float64 `yaml:"reply_weight,omitempty" json:"reply_weight,omitempty"`
	ReportPenalty     *float64 `yaml:"report_penalty,omitempty" json:"report_penalty,omitempty"`
	ZapAmountWeight   *float64 `yaml:"zap_amount_weight,omitempty" json:"zap_amount_weight,omitempty"`
	ZapperCountWeight *float64 `yaml:"zapper_count_weight,omitempty" json:"zapper_count_weight,omitempty"`
	ZapAvgWeight      *float64 `yaml:"zap_avg_weight,omitempty" json:"zap_avg_weight,omitempty"`
	ZapTotalWeight    *float64 `yaml:"zap_total_weight,omitempty" json:"zap_total_weight,omitempty"`
	MaxSats           *float64 `yaml:"max_sats,omitempty" json:"max_sats,omitempty"`
}

// GetOrDefault returns the value if set, otherwise returns the default value
func (w *ScoringWeights) GetOrDefault(field *float64, defaultValue float64) float64 {
	if field == nil {
		return defaultValue
	}
	return *field
}

type TagDiscoverySettings struct {
	FetchCountPerTag     int            `yaml:"fetch_count_per_tag"`
	PopularArticlesCount int            `yaml:"popular_articles_count"`
	Weights              ScoringWeights `yaml:"weights"`
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
