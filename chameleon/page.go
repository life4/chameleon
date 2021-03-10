package chameleon

import (
	"io"
)

type Page interface {
	Render(io.Writer) error
	Inc() error
	Status() int
}
