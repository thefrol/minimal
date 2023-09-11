package credentials

type Profile struct {
	AwsKey     string
	AwsSecret  string
	AwsSession string
}

func (p Profile) AWSCredentials() (string, string, string) {
	return p.AwsKey, p.AwsSecret, p.AwsSession
}
