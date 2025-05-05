package helpers

import "nostr-static/src/types"

func NameOrDisplayName(parsedProfile *types.ParsedProfile) string {
	if parsedProfile == nil {
		return ""
	}

	if parsedProfile.Name != "" {
		return parsedProfile.Name
	}

	return parsedProfile.DisplayName
}

func PictureOrFallback(parsedProfile *types.ParsedProfile, pubkey string) string {
	if parsedProfile == nil {
		return ""
	}

	if parsedProfile.Picture != "" {
		return parsedProfile.Picture
	}

	return "https://robohash.org/" + pubkey
}
