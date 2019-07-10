package cfg

import (
	"log"
	"os"
)

const (
	url = "URL"
	env = "ENV"

	dbURI  = "DB_URI"
	dbName = "DB_NAME"

	oauthClientID     = "OAUTH_CLIENT_ID"
	oauthClientSecret = "OAUTH_CLIENT_SECRET"
	oauthRedirectURL  = "OAUTH_REDIRECT_URL"
	oauthSignOutURL   = "OAUTH_SIGN_OUT_URL"
	oauthAuthURL      = "OAUTH_AUTH_URL"
	oauthTokenURL     = "OAUTH_TOKEN_URL"
	oauthSessSign     = "OAUTH_SESSION_SIGN"

	sqsChangeQueue = "SQS_CHANGE_QUEUE"
	sqsAlarmQueue  = "SQS_ALARM_QUEUE"
	sqsS3Queue     = "SQS_S3_QUEUE"

	snsAlarmTopicArn = "SNS_ALARM_TOPIC_ARN"

	consoleLink = "CONSOLE_LINK"
)

// Get ...
var Get map[string]string

func init() {
	Get = map[string]string{}

	v, exists := os.LookupEnv(url)
	if !exists {
		log.Fatalf("environment variable %s required", url)
	}
	Get[url] = v

	v, exists = os.LookupEnv(env)
	if !exists {
		log.Fatalf("environment variable %s required", env)
	}
	Get[env] = v

	v, exists = os.LookupEnv(dbURI)
	if !exists {
		log.Fatalf("environment variable %s required", dbURI)
	}
	Get[dbURI] = v

	v, exists = os.LookupEnv(dbName)
	if !exists {
		log.Fatalf("environment variable %s required", dbName)
	}
	Get[dbName] = v

	v, exists = os.LookupEnv(oauthClientID)
	if !exists {
		log.Fatalf("environment variable %s required", oauthClientID)
	}
	Get[oauthClientID] = v

	v, exists = os.LookupEnv(oauthClientSecret)
	if !exists {
		log.Fatalf("environment variable %s required", oauthClientSecret)
	}
	Get[oauthClientSecret] = v

	v, exists = os.LookupEnv(oauthRedirectURL)
	if !exists {
		log.Fatalf("environment variable %s required", oauthRedirectURL)
	}
	Get[oauthRedirectURL] = v

	v, exists = os.LookupEnv(oauthSignOutURL)
	if !exists {
		log.Fatalf("environment variable %s required", oauthSignOutURL)
	}
	Get[oauthSignOutURL] = v

	v, exists = os.LookupEnv(oauthAuthURL)
	if !exists {
		log.Fatalf("environment variable %s required", oauthAuthURL)
	}
	Get[oauthAuthURL] = v

	v, exists = os.LookupEnv(oauthTokenURL)
	if !exists {
		log.Fatalf("environment variable %s required", oauthTokenURL)
	}
	Get[oauthTokenURL] = v

	v, exists = os.LookupEnv(oauthSessSign)
	if !exists {
		log.Fatalf("environment variable %s required", oauthSessSign)
	}
	Get[oauthSessSign] = v

	v, exists = os.LookupEnv(sqsChangeQueue)
	if !exists {
		log.Fatalf("environment variable %s required", sqsChangeQueue)
	}
	Get[sqsChangeQueue] = v

	v, exists = os.LookupEnv(sqsS3Queue)
	if !exists {
		log.Fatalf("environment variable %s required", sqsS3Queue)
	}
	Get[sqsS3Queue] = v

	v, exists = os.LookupEnv(sqsAlarmQueue)
	if !exists {
		log.Fatalf("environment variable %s required", sqsAlarmQueue)
	}
	Get[sqsAlarmQueue] = v

	v, exists = os.LookupEnv(snsAlarmTopicArn)
	if !exists {
		log.Fatalf("environment variable %s required", snsAlarmTopicArn)
	}
	Get[snsAlarmTopicArn] = v

	v, exists = os.LookupEnv(consoleLink)
	if !exists {
		log.Fatalf("environment variable %s required", consoleLink)
	}
	Get[consoleLink] = v
}
