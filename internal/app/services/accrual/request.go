package accrual

type request struct {
	order    string
	user     string
	attempts uint
}

func (r *request) isAttemptsExceeded() bool {
	return r.attempts >= maxAttempts-1
}

func (s *accrualService) retry(req *request) {
	if !req.isAttemptsExceeded() {
		s.requestCh <- &request{
			order:    req.order,
			attempts: req.attempts + 1,
		}
	}
}
