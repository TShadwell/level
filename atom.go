package level
func (l *Level) NewAtom() *Atom{
	return &Atom{
		l.NewWriteBatch(),
	}
}

/*
	Empty the writes and deletes of this Atom.
*/
func (a *Atom) Clear() *Atom {
	a.UnderlyingWriteBatch.Clear()
	return a
}

func (a *Atom) Close() *Atom {
	a.UnderlyingWriteBatch.Close()
	return a
}

/*
	Delete a Value from the UnderlyingDatabase.
*/
func (a *Atom) Delete(k Key) *Atom {
	a.UnderlyingWriteBatch.Delete(k)
	return a
}

/*
	Store a Value at Key.
*/
func (a *Atom) Put(k Key, v Value) *Atom {
	a.UnderlyingWriteBatch.Put(k, v)
	return a
}
