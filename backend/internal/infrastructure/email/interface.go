package email

type Sender interface {
	SendVerification(
		to string,
		code string,
	) error
}
