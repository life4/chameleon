package chameleon

import (
	"io"
)

type PageAsset struct {
	Path Path
}

func (page PageAsset) Render(w io.Writer) error {
	f, err := page.Path.Open()
	if err != nil {
		return err
	}
	_, err = io.Copy(w, f)
	err2 := f.Close()
	if err2 != nil {
		return err2
	}
	return err
}

func (page PageAsset) Inc() error {
	return nil
}
