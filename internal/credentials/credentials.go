package credentials

type Credentials struct {
	defaultProfileName string
	Profiles           map[string]Profile
}

func (c Credentials) DefaultProfile() Profile {
	return c.Profile(c.defaultProfileName)
}

func (c Credentials) Profile(name string) Profile {
	return c.Profiles[name]
}
