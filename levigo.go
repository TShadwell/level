// +build !purego

package level

import (
	"github.com/jmhodges/levigo"
)

func newLRUCache(capacity int) cache {
	return levigo.NewLRUCache(capacity)
}

func destroyDatabase(name string, o options) error {
	return levigo.DestroyDatabase(name, o.(opts).Options)
}
func repairDatabase(name string, o options) error {
	return levigo.RepairDatabase(name, o.(opts).Options)
}
func openDatabase(name string, o options) (database, error) {
	dtb, e := levigo.Open(name, o.(opts).Options)
	return db{dtb}, e
}
func newOptions() options {
	return opts{levigo.NewOptions()}
}
func newReadOptions() readOptions {
	return levigo.NewReadOptions()
}
func newWriteOptions() writeOptions {
	return levigo.NewWriteOptions()
}
func newWriteBatch() writeBatch {
	return wtb{levigo.NewWriteBatch()}
}

type db struct {
	*levigo.DB
}

func (d db) Delete(w writeOptions, k Key) error {
	return d.DB.Delete(w.(*levigo.WriteOptions), k)
}
func (d db) Put(w writeOptions, k Key, v Value) error {
	return d.DB.Put(w.(*levigo.WriteOptions), k, v)
}

func (d db) Write(w writeOptions, wb writeBatch) error {
	return d.DB.Write(w.(*levigo.WriteOptions), wb.(wtb).WriteBatch)
}

func (d db) Get(r readOptions, k Key) (Value, error) {
	return d.DB.Get(r.(*levigo.ReadOptions), k)
}

type wtb struct {
	*levigo.WriteBatch
}

func (w wtb) Delete(k Key) {
	w.WriteBatch.Delete(k)
}

func (w wtb) Put(k Key, v Value) {
	w.WriteBatch.Put(k, v)
}

type opts struct {
	*levigo.Options
}

func (o *opts) U() *levigo.Options {
	return o.Options
}

func (o opts) SetCache(c cache) {
	o.U().SetCache(c.(*levigo.Cache))
}
