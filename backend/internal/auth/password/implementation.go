package password

type service struct{}

func New() Service {
	return &service{}
}

func (s *service) Hash(
	password string,
) (string, error) {
	return Hash(password)
}

func (s *service) Verify(
	password string,
	hash string,
) bool {
	return Compare(password, hash)
}
