package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Merge creates and merges a list of configs into a single
// AWS session allowing the session to assume multiple roles.
func Merge(roles []string, s *session.Session) []*aws.Config {
	configs := []*aws.Config{aws.NewConfig()}

	for _, r := range roles {
		if len(r) == 0 {
			continue
		}

		creds := stscreds.NewCredentials(s, r)

		configs = append(configs, aws.NewConfig().WithCredentials(creds))
	}

	s.Config.MergeIn(configs...)

	return configs
}
