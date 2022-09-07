package mock

import (
	"bytes"
	"encoding/binary"
	"time"
)

type MessageI interface {
	ID() int
	Payload() []byte
	Latency() int64
	ToBytes() []byte
	Len() uint32
}

type mockMessage struct {
	id         int
	hour       int
	minute     int
	second     int
	nanosecond int
	payload    []byte
}

func (mm *mockMessage) ID() int {
	return mm.id
}

func (mm *mockMessage) Payload() []byte {
	return mm.payload
}

func (mm *mockMessage) Latency() int64 {
	t := time.Now()

	return t.Sub(time.Date(t.Year(), t.Month(), t.Day(), mm.hour, mm.minute, mm.second, mm.nanosecond, t.Location())).Milliseconds()
}

func (mm *mockMessage) Len() uint32 {
	return uint32(len(mm.payload) + 40)
}

func (mm *mockMessage) ToBytes() []byte {
	buf := &bytes.Buffer{}
	writeInt(buf, mm.id)
	writeInt(buf, mm.hour)
	writeInt(buf, mm.minute)
	writeInt(buf, mm.second)
	writeInt(buf, mm.nanosecond)
	buf.Write(mm.payload)

	return buf.Bytes()
}

func writeInt(writer *bytes.Buffer, value int) {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(value))
	writer.Write(bs)
}

func ParseMsg(p []byte) MessageI {
	msg := &mockMessage{}
	msg.id = int(binary.LittleEndian.Uint64(p[:8]))
	msg.hour = int(binary.LittleEndian.Uint64(p[8:16]))
	msg.minute = int(binary.LittleEndian.Uint64(p[16:24]))
	msg.second = int(binary.LittleEndian.Uint64(p[24:32]))
	msg.nanosecond = int(binary.LittleEndian.Uint64(p[32:40]))
	msg.payload = p[40:]

	return msg
}

func NewMessage(id, size int) MessageI {
	size -= 40
	payloadBytes := make([]byte, size)

	for i := 0; i < size; i++ {
		payloadBytes[i] = 'a'
	}

	t := time.Now()
	return &mockMessage{
		id:         id,
		hour:       t.Hour(),
		minute:     t.Minute(),
		second:     t.Second(),
		nanosecond: t.Nanosecond(),
		payload:    payloadBytes,
	}
}
