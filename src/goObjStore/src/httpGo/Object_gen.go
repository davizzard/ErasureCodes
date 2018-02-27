package httpGo

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Object) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "size":
			z.Size, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "PartsNum":
			z.PartsNum, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "ParityNum":
			z.ParityNum, err = dc.ReadInt()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Object) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "name"
	err = en.Append(0x84, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "size"
	err = en.Append(0xa4, 0x73, 0x69, 0x7a, 0x65)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Size)
	if err != nil {
		return
	}
	// write "PartsNum"
	err = en.Append(0xa8, 0x50, 0x61, 0x72, 0x74, 0x73, 0x4e, 0x75, 0x6d)
	if err != nil {
		return
	}
	err = en.WriteInt(z.PartsNum)
	if err != nil {
		return
	}
	// write "ParityNum"
	err = en.Append(0xa9, 0x50, 0x61, 0x72, 0x69, 0x74, 0x79, 0x4e, 0x75, 0x6d)
	if err != nil {
		return
	}
	err = en.WriteInt(z.ParityNum)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Object) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "name"
	o = append(o, 0x84, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "size"
	o = append(o, 0xa4, 0x73, 0x69, 0x7a, 0x65)
	o = msgp.AppendInt(o, z.Size)
	// string "PartsNum"
	o = append(o, 0xa8, 0x50, 0x61, 0x72, 0x74, 0x73, 0x4e, 0x75, 0x6d)
	o = msgp.AppendInt(o, z.PartsNum)
	// string "ParityNum"
	o = append(o, 0xa9, 0x50, 0x61, 0x72, 0x69, 0x74, 0x79, 0x4e, 0x75, 0x6d)
	o = msgp.AppendInt(o, z.ParityNum)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Object) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "size":
			z.Size, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "PartsNum":
			z.PartsNum, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "ParityNum":
			z.ParityNum, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Object) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Name) + 5 + msgp.IntSize + 9 + msgp.IntSize + 10 + msgp.IntSize
	return
}
