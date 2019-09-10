package httpc

type retryErr struct {
	error
}

func (r *retryErr) Retry() bool { return true }
