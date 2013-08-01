package golevel

import (
	"github.com/TShadwell/level"
	"github.com/syndtr/goleveldb/leveldb"
	C "github.com/syndtr/goleveldb/leveldb/cache"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

var Level *level.Level

func init() {
	Level = level.New(ulevel{})
}

type ulevel struct{}

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

func (d db) Put(w level.UnderlyingWriteOptions, k level.Key, v level.Value) error {
	return d.DB.Put(k, v, w.(wopts).WriteOptions)
}

func (d db) Delete(w level.UnderlyingWriteOptions, k level.Key) error {
	return d.DB.Delete(k, w.(wopts).WriteOptions)
}

func (d db) Get(ro level.UnderlyingReadOptions, k level.Key) (v level.Value, e error) {
	v, e = d.DB.Get(k, ro.(ropts).ReadOptions)
	if e == errors.ErrNotFound {
		//so it works like levigo
		e = nil
	}
	return
}

func (d db) Write(w level.UnderlyingWriteOptions, a level.UnderlyingWriteBatch) error {
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

func (w wb) Delete(k level.Key) {
	w.batch().Delete(k)
}

func (w wb) Put(k level.Key, v level.Value) {
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

func (o opts) SetCache(c level.UnderlyingCache) {
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

func (ulevel) OpenDatabase(name string, o level.UnderlyingOptions) (dtb level.UnderlyingDatabase, err error) {
	stor, err := storage.OpenFile(name)
	if err != nil {
		return
	}
	var dtbe *leveldb.DB
	dtbe, err = leveldb.Open(stor, o.(opts).Options)
	dtb = db{dtbe, stor}
	return
}

func (ulevel) NewLRUCache(capacity int) level.UnderlyingCache {
	return che{
		Cache: C.NewLRUCache(capacity),
	}
}

//BUG: Destroy database not written
func (ulevel) DestroyDatabase(name string, o level.UnderlyingOptions) error {
	return nil
}

//BUG: Repair database not in go-leveldb
func (ulevel) RepairDatabase(name string, o level.UnderlyingOptions) error {
	return nil
}

func (ulevel) NewOptions() level.UnderlyingOptions {
	return opts{
		new(opt.Options),
	}
}
func (ulevel) NewReadOptions() level.UnderlyingReadOptions {
	return ropts{
		ReadOptions: new(opt.ReadOptions),
	}
}
func (ulevel) NewWriteOptions() level.UnderlyingWriteOptions {
	return wopts{
		new(opt.WriteOptions),
	}
}
func (ulevel) NewWriteBatch() level.UnderlyingWriteBatch {
	return wb{
		new(leveldb.Batch),
	}
}
