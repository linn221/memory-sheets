package models

type Session struct {
	message string
}

func (s *Session) SetMessage(message string) {
	s.message = message
}

func (s *Session) HasMessage() bool {
	return s.message != ""
}

func (s *Session) FlushMessage() string {
	message := s.message
	s.message = ""
	return message
}
