package utils

func ParseDuration(duration string) string {
	switch duration {
	case "-1h":
		return "1 hour"
	case "-2h":
		return "2 hours"
	case "-3h":
		return "3 hours"
	default:
		return "1 hour"
	}
}
