package notify

func Color(status string) string {
	switch status {
	case "pending":
		return "#ffc845"
	case "stopped", "error":
		return "#ff4c4c"
	case "running", "deployed":
		return "#34bf49"
	default:
		return "#52565e"
	}
}
