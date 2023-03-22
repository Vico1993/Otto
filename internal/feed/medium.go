package feed

var baseUrl string = "https://medium.com/feed/tag"

// Generate the list of medium article based on interested tag
func buildMediumFeedBasedOnTag() []string {
	mediumUrl := []string{}
	for _, tag := range tags {
		mediumUrl = append(mediumUrl, baseUrl+"/"+tag)
	}

	return mediumUrl
}
