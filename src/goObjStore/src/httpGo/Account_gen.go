package httpGo

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Account) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "containers":
			var zb0002 uint32
			zb0002, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Containers == nil && zb0002 > 0 {
				z.Containers = make(map[string]Container, zb0002)
			} else if len(z.Containers) > 0 {
				for key, _ := range z.Containers {
					delete(z.Containers, key)
				}
			}
			for zb0002 > 0 {
				zb0002--
				var za0001 string
				var za0002 Container
				za0001, err = dc.ReadString()
				if err != nil {
					return
				}
				err = za0002.DecodeMsg(dc)
				if err != nil {
					return
				}
				z.Containers[za0001] = za0002
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
func (z *Account) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "name"
	err = en.Append(0x82, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "containers"
	err = en.Append(0xaa, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x73)
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.Containers)))
	if err != nil {
		return
	}
	for za0001, za0002 := range z.Containers {
		err = en.WriteString(za0001)
		if err != nil {
			return
		}
		err = za0002.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Account) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "name"
	o = append(o, 0x82, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "containers"
	o = append(o, 0xaa, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x73)
	o = msgp.AppendMapHeader(o, uint32(len(z.Containers)))
	for za0001, za0002 := range z.Containers {
		o = msgp.AppendString(o, za0001)
		o, err = za0002.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Account) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "containers":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.Containers == nil && zb0002 > 0 {
				z.Containers = make(map[string]Container, zb0002)
			} else if len(z.Containers) > 0 {
				for key, _ := range z.Containers {
					delete(z.Containers, key)
				}
			}
			for zb0002 > 0 {
				var za0001 string
				var za0002 Container
				zb0002--
				za0001, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				bts, err = za0002.UnmarshalMsg(bts)
				if err != nil {
					return
				}
				z.Containers[za0001] = za0002
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
func (z *Account) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Name) + 11 + msgp.MapHeaderSize
	if z.Containers != nil {
		for za0001, za0002 := range z.Containers {
			_ = za0002
			s += msgp.StringPrefixSize + len(za0001) + za0002.Msgsize()
		}
	}
	return
}
