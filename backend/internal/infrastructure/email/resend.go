package email

import (
	"fmt"
	"net/smtp"
)

type ResendService struct {
	host     string
	port     int
	email    string
	password string
}

func NewResend(
	host string,
	port int,
	email string,
	password string,
) *ResendService {

	return &ResendService{
		host:     host,
		port:     port,
		email:    email,
		password: password,
	}
}

func (s *ResendService) SendVerification(
	to string,
	code string,
) error {

	auth := smtp.PlainAuth(
		"",
		s.email,
		s.password,
		s.host,
	)

	subject := "Subject: Verify your email address\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	htmlBody := s.buildHTMLTemplate(code)

	msg := []byte(subject + mime + htmlBody)

	addr := fmt.Sprintf("%s:%v", s.host, s.port)

	return smtp.SendMail(
		addr,
		auth,
		s.email,
		[]string{to},
		msg,
	)
}

func (s *ResendService) buildHTMLTemplate(code string) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Email Verification</title>
	</head>
	<body style="margin: 0; padding: 0; background-color: #f4f6f8; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;">
		<table border="0" cellpadding="0" cellspacing="0" width="100%%" style="table-layout: fixed;">
			<tr>
				<td align="center" style="padding: 40px 0;">
					<table border="0" cellpadding="0" cellspacing="0" width="100%%" style="max-width: 500px; background-color: #ffffff; border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.05); overflow: hidden;">
						<tr>
							<td align="center" style="background-color: #4f46e5; padding: 30px 20px;">
								<h2 style="margin: 0; color: #ffffff; font-size: 24px; font-weight: 600; letter-spacing: 0.5px;">Confirm Your Email</h2>
							</td>
						</tr>
						<tr>
							<td style="padding: 40px 30px; text-align: center;">
								<p style="margin: 0 0 24px 0; color: #4b5563; font-size: 16px; line-height: 24px;">
									Thank you for registering! Please use the verification code below to activate your account:
								</p>
								
								<div style="background-color: #f3f4f6; border-radius: 6px; padding: 16px; margin: 24px 0; display: inline-block; letter-spacing: 6px; font-weight: bold;">
									<span style="font-size: 32px; color: #1f2937; font-family: monospace;">%s</span>
								</div>

								<p style="margin: 24px 0 0 0; color: #9ca3af; font-size: 14px;">
									The code will expire in <strong style="color: #6b7280;">15 minutes</strong>.
								</p>
							</td>
						</tr>
						<tr>
							<td style="background-color: #fafafa; padding: 20px 30px; text-align: center; border-top: 1px solid #f0f0f0;">
								<p style="margin: 0; color: #9ca3af; font-size: 12px; line-height: 18px;">
									If you did not request this email, you can safely ignore it.
								</p>
							</td>
						</tr>
					</table>
				</td>
			</tr>
		</table>
	</body>
	</html>
	`, code)
}
