package msg

import "bytes"

type Packer struct {
	separator byte
}

func NewPacker(separator byte) Packer {
	return Packer{
		separator: separator,
	}
}

func (p Packer) Unpack(data []byte) (name string, payload []byte, err error) {
	if id := bytes.IndexByte(data, p.separator); id == -1 {
		name = string(data)
	} else {
		name = string(data[:id])
		payload = data[id+1:]
	}
	return
}

func (p Packer) Pack(name string, payload []byte) (data []byte, err error) {
	data = append(append([]byte(name), p.separator), payload...)
	return
}
