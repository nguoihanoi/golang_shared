package languages

func GetText(inKey string, inLanguageCode string) string {
	if inKey == "" {
		return ""
	}
	temKey := "codes:" + inLanguageCode
	cacheData := languageCache.HGet(temKey, inKey)
	if strVal, ok := cacheData.(string); ok && strVal != "" {
		return strVal
	}
	languageCache.HSet(temKey, inKey, inKey)
	return inKey
}
