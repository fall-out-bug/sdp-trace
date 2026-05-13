package prreview

import (
	"io"
)

func Copy(r io.Reader, w io.Writer) error {
	_, err := io.Copy(w, r)
	return err
}
