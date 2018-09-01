package pulli

func truncateString(s string, limit int) string {
	if len(s) < limit || limit < 0 {
		return s
	}

	return s[:limit]
}
