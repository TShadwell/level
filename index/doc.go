/*
	Package index abstracts level's keys to allow easier retrieval
	and storage of separate categories of data
*/
package index

import (
	"bytes"
	"encoding/binary"
	"github.com/TShadwell/NHTGD2013/database/level"
)

type Index uint

var byteorder = binary.LittleEndian

func Key(in Index, km level.KeyMarshaler) (k []byte, err error) {
	var buf bytes.Buffer
	err = binary.Write(&buf, byteorder, in)
	if err != nil {
		return
	}

	err = binary.Write(&buf, byteorder, km)
	if err != nil {
		return
	}

	k = buf.Bytes()
	return
}

func (in Index) Key(km level.KeyMarshaler) (k []byte, err error) {
	return Key(in, km)
}
