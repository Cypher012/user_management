package token

type TokenType string

const (
	VerifyEmailTokenType    TokenType = "verify-email"
	ForgetPasswordTokenType TokenType = "forget-password"
)
