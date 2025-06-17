package containers

func ShortenID(id string) string {
	return shortenWithAmount(id, 12)
}

func ShortenName(name string) string {
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	return shortenWithAmount(name, 25)
}

func shortenWithAmount(text string, amount int) string {
	if len(text) > amount {
		return text[:amount]
	}
	return text
}
