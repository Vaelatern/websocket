package websocket

import (
	"fmt"
	"io"

	"golang.org/x/xerrors"
)

type limitedReader struct {
	c     *Conn
	r     io.Reader
	left  int64
	limit int64
}

func (lr *limitedReader) Read(p []byte) (int, error) {
	if lr.limit == 0 {
		lr.limit = lr.left
	}

	if lr.left <= 0 {
		msg := fmt.Sprintf("read limited at %v bytes", lr.limit)
		lr.c.Close(StatusPolicyViolation, msg)
		return 0, xerrors.Errorf(msg)
	}

	if int64(len(p)) > lr.left {
		p = p[:lr.left]
	}
	n, err := lr.r.Read(p)
	lr.left -= int64(n)
	return n, err
}