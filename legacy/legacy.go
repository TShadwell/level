/*
	Package legacy provides a replacement package for packages that used buildtag dependant versions of level.

	UNFINISHED
*/
package legacy

import (
	"github.com/TShadwell/level"
)

var lvl *level.Level

type (
	BytesSize     level.BytesSize
	Key           level.Key
	Value         level.Value
	Options       level.Options
	Cache         level.Cache
	WriteOptions  level.WriteOptions
	ReadOptions   level.ReadOptions
	Database      level.Database
	Atom          level.Atom
	InterfaceAtom struct {
		*Atom
	}

	KeyMarshaler interface {
		MarshalKey() Key
	}

	ValueMarshaler interface {
		MarshalValue() Value
	}

	KeyValueMarshaler interface {
		KeyMarshaler
		ValueMarshaler
	}

	options      level.UnderlyingOptions
	cache        level.UnderlyingCache
	writeOptions level.UnderlyingWriteOptions
)

func (v Value) MarshalValue() Value {
	return v
}
func (k Key) MarshalKey() Key {
	return k
}

//Options functions

func (o *Options) Inner() options {
	if o == nil {
		o = (*Options)(lvl.NewOptions())
	}
	return o.UnderlyingOptions
}

func (o *Options) down() *level.Options {
	if o == nil {
		o = (*Options)(lvl.NewOptions())
	}
	return (*level.Options)(o)
}

func (o *Options) SetCreateIfMissing(b bool) *Options {
	o.down().SetCreateIfMissing(b)
	return o
}

func (o *Options) SetCache(c *Cache) *Options {
	o.down().SetCache((*level.Cache)(c))
	return o
}

func (o *Options) SetCacheSize(size BytesSize) *Options {
	o.SetCache(new(Cache).Size(size))
	return o
}

//Cache functions

func (c *Cache) down() *level.Cache {
	return (*level.Cache)(c)
}

func (c *Cache) Inner() cache {
	return c.down().UnderlyingCache
}

func (c *Cache) Size(b BytesSize) *Cache {
	c = (*Cache)(lvl.NewCache((level.BytesSize)(b)))
	return c
}

//WriteOptions functions
func (w *WriteOptions) down() *level.WriteOptions {
	if w == nil {
		w = (*WriteOptions)(lvl.NewWriteOptions())
	}
	return (*level.WriteOptions)(w)
}

func (w *WriteOptions) Inner() writeOptions {
	return w.down().UnderlyingWriteOptions
}

func (w *WriteOptions) SetSync(b bool) *WriteOptions {
	w.down().SetSync(b)
	return w
}
