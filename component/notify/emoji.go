package notify

func Emoji(status string) string {
	switch status {
	case "pending":
		return "hourglass"
	case "stopped", "error":
		return "thumbsdown"
	case "running", "deployed":
		return "thumbsup"
	default:
		return "confounded"
	}
}
