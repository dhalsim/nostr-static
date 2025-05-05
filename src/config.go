package main

import (
	"os"

	"gopkg.in/yaml.v3"

	"nostr-static/src/types"
)

func setDefaultWeights(weights *types.ScoringWeights) {
	if weights.ReactionWeight == nil {
		weights.ReactionWeight = &[]float64{types.DefaultReactionWeight}[0]
	}
	if weights.RepostWeight == nil {
		weights.RepostWeight = &[]float64{types.DefaultRepostWeight}[0]
	}
	if weights.ReplyWeight == nil {
		weights.ReplyWeight = &[]float64{types.DefaultReplyWeight}[0]
	}
	if weights.ReportPenalty == nil {
		weights.ReportPenalty = &[]float64{types.DefaultReportPenalty}[0]
	}
	if weights.ZapAmountWeight == nil {
		weights.ZapAmountWeight = &[]float64{types.DefaultZapAmountWeight}[0]
	}
	if weights.ZapperCountWeight == nil {
		weights.ZapperCountWeight = &[]float64{types.DefaultZapperCountWeight}[0]
	}
	if weights.ZapAvgWeight == nil {
		weights.ZapAvgWeight = &[]float64{types.DefaultZapAvgWeight}[0]
	}
	if weights.ZapTotalWeight == nil {
		weights.ZapTotalWeight = &[]float64{types.DefaultZapTotalWeight}[0]
	}
	if weights.MaxSats == nil {
		weights.MaxSats = &[]float64{types.DefaultMaxSats}[0]
	}
}

func LoadConfig(path string) (*types.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config types.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Set default theme if not specified
	if config.Layout.Color == "" {
		config.Layout.Color = "light"
	}

	// Set default weights if tag discovery is enabled
	if config.Features.TagDiscovery {
		setDefaultWeights(&config.Settings.TagDiscovery.Weights)
	}

	return &config, nil
}
