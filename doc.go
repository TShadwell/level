
/*
	Package level abstracts a C and Go implimentation of levelDB through use of
	compile-time tags.

	The tag 'purego' can be used to compile with the go implimentation,
	github.com/syndtr/goleveldb , otherwise the C- implimentation,
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
	to generate keys and values respectively, it can be used more directly with the database:

		testAtom.Object().Delete(
			//This can be deleted, given it impliments KeyMarshaler
			things,
		)

	Types that impliment both KeyMarshaler and ValueMarshaler can be placed directly
	in the database.
	
		testAtom.Object().Place(
			val,
		)

	As well as being Written to the database, atoms can be committed, which
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

//These represent the interfaces to which implimentations must conform
//these will be extended and abstracted by their exported versions
type (
	//A Key, for the database
	Key []byte
	//A Value, for the database
	Value  []byte
	closer interface {
		Close()
	}
	options interface {
		SetCreateIfMissing(yes bool)
		SetCache(cache)
		closer
	}
	database interface {
		closer
		Delete(writeOptions, Key) error
		Put(writeOptions, Key, Value) error
		Write(writeOptions, writeBatch) error
		Get(readOptions, Key) (Value, error)
	}
	writeOptions interface {
		closer
		SetSync(sync bool)
	}
	readOptions interface {
		closer
		SetVerifyChecksums(yes bool)
	}
	writeBatch interface {
		closer
		Clear()
		Delete(Key)
		Put(Key, Value)
	}
	cache interface {
		closer
	}
)

//Define the abstract implimentations of the interfaces.
type (
	//Database options
	Options struct {
		options
	}
	//LRU Cache
	Cache struct {
		cache
	}
	//General write options
	WriteOptions struct {
		writeOptions
	}
	//General read options
	ReadOptions struct {
		readOptions
	}
	//A levelDB database
	Database struct {
		database
		Cache Cache
		*Options
		*ReadOptions
		*WriteOptions
	}
	//type Atom represents series of deletions and writes that all fail and
	//do not commit if one fails.
	Atom struct {
		writeBatch
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
		Inner() writeBatch
	}
)

func (v Value) MarshalValue() Value {
	return v
}

func (k Key) MarshalKey() Key {
	return k
}

/*
	=== Options Functions ===
*/

/*
	Gets the underlying implimentation of Options as an interface,
	creating it if it doesn't already exist.
*/
func (o *Options) Inner() options {
	if o == nil {
		o = new(Options)
	}
	if o.options == nil {
		o.options = newOptions()
	}
	return o.options
}

/*
	Function SetCreateIfMissing causes an attempt
	to open a database to also create it if it did not exist.
*/
func (o *Options) SetCreateIfMissing(b bool) *Options {
	o.Inner().SetCreateIfMissing(b)
	return o
}

/*
	Function SetCache sets the cache object for the database
*/
func (o *Options) SetCache(c *Cache) *Options {
	o.Inner().SetCache(c.Inner())
	return o
}

/*
	Function SetCacheSize sets the cache object for the database to a new cache of given size.
*/
func (o *Options) SetCacheSize(size BytesSize) *Options {
	o.SetCache(new(Cache).Size(size))
	return o
}

/*
	=== Cache Functions ===
*/

/*
	Function Inner returns the underlying implimentation of the Cache.

	Unlike other Inner Functions, this may return nil, since LRUCaches must
	be created with given size.

	Therefore, Size(BytesSize) should be called before this.
*/
func (c *Cache) Inner() cache {
	return c.cache
}

/*
	Function Size sets the size of the underlying LRUCache.
*/
func (c *Cache) Size(b BytesSize) *Cache {
	c.cache = newLRUCache(int(b))
	return c
}

/*
	=== Write Options Functions ===
*/

/*
	Returns the underlying writeOptions interface
	of this WriteOptions, creating it if it does not
	exist.
*/
func (w *WriteOptions) Inner() writeOptions {
	if w == nil {
		w = new(WriteOptions)
	}
	if w.writeOptions == nil {
		w.writeOptions = newWriteOptions()
	}
	return w.writeOptions
}

/*
	Function SetSync sets whether these writes will be flushed
	immediately from the buffer cache. This slows down writes
	but has better crash semantics.
*/
func (w *WriteOptions) SetSync(b bool) *WriteOptions {
	w.Inner().SetSync(b)
	return w
}

/*
	=== Read Options Functions ===
*/

/*
	Returns the underlying readOptions interface
	of this ReadOptions, creating it if it does not
	exist.
*/
func (r *ReadOptions) Inner() readOptions {
	if r == nil {
		r = new(ReadOptions)
	}
	if r.readOptions == nil {
		r.readOptions = newReadOptions()
	}
	return r.readOptions
}

