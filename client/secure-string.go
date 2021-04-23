package client

type SecurePassword string

func (s SecurePassword) String() string {
	return "****SECURE****"
}

func (s SecurePassword) Value() string {
	return string(s)
}
