package level

func (l *Level) NewCache(capacity BytesSize) *Cache{
	return &Cache{
		l.NewLRUCache(int(capacity)),
	}
}
