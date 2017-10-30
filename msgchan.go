package vilka

import (
	"fmt"
	"io"
	"os"
)

// MsgChan is an interface to represent required functionality for
// inter-process messaging
type MsgChan interface {
}

type msgChan struct {
	in  io.Reader
	out io.Writer
}

func msgChanFromFd(fd uintptr) (MsgChan, error) {
	file := os.NewFile(fd, "<msgchan>")
	if file != nil {
		return &msgChan{
			in:  file,
			out: file,
		}, nil
	}

	return nil, fmt.Errorf("invalid file descriptor: %d", fd)
}
