package password

type Service interface {
	Hash(password string) (string, error)
	Verify(password string, hash string) bool
}
