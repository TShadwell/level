package dex

import (
	"github.com/TShadwell/level"
	gl "github.com/TShadwell/level/golevel"
	"testing"
)

const CatType Type = iota

type Cat struct {
	Name string
}

func (Cat) TypeDex() Type{
	return CatType
}

func (c Cat) MarshalDex() []byte{
	return []byte(c.Name)
}

func (c *Cat) UnmarshalDex(b []byte) error{
	c.Name = string(b)
	return nil
}

func TestDex(t *testing.T) {
	db := level.Database{
		Cache: gl.Level.NewCache(500 * level.Megabyte),
		Options: gl.Level.NewOptions().SetCreateIfMissing(
			true,
		),
	}

	if err := gl.Level.OpenDatabase(&db, "leveldb"); err != nil{
		panic(err)
	}

	dx := Dex{
		&db,
	}

	const catName = "Michael"

	if err := dx.Store(Cat{catName}, 0); err != nil{
		panic(err)
	}

	var Michael Cat
	if err := dx.Retrieve(&Michael, 0); err != nil{
		panic(err)
	}


}
