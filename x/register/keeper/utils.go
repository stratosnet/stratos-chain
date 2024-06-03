package keeper

func hasStringValue(items []string, item string) bool {
	for _, eachItem := range items {
		if len(item) > 0 && eachItem == item {
			return true
		}
	}
	return false
}
