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
	CheckSimpleErr(err, nil, true)
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		CheckSimpleErr(err, nil, true)
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, err = dc.ReadString()
			CheckSimpleErr(err, nil, true)
		case "size":
			z.Size, err = dc.ReadInt()
			CheckSimpleErr(err, nil, true)
		case "PartsNum":
			z.PartsNum, err = dc.ReadInt()
			CheckSimpleErr(err, nil, true)
		case "ParityNum":
			z.ParityNum, err = dc.ReadInt()
			CheckSimpleErr(err, nil, true)
		default:
			err = dc.Skip()
			CheckSimpleErr(err, nil, true)
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Object) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "name"
	err = en.Append(0x84, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	CheckSimpleErr(err, nil, true)
	err = en.WriteString(z.Name)
	CheckSimpleErr(err, nil, true)
	// write "size"
	err = en.Append(0xa4, 0x73, 0x69, 0x7a, 0x65)
	CheckSimpleErr(err, nil, true)
	err = en.WriteInt(z.Size)
	CheckSimpleErr(err, nil, true)
	// write "PartsNum"
	err = en.Append(0xa8, 0x50, 0x61, 0x72, 0x74, 0x73, 0x4e, 0x75, 0x6d)
	CheckSimpleErr(err, nil, true)
	err = en.WriteInt(z.PartsNum)
	CheckSimpleErr(err, nil, true)
	// write "ParityNum"
	err = en.Append(0xa9, 0x50, 0x61, 0x72, 0x69, 0x74, 0x79, 0x4e, 0x75, 0x6d)
	CheckSimpleErr(err, nil, true)
	err = en.WriteInt(z.ParityNum)
	CheckSimpleErr(err, nil, true)
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
	CheckSimpleErr(err, nil, true)
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		CheckSimpleErr(err, nil, true)
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			CheckSimpleErr(err, nil, true)
		case "size":
			z.Size, bts, err = msgp.ReadIntBytes(bts)
			CheckSimpleErr(err, nil, true)
		case "PartsNum":
			z.PartsNum, bts, err = msgp.ReadIntBytes(bts)
			CheckSimpleErr(err, nil, true)
		case "ParityNum":
			z.ParityNum, bts, err = msgp.ReadIntBytes(bts)
			CheckSimpleErr(err, nil, true)
		default:
			bts, err = msgp.Skip(bts)
			CheckSimpleErr(err, nil, true)
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
