// +build purego

package level

import (
	"github.com/syndtr/goleveldb/leveldb"
	C "github.com/syndtr/goleveldb/leveldb/cache"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

type db struct {
	*leveldb.DB
	*storage.FileStorage
}

func (d db) Close() {
	e := d.DB.Close()
	//Fix this nasty kludge
	if e != nil {
		panic(e)
	}
	e = d.FileStorage.Close()
	if e != nil {
		panic(e)
	}
}

func (d db) Put(w writeOptions, k Key, v Value) error {
	return d.DB.Put(k, v, w.(wopts).WriteOptions)
}

func (d db) Delete(w writeOptions, k Key) error {
	return d.DB.Delete(k, w.(wopts).WriteOptions)
}

func (d db) Get(ro readOptions, k Key) (v Value, e error) {
	v, e = d.DB.Get(k, ro.(ropts).ReadOptions)
	if e == errors.ErrNotFound {
		//so it works like levigo
		e = nil
	}
	return
}

func (d db) Write(w writeOptions, a writeBatch) error {
	return d.DB.Write(a.(wb).Batch, w.(wopts).WriteOptions)
}

type wb struct {
	*leveldb.Batch
}

func (w wb) batch() *leveldb.Batch {
	if w.Batch == nil {
		w.Batch = new(leveldb.Batch)
	}
	return w.Batch
}

func (w wb) Delete(k Key) {
	w.batch().Delete(k)
}

func (w wb) Put(k Key, v Value) {
	w.batch().Put(k, v)
}

func (w wb) Clear() {
	w.batch().Reset()
}

func (w wb) Close() {
	w.Batch = nil
}

type ropts struct {
	*opt.ReadOptions
}

func (r ropts) Close() {
	r.ReadOptions = nil
}

func (r ropts) readOptions() *opt.ReadOptions {
	if r.ReadOptions == nil {
		r.ReadOptions = new(opt.ReadOptions)
	}
	return r.ReadOptions
}

func (r ropts) SetVerifyChecksums(b bool) {
	if b {
		r.readOptions().Flag |= opt.RFVerifyChecksums
	} else {
		r.readOptions().Flag &^= opt.RFVerifyChecksums
	}
}

type opts struct {
	*opt.Options
}

type wopts struct {
	*opt.WriteOptions
}

func (w wopts) Close() {
	w.WriteOptions = nil
}

func (w wopts) writeOptions() *opt.WriteOptions {
	if w.WriteOptions == nil {
		w.WriteOptions = new(opt.WriteOptions)
	}
	return w.WriteOptions
}

func (w wopts) SetSync(b bool) {
	if b {
		w.writeOptions().Flag |= opt.WFSync
	} else {
		w.writeOptions().Flag &^= opt.WFSync
	}
}

func (o opts) Close() {
	o.Options = nil
}

func (o opts) options() *opt.Options {
	if o.Options == nil {
		o.Options = new(opt.Options)
	}

	return o.Options
}

func (o opts) SetCreateIfMissing(b bool) {
	if b {
		o.options().Flag |= opt.OFCreateIfMissing
	} else {
		o.options().Flag &^= opt.OFCreateIfMissing
	}
}

func (o opts) SetCache(c cache) {
	o.Options.BlockCache = c.(che).Cache
}

type che struct {
	C.Cache
}

func (c che) Close() {
	c.Cache.Purge(func() {
		c.Cache = nil
	})
}

func openDatabase(name string, o options) (dtb database, err error) {
	stor, err := storage.OpenFile(name)
	if err != nil {
		return
	}
	var dtbe *leveldb.DB
	dtbe, err = leveldb.Open(stor, o.(opts).Options)
	dtb = db{dtbe, stor}
	return
}

func newLRUCache(capacity int) cache {
	return che{
		Cache: C.NewLRUCache(capacity),
	}
}

//BUG: Destroy database not written
func destroyDatabase(name string, o options) error {
	return nil
}

//BUG: Repair database not in go-leveldb
func repairDatabase(name string, o options) error {
	return nil
}

func newOptions() options {
	return opts{
		new(opt.Options),
	}
}
func newReadOptions() readOptions {
	return ropts{
		ReadOptions: new(opt.ReadOptions),
	}
}
func newWriteOptions() writeOptions {
	return wopts{
		new(opt.WriteOptions),
	}
}
func newWriteBatch() writeBatch {
	return wb{
		new(leveldb.Batch),
	}
}
