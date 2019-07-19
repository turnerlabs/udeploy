package commit

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/turnerlabs/udeploy/model"

	"github.com/turnerlabs/udeploy/component/integration/github"
)

// Change ...
type Change struct {
	Tag string `json:"tag"`

	Message string `json:"message"`

	URL string `json:"url"`
	SHA string `json:"sha"`
}

// BuildRelease ...
func BuildRelease(org, repo, targetTag, currentTag, url, accessToken string, maxCommits int, commitConfig model.CommitConfig) ([]Change, error) {

	tags, err := github.GetTags(org, repo, url, accessToken)
	if err != nil {
		return []Change{}, err
	}

	tag, found := tags[targetTag]
	if !found {
		return []Change{}, fmt.Errorf("tag %s not found", targetTag)
	}

	commit, err := github.GetCommit(org, repo, tag.Commit.SHA, url, accessToken)
	if err != nil {
		return []Change{}, err
	}

	changes := []Change{}
	for x := 1; x <= commitConfig.Limit; x++ {

		if isCommitImportant(commitConfig.Filter, commit.Commit) {
			commit.Commit.Message = mutate(commit.Commit.Message, commitConfig.ExistingValue, commitConfig.NewValue)

			commit.Commit.Message = tokenize(commit, targetTag)

			if nextTag, err := getCommitTag(commit.SHA, tags); err == nil {
				tag = nextTag
			}

			changes = append(changes, Change{
				Tag: tag.Name,

				URL: buildLink(commit, org, repo),
				SHA: commit.SHA,

				Message: commit.Commit.Message,
			})
		} else {
			x--
		}

		if isRange(currentTag, targetTag) {
			if reachedCurrentTag(commit.Parents[0].SHA, currentTag, tags) {
				break
			}
		}

		commit, err = github.GetCommit(org, repo, commit.Parents[0].SHA, url, accessToken)
		if err != nil {
			return changes, err
		}
	}

	return changes, nil
}

func getCommitTag(SHA string, tags map[string]github.Tag) (github.Tag, error) {
	for _, t := range tags {
		if t.Commit.SHA == SHA {
			return t, nil
		}
	}

	return github.Tag{}, errors.New("not tagged")
}

func isRange(start, end string) bool {
	return start != end
}

func buildLink(commit github.Commit, org, repo string) string {
	return fmt.Sprintf("https://github.com/%s/%s/commit/%s", org, repo, commit.SHA)
}

func tokenize(commit github.Commit, tag string) string {

	commit.Commit.Message = strings.Replace(commit.Commit.Message, "{VERSION}", tag, -1)

	commit.Commit.Message = strings.Replace(commit.Commit.Message, "{COMMIT_HASH}", commit.SHA, -1)

	return commit.Commit.Message
}

func isCommitImportant(filter string, commit github.CommitCommit) bool {
	if len(filter) == 0 {
		return true
	}

	re := regexp.MustCompile(filter)

	return re.MatchString(commit.Message)
}

func reachedCurrentTag(SHA, tag string, tags map[string]github.Tag) bool {
	for _, t := range tags {
		if t.Name == tag && t.Commit.SHA == SHA {
			return true
		}
	}

	return false
}
