package mock

import (
	"errors"
	"io"
)

func ReadPayload(reader io.Reader) ([]byte, error) {
	lenBuffer := make([]byte, 4)

	for lenRead := 0; lenRead < 4; {
		n, err := reader.Read(lenBuffer[lenRead:])
		if err != nil {
			return nil, err
		}

		lenRead += n
	}

	length := int32(0)
	length = int32(lenBuffer[3])
	length = length<<8 + int32(lenBuffer[2])
	length = length<<8 + int32(lenBuffer[1])
	length = length<<8 + int32(lenBuffer[0])

	payload := make([]byte, length)
	payloadRead := 0

	for payloadRead < int(length) {
		n, err := reader.Read(payload[payloadRead:])
		if err != nil {
			return nil, err
		}

		payloadRead += n
	}

	if payloadRead != int(length) {
		return nil, errors.New("different size is received")
	}

	return payload, nil
}

func WritePayload(writer io.Writer, msg MessageI) {
	msgLength := msg.Len()

	// lengthBytes := &bytes.Buffer{}
	// lengthBytes.WriteByte(byte(msgLength))
	// lengthBytes.WriteByte(byte(msgLength >> 8))
	// writer.Write(lengthBytes.Bytes())

	length := make([]byte, 4)
	length[0] = byte(msgLength)
	length[1] = byte(msgLength >> 8)
	length[2] = byte(msgLength >> 16)
	length[3] = byte(msgLength >> 24)
	writer.Write(length)

	writer.Write(msg.ToBytes())
}
