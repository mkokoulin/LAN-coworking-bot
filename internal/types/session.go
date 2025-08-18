package types

func (s *Session) ResetFlow() {
	s.Flow = ""
	s.Step = ""
}
