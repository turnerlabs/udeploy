package commit

import (
	"fmt"
	"testing"
)

func TestEnsureCommitsAreMutated(t *testing.T) {

	comment := "MYPRJCT-422: Feature or bug comments can be added here."

	expected := "<http://tickets.organization.com/browse/MYPRJCT-422|MYPRJCT-422>: Feature or bug comments can be added here."

	actual := mutate(comment, `[A-Z]{7}-\d*`, fmt.Sprintf("<http://tickets.organization.com/browse/%s|%s>", existingValueToken, existingValueToken))

	if actual != expected {
		t.Error("failed to insert link")
	}
}
