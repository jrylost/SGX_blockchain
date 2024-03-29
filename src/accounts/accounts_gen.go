package accounts

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ExternalAccount) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "version":
			z.Version, err = dc.ReadInt64()
			if err != nil {
				err = msgp.WrapError(err, "Version")
				return
			}
		case "id":
			z.Id, err = dc.ReadBytes(z.Id)
			if err != nil {
				err = msgp.WrapError(err, "Id")
				return
			}
		case "count":
			z.Count, err = dc.ReadInt64()
			if err != nil {
				err = msgp.WrapError(err, "Count")
				return
			}
		case "nonce":
			z.Nonce, err = dc.ReadInt64()
			if err != nil {
				err = msgp.WrapError(err, "Nonce")
				return
			}
		case "balance":
			z.Balance, err = dc.ReadInt64()
			if err != nil {
				err = msgp.WrapError(err, "Balance")
				return
			}
		case "file":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "File")
				return
			}
			if cap(z.File) >= int(zb0002) {
				z.File = (z.File)[:zb0002]
			} else {
				z.File = make([]string, zb0002)
			}
			for za0001 := range z.File {
				z.File[za0001], err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "File", za0001)
					return
				}
			}
		case "contract":
			var zb0003 uint32
			zb0003, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "Contract")
				return
			}
			if cap(z.Contract) >= int(zb0003) {
				z.Contract = (z.Contract)[:zb0003]
			} else {
				z.Contract = make([]string, zb0003)
			}
			for za0002 := range z.Contract {
				z.Contract[za0002], err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "Contract", za0002)
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *ExternalAccount) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 7
	// write "version"
	err = en.Append(0x87, 0xa7, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.Version)
	if err != nil {
		err = msgp.WrapError(err, "Version")
		return
	}
	// write "id"
	err = en.Append(0xa2, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Id)
	if err != nil {
		err = msgp.WrapError(err, "Id")
		return
	}
	// write "count"
	err = en.Append(0xa5, 0x63, 0x6f, 0x75, 0x6e, 0x74)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.Count)
	if err != nil {
		err = msgp.WrapError(err, "Count")
		return
	}
	// write "nonce"
	err = en.Append(0xa5, 0x6e, 0x6f, 0x6e, 0x63, 0x65)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.Nonce)
	if err != nil {
		err = msgp.WrapError(err, "Nonce")
		return
	}
	// write "balance"
	err = en.Append(0xa7, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.Balance)
	if err != nil {
		err = msgp.WrapError(err, "Balance")
		return
	}
	// write "file"
	err = en.Append(0xa4, 0x66, 0x69, 0x6c, 0x65)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.File)))
	if err != nil {
		err = msgp.WrapError(err, "File")
		return
	}
	for za0001 := range z.File {
		err = en.WriteString(z.File[za0001])
		if err != nil {
			err = msgp.WrapError(err, "File", za0001)
			return
		}
	}
	// write "contract"
	err = en.Append(0xa8, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Contract)))
	if err != nil {
		err = msgp.WrapError(err, "Contract")
		return
	}
	for za0002 := range z.Contract {
		err = en.WriteString(z.Contract[za0002])
		if err != nil {
			err = msgp.WrapError(err, "Contract", za0002)
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ExternalAccount) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 7
	// string "version"
	o = append(o, 0x87, 0xa7, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendInt64(o, z.Version)
	// string "id"
	o = append(o, 0xa2, 0x69, 0x64)
	o = msgp.AppendBytes(o, z.Id)
	// string "count"
	o = append(o, 0xa5, 0x63, 0x6f, 0x75, 0x6e, 0x74)
	o = msgp.AppendInt64(o, z.Count)
	// string "nonce"
	o = append(o, 0xa5, 0x6e, 0x6f, 0x6e, 0x63, 0x65)
	o = msgp.AppendInt64(o, z.Nonce)
	// string "balance"
	o = append(o, 0xa7, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65)
	o = msgp.AppendInt64(o, z.Balance)
	// string "file"
	o = append(o, 0xa4, 0x66, 0x69, 0x6c, 0x65)
	o = msgp.AppendArrayHeader(o, uint32(len(z.File)))
	for za0001 := range z.File {
		o = msgp.AppendString(o, z.File[za0001])
	}
	// string "contract"
	o = append(o, 0xa8, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Contract)))
	for za0002 := range z.Contract {
		o = msgp.AppendString(o, z.Contract[za0002])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ExternalAccount) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "version":
			z.Version, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Version")
				return
			}
		case "id":
			z.Id, bts, err = msgp.ReadBytesBytes(bts, z.Id)
			if err != nil {
				err = msgp.WrapError(err, "Id")
				return
			}
		case "count":
			z.Count, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Count")
				return
			}
		case "nonce":
			z.Nonce, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Nonce")
				return
			}
		case "balance":
			z.Balance, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Balance")
				return
			}
		case "file":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "File")
				return
			}
			if cap(z.File) >= int(zb0002) {
				z.File = (z.File)[:zb0002]
			} else {
				z.File = make([]string, zb0002)
			}
			for za0001 := range z.File {
				z.File[za0001], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "File", za0001)
					return
				}
			}
		case "contract":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Contract")
				return
			}
			if cap(z.Contract) >= int(zb0003) {
				z.Contract = (z.Contract)[:zb0003]
			} else {
				z.Contract = make([]string, zb0003)
			}
			for za0002 := range z.Contract {
				z.Contract[za0002], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Contract", za0002)
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *ExternalAccount) Msgsize() (s int) {
	s = 1 + 8 + msgp.Int64Size + 3 + msgp.BytesPrefixSize + len(z.Id) + 6 + msgp.Int64Size + 6 + msgp.Int64Size + 8 + msgp.Int64Size + 5 + msgp.ArrayHeaderSize
	for za0001 := range z.File {
		s += msgp.StringPrefixSize + len(z.File[za0001])
	}
	s += 9 + msgp.ArrayHeaderSize
	for za0002 := range z.Contract {
		s += msgp.StringPrefixSize + len(z.Contract[za0002])
	}
	return
}
