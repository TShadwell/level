package level

func (l *Level) OpenDatabase(d *Database, location string) (err error) {
	if d.Options == nil {
		d.Options = l.NewOptions()
	}

	if d.ReadOptions == nil {
		d.ReadOptions = l.NewReadOptions()
	}

	if d.WriteOptions == nil {
		d.WriteOptions = l.NewWriteOptions()
	}
	d.UnderlyingDatabase, err = l.UnderlyingLevel.OpenDatabase(location, d.Options.UnderlyingOptions)
	if err != nil {
		return
	}
	return
}

func (d *Database) Close() {
	d.UnderlyingDatabase.Close()
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
	Returns the UnderlyingDatabase of the Database.
	If the UnderlyingDatabase has not been opened, the Not_Open error
	will be returned.
*/
func (d *Database) Inner() UnderlyingDatabase {
	return d.UnderlyingDatabase
}

/*
	Deletes a single value from the UnderlyingDatabase.
	For batch deletions, use an Atom.
*/
func (d *Database) Delete(k Key) error {
	return d.UnderlyingDatabase.Delete(d.WriteOptions.UnderlyingWriteOptions, k)
}

/*
	Puts a single value into the UnderlyingDatabase.
	For batch puts, use an Atom.
*/
func (d *Database) Put(k Key, v Value) error {
	return d.UnderlyingDatabase.Put(d.WriteOptions.UnderlyingWriteOptions, k, v)
}

/*
	Gets a single value from the UnderlyingDatabase.
*/
func (d *Database) Get(k Key) (Value, error) {
	return d.UnderlyingDatabase.Get(d.ReadOptions.UnderlyingReadOptions, k)
}

/*
	Write an Atom to the Database.
*/
func (d *Database) Write(an *Atom) error {
	return d.UnderlyingDatabase.Write(d.WriteOptions.UnderlyingWriteOptions, an.UnderlyingWriteBatch)
}

/*
	Write an Atom or InterfaceAtom to the Database,
	closing it afterward.
*/
func (d *Database) Commit(an *Atom) error {
	defer an.UnderlyingWriteBatch.Close()
	return d.Write(an)
}
