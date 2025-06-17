package http

func IsStaticFile(path string) bool {
	staticPaths := []string{"/javascript/", "/assets/"}
	for _, staticPath := range staticPaths {
		if len(path) >= len(staticPath) && path[:len(staticPath)] == staticPath {
			return true
		}
	}
	return false
}
