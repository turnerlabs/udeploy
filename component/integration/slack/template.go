package slack

import (
	"fmt"
)

func Template(subject, text, color, linkText, linkURL string) string {
	return fmt.Sprintf(`
           {
               "text": "%s",
               "attachments": [
                   {
                       "fallback": "%s",
                       "text": "%s",
                       "color": "%s"
                   },
                   {
                        "title": "%s",
                        "title_link": "%s"
                   }
               ]
           }`, subject, text, text, color, linkText, linkURL)
}
