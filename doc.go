/*
	Package level abstracts a C and Go implementation of levelDB through use of
	compile-time tags.

	The tag 'purego' can be used to compile with the go implementation,
	github.com/syndtr/goleveldb , otherwise the C- implementation,
	github.com/jmhodges/levigo is used.

	A number of different syntaxes can be used with level,
	it is designed to be friendly with new() syntax and function
	chaining.

		//Open a Database
		db, err := new(Database).SetOptions(
			new(Options).SetCreateIfMissing(
				true,
			).SetCacheSize(
				500 * Megabyte,
			),
		).OpenDB(path + "/leveldb")

	Atoms can be used for atomic writes and deletions

		testAtom := new(Atom).Put(
			[]byte("beans"),
			[]byte("can),
		)

	Atoms are also abstracted using the interfaces KeyMarshaler and ValueMarshaler,
	if a type impliments the methods MarshalKey() Key and MarshalValue() Value,
	to generate keys and values respectively, it can be used more directly with the UnderlyingDatabase:

		testAtom.Object().Delete(
			//This can be deleted, given it impliments KeyMarshaler
			things,
		)

	Types that impliment both KeyMarshaler and ValueMarshaler can be placed directly
	in the UnderlyingDatabase.

		testAtom.Object().Place(
			val,
		)

	As well as being Written to the UnderlyingDatabase, atoms can be committed, which
	closes the underlying structure.

		err = db.Commit(testAtom)


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
	Value  []byte
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
