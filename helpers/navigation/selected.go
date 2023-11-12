package navigation

type Selected uint

func (s *Selected) Next(max int) {
	if int(*s) < max {
		*s++
	}
}

func (s *Selected) Prev() {
	if *s > 0 {
		*s--
	}
}
