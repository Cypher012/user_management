package email

type Service struct {
	sender    Sender
	templates *Template
}

func NewService(sender Sender, templates *Template) *Service {
	return &Service{
		sender:    sender,
		templates: templates,
	}
}

func (s *Service) SendWelcome(to string) error {
	html, err := render(s.templates.welcome, nil)
	if err != nil {
		return err
	}
	return s.sender.Send(to, "Welcome to my app", html)
}

func (s *Service) SendVerifyEmail(to, token string) error {
	html, err := render(s.templates.verify, map[string]string{
		"Token": token,
	})
	if err != nil {
		return err
	}
	return s.sender.Send(to, "Verify your email", html)
}

func (s *Service) SendResetPassword(to, token string) error {
	html, err := render(s.templates.reset, map[string]string{
		"Token": token,
	})
	if err != nil {
		return err
	}
	return s.sender.Send(to, "Reset your password", html)
}
