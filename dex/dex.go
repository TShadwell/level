package dex

import (
	"bytes"
	"encoding/binary"
	"github.com/TShadwell/level"
)

type (
	Type uint64
	Index uint64

	Marshaler interface{
		MarshalDex() []byte
	}

	Typer interface {
		TypeDex() Type
	}

	MarshalTyper interface {
		Marshaler
		Typer
	}

	Unmarshaler interface {
		UnmarshalDex([]byte) error
	}

	UnmarshalTyper interface{
		Unmarshaler
		Typer
	}
)

func itmKey(tp Type, i Index) (level.Key, error){
	var b bytes.Buffer

	if err := binary.Write(&b, binary.LittleEndian, tp); err != nil{
		return nil, err
	}

	if err := binary.Write(&b, binary.LittleEndian, i); err != nil{
		return nil, err
	}
	return b.Bytes(), nil
}

func (t Type) Key(i Index) (level.Key, error){
	return itmKey(t, i)
}

func (i Index) Key(t Type) (level.Key, error){
	return itmKey(t, i)
}

type Dex struct {
	*level.Database
}

func (d Dex) StoreWithType(s Marshaler, t Type, i Index) (err error){
	var k level.Key
	if k, err = i.Key(t); err != nil{
		return
	}
	return d.Put(k, s.MarshalDex())
}

func (d Dex) Store(s MarshalTyper, i Index) error{
	return d.StoreWithType(s, s.TypeDex(), i)
}

func (d Dex) RetrieveWithType(r Unmarshaler, t Type, i Index) (err error) {
	var k level.Key
	if k, err = i.Key(t); err != nil{
		return
	}

	var v level.Value
	if v, err = d.Get(k); err != nil{
		return
	}

	return r.UnmarshalDex(v)
}

func (d Dex) Retrieve(r UnmarshalTyper, i Index) error {
	return d.RetrieveWithType(r, r.TypeDex(), i)
}
