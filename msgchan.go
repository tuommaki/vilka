package vilka

import (
	"encoding/gob"
	"fmt"
	"os"
)

// MsgChan is an interface to represent required functionality for
// inter-process messaging
type MsgChan interface {
	// Send message via message channel
	Send(msg interface{}) error

	// Receive message via message channel
	// msg must be a pointer to work correctly
	Recv(msg interface{}) error

	// Close message channel
	Close()
}

type msgChan struct {
	dec  *gob.Decoder
	enc  *gob.Encoder
	file *os.File
}

func msgChanFromFd(fd uintptr) (MsgChan, error) {
	file := os.NewFile(fd, "<msgchan>")
	if file == nil {
		return nil, fmt.Errorf("invalid file descriptor: %d", fd)
	}

	return &msgChan{
		dec:  gob.NewDecoder(file),
		enc:  gob.NewEncoder(file),
		file: file,
	}, nil
}

func (mc *msgChan) Send(msg interface{}) error {
	return mc.enc.Encode(msg)
}

func (mc *msgChan) Recv(msg interface{}) error {
	return mc.dec.Decode(msg)
}

func (mc *msgChan) Close() {
	mc.file.Close()
}
