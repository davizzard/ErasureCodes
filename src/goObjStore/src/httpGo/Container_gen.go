package httpGo

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Container) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "objs":
			var zb0002 uint32
			zb0002, err = dc.ReadMapHeader()
			CheckSimpleErr(err, nil, true)
			if z.Objs == nil && zb0002 > 0 {
				z.Objs = make(map[string]Object, zb0002)
			} else if len(z.Objs) > 0 {
				for key, _ := range z.Objs {
					delete(z.Objs, key)
				}
			}
			for zb0002 > 0 {
				zb0002--
				var za0001 string
				var za0002 Object
				za0001, err = dc.ReadString()
				CheckSimpleErr(err, nil, true)
				err = za0002.DecodeMsg(dc)
				CheckSimpleErr(err, nil, true)
				z.Objs[za0001] = za0002
			}
		case "policy":
			z.Policy, err = dc.ReadString()
			CheckSimpleErr(err, nil, true)
		default:
			err = dc.Skip()
			CheckSimpleErr(err, nil, true)
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Container) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "name"
	err = en.Append(0x83, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	CheckSimpleErr(err, nil, true)
	err = en.WriteString(z.Name)
	CheckSimpleErr(err, nil, true)
	// write "objs"
	err = en.Append(0xa4, 0x6f, 0x62, 0x6a, 0x73)
	CheckSimpleErr(err, nil, true)
	err = en.WriteMapHeader(uint32(len(z.Objs)))
	CheckSimpleErr(err, nil, true)
	for za0001, za0002 := range z.Objs {
		err = en.WriteString(za0001)
		CheckSimpleErr(err, nil, true)
		err = za0002.EncodeMsg(en)
		CheckSimpleErr(err, nil, true)
	}
	// write "policy"
	err = en.Append(0xa6, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79)
	CheckSimpleErr(err, nil, true)
	err = en.WriteString(z.Policy)
	CheckSimpleErr(err, nil, true)
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Container) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "name"
	o = append(o, 0x83, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "objs"
	o = append(o, 0xa4, 0x6f, 0x62, 0x6a, 0x73)
	o = msgp.AppendMapHeader(o, uint32(len(z.Objs)))
	for za0001, za0002 := range z.Objs {
		o = msgp.AppendString(o, za0001)
		o, err = za0002.MarshalMsg(o)
		CheckSimpleErr(err, nil, true)
	}
	// string "policy"
	o = append(o, 0xa6, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79)
	o = msgp.AppendString(o, z.Policy)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Container) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "objs":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			CheckSimpleErr(err, nil, true)
			if z.Objs == nil && zb0002 > 0 {
				z.Objs = make(map[string]Object, zb0002)
			} else if len(z.Objs) > 0 {
				for key, _ := range z.Objs {
					delete(z.Objs, key)
				}
			}
			for zb0002 > 0 {
				var za0001 string
				var za0002 Object
				zb0002--
				za0001, bts, err = msgp.ReadStringBytes(bts)
				CheckSimpleErr(err, nil, true)
				bts, err = za0002.UnmarshalMsg(bts)
				CheckSimpleErr(err, nil, true)
				z.Objs[za0001] = za0002
			}
		case "policy":
			z.Policy, bts, err = msgp.ReadStringBytes(bts)
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
func (z *Container) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Name) + 5 + msgp.MapHeaderSize
	if z.Objs != nil {
		for za0001, za0002 := range z.Objs {
			_ = za0002
			s += msgp.StringPrefixSize + len(za0001) + za0002.Msgsize()
		}
	}
	s += 7 + msgp.StringPrefixSize + len(z.Policy)
	return
}
