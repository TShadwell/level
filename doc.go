/*
	Package level abstracts levelDB, providing a midway point through which several implementations
	of levelDB may be used. In previous versions this relied on build tags, but this restricted
	how buildtags could be used when this package was imported, and extensibility, this way an
	importer may specify that heroku builds must use goleveldb for example.

	The subpackage /golevel provides a *level.Level corresponding to github.com/syndtr/goleveldb,
	and the subpackage /levigo provides one corresponding to github.com/jmhodges/levigo

	It is important to note that there is no system that allows compatibility between the abstracted
	types of different implimentations, trying to mix them will usually cause assertion runtime panics.

		var lvl *level.Level

		//Import the package of your choice
		lvl = levigo.Level()

		db := &level.Database{
			Cache: lvl.NewCache(500 *level.Megabyte),
			Options: lvl.NewOptions().SetCreateIfMissing(
				true,
			),
		}

		if err := lvl.OpenDatabase(db, path+"/leveldb/"); err != nil{
			t.Fatal("Error whilst loading DB:", err)
		}

	Atoms can be used for atomic writes and deletions

		testAtom := lvl.NewAtom().Put(
			[]byte("beans"),
			[]byte("can"),
		)

	As well as being Written to the UnderlyingDatabase, atoms can be committed, which
	closes the underlying structure.

		err = db.Commit(testAtom)

	The /legacy package has the same interface as previous versions, which used build tags.

*/
package level

//Welcome to wrapper central

/*
	A number of bytes, for sizing the LRU cache.
*/
type BytesSize uint

const (
	Byte = 1 << (10 * iota)
	Kilobyte
	Megabyte
)

//The interfaces to which implementations must conform,
//these will be extended and abstracted by their exported versions
type (
	//A Key, for the UnderlyingDatabase
	Key []byte
	//A Value, for the UnderlyingDatabase
	Value             []byte
	UnderlyingOptions interface {
		SetCreateIfMissing(yes bool)
		SetCache(UnderlyingCache)
		Close()
	}
	UnderlyingDatabase interface {
		Close()
		Delete(UnderlyingWriteOptions, Key) error
		Put(UnderlyingWriteOptions, Key, Value) error
		Write(UnderlyingWriteOptions, UnderlyingWriteBatch) error
		Get(UnderlyingReadOptions, Key) (Value, error)
	}
	UnderlyingWriteOptions interface {
		Close()
		SetSync(sync bool)
	}
	UnderlyingReadOptions interface {
		Close()
		SetVerifyChecksums(yes bool)
	}
	UnderlyingWriteBatch interface {
		Close()
		Clear()
		Delete(Key)
		Put(Key, Value)
	}
	UnderlyingCache interface {
		Close()
	}
)

//Define the abstract implementations of the interfaces.
type (
	//Database UnderlyingOptions
	Options struct {
		UnderlyingOptions
	}
	//LRU Cache
	Cache struct {
		UnderlyingCache
	}
	//General write UnderlyingOptions
	WriteOptions struct {
		UnderlyingWriteOptions
	}
	//General read UnderlyingOptions
	ReadOptions struct {
		UnderlyingReadOptions
	}
	//A levelDB UnderlyingDatabase
	Database struct {
		UnderlyingDatabase
		Cache *Cache
		*Options
		*ReadOptions
		*WriteOptions
	}
	//type Atom represents series of deletions and writes that all fail and
	//do not commit if one fails.
	Atom struct {
		UnderlyingWriteBatch
	}

	//type InterfaceAtom abstracts Puts and Deletes
	//to an atom to allow more direct use of interfacing.
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

	atom interface {
		Inner() UnderlyingWriteBatch
	}
	UnderlyingLevel interface {
		NewLRUCache(capacity int) UnderlyingCache
		DestroyDatabase(name string, o UnderlyingOptions) error
		RepairDatabase(name string, o UnderlyingOptions) error
		OpenDatabase(name string, o UnderlyingOptions) (UnderlyingDatabase, error)
		NewOptions() UnderlyingOptions
		NewReadOptions() UnderlyingReadOptions
		NewWriteOptions() UnderlyingWriteOptions
		NewWriteBatch() UnderlyingWriteBatch
	}
)

func (v Value) MarshalValue() Value {
	return v
}

func (k Key) MarshalKey() Key {
	return k
}

func (c *Cache) Close() {
	if c != nil && c.UnderlyingCache != nil {
		c.UnderlyingCache.Close()
	}
}

func (o *Options) Close() {
	if o != nil && o.UnderlyingOptions != nil {
		o.UnderlyingOptions.Close()
	}
}

func (r *ReadOptions) Close() {
	if r != nil && r.UnderlyingReadOptions != nil {
		r.UnderlyingReadOptions.Close()
	}
}

func (w *WriteOptions) Close() {
	if w != nil && w.UnderlyingWriteOptions != nil {
		w.UnderlyingWriteOptions.Close()
	}
}