func (r *ReadOptions) SetVerifyChecksums(b bool) *ReadOptions {
	r.Inner().SetVerifyChecksums(b)
	return r
}

/*
	=== Database Functions ===
*/

/*
	Function open opens a database for writes at location.
*/
func (d *Database) Open(location string) (err error) {
	if d.database != nil {
		err = Already_Open
		return
	}
	var dt database
	dt, err = openDatabase(location, d.Options.Inner())
	d.database = dt
	return
}

/*
	Function OpenDB opens a database for writes at location, function chaining
	syntax.
*/
func (d *Database) OpenDB(location string) (*Database, error) {
	return d, d.Open(location)
}

func (d *Database) Close() {
	d.database.Close()
	d.Cache.Close()
	d.Options.Close()
	d.ReadOptions.Close()
	d.WriteOptions.Close()
}

/*
	Function SetOptions sets the Options of this Database.
*/
func (d *Database) SetOptions(o *Options) *Database {
	if d.Options != nil {
		panic("Options already set!")
	}
	d.Options = o
	return d
}

/*
	Returns the underlying database of the Database.
	If the database has not been opened, the Not_Open error
	will be returned.
*/
func (d *Database) Inner() (db database, err error) {
	db = d.database
	if db == nil {
		err = Not_Opened
	}
	return
}

/*
	Deletes a single value from the database.
	For batch deletions, use an Atom.
*/
func (d *Database) Delete(k Key) error {
	db, err := d.Inner()
	if err != nil {
		return err
	}
	return db.Delete(d.WriteOptions.Inner(), k)
}

/*
	Puts a single value into the database.
	For batch puts, use an Atom.
*/
func (d *Database) Put(k Key, v Value) error {
	db, err := d.Inner()
	if err != nil {
		return err
	}
	return db.Put(d.WriteOptions.Inner(), k, v)
}

/*
	Gets a single value from the database.
*/
func (d *Database) Get(k Key) (Value, error) {
	db, err := d.Inner()
	if err != nil {
		return nil, err
	}
	return db.Get(d.ReadOptions.Inner(), k)
}

/*
	Write an Atom or InterfaceAtom to the Database.
*/
func (d *Database) Write(an atom) error {
	db, err := d.Inner()
	if err != nil {
		return err
	}
	return db.Write(d.WriteOptions.Inner(), an.Inner())
}

/*
	Write an Atom or InterfaceAtom to the Database,
	closing it afterward.
*/
func (d *Database) Commit(an atom) error {
	defer an.Inner().Close()
	return d.Write(an)
}

func (c *Cache) Close() {
	if c != nil && c.cache != nil {
		c.cache.Close()
	}
}

func (o *Options) Close() {
	if o != nil && o.options != nil {
		o.options.Close()
	}
}

func (r *ReadOptions) Close() {
	if r != nil && r.readOptions != nil {
		r.readOptions.Close()
	}
}

func (w *WriteOptions) Close() {
	if w != nil && w.writeOptions != nil {
		w.writeOptions.Close()
	}
}

/*
	Returns the underlying writeBatch of this Atom,
	creating it if it does not exist.
*/
func (a *Atom) Inner() writeBatch {
	if a.writeBatch == nil {
		a.writeBatch = newWriteBatch()
	}
	return a.writeBatch
}

/*
	Empty the writes and deletes of this Atom.
*/
func (a *Atom) Clear() *Atom {
	a.Inner().Clear()
	return a
}

func (a *Atom) Close() *Atom {
	a.Inner().Close()
	return a
}

/*
	Delete a Value from the database.
*/
func (a *Atom) Delete(k Key) *Atom {
	a.Inner().Delete(k)
	return a
}

/*
	Store a Value at Key.
*/
func (a *Atom) Put(k Key, v Value) *Atom {
	a.Inner().Put(k, v)
	return a
}

/*
	Returns an InterfaceAtom that allows more abstracted puts and deletions
	of values.

	An OOAtom is a reference type.
*/
func (a *Atom) Object() InterfaceAtom {
	return InterfaceAtom{a}
}

func (o InterfaceAtom) Delete(k KeyMarshaler) InterfaceAtom {
	o.Atom.Delete(k.MarshalKey())
	return o
}

func (o InterfaceAtom) Place(kv KeyValueMarshaler) InterfaceAtom {
	o.Atom.Put(kv.MarshalKey(), kv.MarshalValue())
	return o
}

func (o InterfaceAtom) Put(k KeyMarshaler, v ValueMarshaler) InterfaceAtom {
	o.Atom.Put(k.MarshalKey(), v.MarshalValue())
	return o
}
