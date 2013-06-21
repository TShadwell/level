package level

func (l *Level) NewOptions() *Options {
	return &Options{
		l.UnderlyingLevel.NewOptions(),
	}
}

/*
	Function SetCreateIfMissing causes an attempt
	to open a UnderlyingDatabase to also create it if it did not exist.
*/
func (o *Options) SetCreateIfMissing(yes bool) *Options {
	o.UnderlyingOptions.SetCreateIfMissing(yes)
	return o
}

/*
	Function SetCache sets the cache object for the UnderlyingDatabase
*/
func (o *Options) SetCache(c *Cache) *Options {
	o.UnderlyingOptions.SetCache(c.UnderlyingCache)
	return o
}
