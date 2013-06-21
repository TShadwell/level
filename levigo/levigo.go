// +build !purego

package level

import (
	"github.com/TShadwell/level"
	"github.com/jmhodges/levigo"
)

var lv *level.Level

type ulevel struct{}

func Level() *level.Level {
	if lv == nil {
		lv = level.New(new(ulevel))
	}
	return lv
}

func (ulevel) NewLRUCache(capacity int) level.UnderlyingCache {
	return levigo.NewLRUCache(capacity)
}

func (ulevel) DestroyDatabase(name string, o level.UnderlyingOptions) error {
	return levigo.DestroyDatabase(name, o.(opts).Options)
}
func (ulevel) RepairDatabase(name string, o level.UnderlyingOptions) error {
	return levigo.RepairDatabase(name, o.(opts).Options)
}
func (ulevel) OpenDatabase(name string, o level.UnderlyingOptions) (level.UnderlyingDatabase, error) {
	dtb, e := levigo.Open(name, o.(opts).Options)
	return db{dtb}, e
}
func (ulevel) NewOptions() level.UnderlyingOptions {
	return opts{levigo.NewOptions()}
}
func (ulevel) NewReadOptions() level.UnderlyingReadOptions {
	return levigo.NewReadOptions()
}
func (ulevel) NewWriteOptions() level.UnderlyingWriteOptions {
	return levigo.NewWriteOptions()
}
func (ulevel) NewWriteBatch() level.UnderlyingWriteBatch {
	return wtb{levigo.NewWriteBatch()}
}

type db struct {
	*levigo.DB
}

func (d db) Delete(w level.UnderlyingWriteOptions, k level.Key) error {
	return d.DB.Delete(w.(*levigo.WriteOptions), k)
}
func (d db) Put(w level.UnderlyingWriteOptions, k level.Key, v level.Value) error {
	return d.DB.Put(w.(*levigo.WriteOptions), k, v)
}

func (d db) Write(w level.UnderlyingWriteOptions, wb level.UnderlyingWriteBatch) error {
	return d.DB.Write(w.(*levigo.WriteOptions), wb.(wtb).WriteBatch)
}

func (d db) Get(r level.UnderlyingReadOptions, k level.Key) (level.Value, error) {
	return d.DB.Get(r.(*levigo.ReadOptions), k)
}

type wtb struct {
	*levigo.WriteBatch
}

func (w wtb) Delete(k level.Key) {
	w.WriteBatch.Delete(k)
}

func (w wtb) Put(k level.Key, v level.Value) {
	w.WriteBatch.Put(k, v)
}

type opts struct {
	*levigo.Options
}

func (o *opts) U() *levigo.Options {
	return o.Options
}

func (o opts) SetCache(c level.UnderlyingCache) {
	o.U().SetCache(c.(*levigo.Cache))
}
