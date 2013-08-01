package tests

import (
	"bitbucket.org/kardianos/osext"
	"bytes"
	"github.com/TShadwell/go-useful/errors"
	"github.com/TShadwell/level"
	glvl "github.com/TShadwell/level/golevel"
	lvigo "github.com/TShadwell/level/levigo"
	"testing"
)

var (
	keyone   = []byte("Alpha")
	keytwo   = []byte("Beta")
	valueone = []byte("x")
	valuetwo = []byte("y")
)

func TestDatabase(t *testing.T) {
	for _, v := range []*level.Level{glvl.Level, lvigo.Level} {
		Tdb(t, v)
	}
}

func Tdb(t *testing.T, lvl *level.Level) {
	path, err := osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}
	db := &level.Database{
		Cache: lvl.NewCache(500 * level.Megabyte),
		Options: lvl.NewOptions().SetCreateIfMissing(
			true,
		),
	}
	err = lvl.OpenDatabase(db, path+"/leveldb/")
	if err != nil {
		t.Fatal("Error whilst loading DB: ", errors.Extend(err))
	}

	writeAtom := lvl.NewAtom().Put(
		keyone,
		valueone,
	).Put(
		keytwo,
		valuetwo,
	)

	err = db.Commit(
		writeAtom,
	)

	if err != nil {
		t.Fatal("Error performing atomic DB write: ", errors.Extend(err))
	}

	t.Log("Atom written: ", writeAtom)

	v, err := db.Get(
		keyone,
	)
	if err != nil {
		t.Fatal("Error retrieving key one: ", errors.Extend(err))
	}

	t.Log("Retrieved value: ", string(v))

	if !bytes.Equal(v, valueone) {
		t.Fatal("Values stored and retrived are not the same!")
	}

	v, err = db.Get(
		keytwo,
	)

	if err != nil {
		t.Fatal("Error retrieving key two: ", errors.Extend(err))
	}

	if !bytes.Equal(v, valuetwo) {
		t.Fatal("Values stored and retrived are not the same!")
	}

	//Delete the values from the DB.

	err = db.Commit(
		lvl.NewAtom().Delete(
			keyone,
		).Delete(
			keytwo,
		),
	)

	if err != nil {
		t.Fatal("Could not delete added keys: ", err)
	}

}
