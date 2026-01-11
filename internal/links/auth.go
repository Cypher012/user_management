package links

import "fmt"

func (l *Links) VerifyEmail(token string) string {
	return fmt.Sprintf("%s/auth/verify-email/%s", l.base, token)
}

func (l *Links) ResetPassword(token string) string {
	return fmt.Sprintf("%s/auth/reset-password/%s", l.base, token)
}
