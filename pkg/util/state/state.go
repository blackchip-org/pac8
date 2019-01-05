package state

import (
	"encoding/gob"
	"io"
)

type Encoder struct {
	Err error
	enc *gob.Encoder
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		enc: gob.NewEncoder(w),
	}
}

func (e *Encoder) Encode(v interface{}) {
	if e.Err != nil {
		return
	}
	e.Err = e.enc.Encode(v)
}

type Decoder struct {
	Err error
	dec *gob.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		dec: gob.NewDecoder(r),
	}
}

func (d *Decoder) Decode(v interface{}) {
	if d.Err != nil {
		return
	}
	d.Err = d.dec.Decode(v)
}
