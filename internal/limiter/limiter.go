package limiter

type Limiter struct {
	count chan struct{}
}

func NewLimiter(count int) *Limiter {
	return &Limiter{
		count: make(chan struct{}, count),
	}
}

func (l *Limiter) Acquire() {
	l.count <- struct{}{}
}

func (l *Limiter) Release() {
	<-l.count
}
