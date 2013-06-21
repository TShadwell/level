package level

func (l *Level) NewReadOptions() *ReadOptions {
	return &ReadOptions{
		l.UnderlyingLevel.NewReadOptions(),
	}
}

func (r *ReadOptions) SetVerifyChecksums(yes bool) *ReadOptions {
	r.UnderlyingReadOptions.SetVerifyChecksums(yes)
	return r
}
