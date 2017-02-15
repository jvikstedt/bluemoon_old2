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

	return name, payload, nil
}

func (p Packer) Pack(name string, payload []byte) (data []byte, err error) {
	name = name + ";"

	return append([]byte(name), payload...), nil
}
