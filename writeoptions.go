package level

func (l *Level) NewWriteOptions() *WriteOptions{
	return &WriteOptions{
		l.UnderlyingLevel.NewWriteOptions(),
	}
}

/*
	Function SetSync sets whether these writes will be flushed
	immediately from the buffer cache. This slows down writes
	but has better crash semantics.
*/
func (w *WriteOptions) SetSync(sync bool) *WriteOptions{
	w.UnderlyingWriteOptions.SetSync(sync)
	return w
}
