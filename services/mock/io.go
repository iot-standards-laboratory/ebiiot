package mock

import (
	"errors"
	"io"
)

func ReadPayload(reader io.Reader) ([]byte, error) {
	lenBuffer := make([]byte, 2)

	for lenRead := 0; lenRead < 2; {
		n, err := reader.Read(lenBuffer[lenRead:])
		if err != nil {
			return nil, err
		}

		lenRead += n
	}

	length := int16(0)
	length = int16(lenBuffer[1])
	length = length<<8 + int16(lenBuffer[0])

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

	length := make([]byte, 2)
	length[0] = byte(msgLength)
	length[1] = byte(msgLength >> 8)
	writer.Write(length)

	writer.Write(msg.ToBytes())
}
