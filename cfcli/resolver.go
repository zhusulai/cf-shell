package cfcli

import (
	"strings"
)

var resolverKeyword = []string{
	"APP_NAME",
	//"BUILDPACK",
	//"BUILDPACK_NAME",
	//"DOMAIN",
	//"FEATURE_NAME",
	"ORG",
	"ORG_NAME",
	//"PLUGIN-NAME",
	//"QUOTA",
	//"SECURITY_GROUP",
	//"SEGMENT_NAME",
	//"SERVICE_BROKER",
	"SERVICE_INSTANCE",
	//"SOURCE_APP",
	"SPACE",
	//"SPACE_QUOTA",
	//"SPACE_QUOTA_NAME",
	//"STACK_NAME",
}

func IsResolvableKeyWord(word string) bool {
	for _, keyword := range resolverKeyword {
		if strings.EqualFold(word, keyword) {
			return true
		}
	}
	return false
}

func ResolveKeyWord(word string) []string {
	switch word {
	case "APP_NAME":
		return context.cache.Apps()
	case "ORG", "ORG_NAME":
		return context.cache.Orgs()
	case "SERVICE_INSTANCE":
		return context.cache.Services()
	case "SPACE":
		return context.cache.Spaces()
	}
	return nil
}
