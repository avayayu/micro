package api

func (s *Status) Error() string {
	return s.Message
}
