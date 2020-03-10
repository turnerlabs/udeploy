package commit

func formatForDisplay(text string) string {
	//replace all new lines with <br /> tags
	text = strings.Replace(text, "\n", "<br />", -1)

	return text
}
