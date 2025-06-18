package containers

func ShortenID(id string) string {
	return shorten(id, 12)
}

func ShortenName(name string) string {
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	return shorten(name, 25)
}

func shorten(text string, amount int) string {
	if len(text) > amount {
		return text[:amount]
	}
	return text
}
