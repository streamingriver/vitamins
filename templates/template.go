package templates

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/google/renameio"
)

func Parse(file string, data interface{}) *Template {
	t := &Template{
		template.Must(template.ParseFiles(file)),
		data,
	}
	return t
}

type Template struct {
	temp *template.Template
	data interface{}
}

func (t *Template) Stdout() {
	t.temp.Execute(os.Stdout, t.data)
}

func (t *Template) Out() string {
	buff := bytes.NewBuffer(nil)
	t.temp.Execute(buff, t.data)
	return buff.String()
}

func (t *Template) Write(file string) error {

	fh, err := renameio.TempFile(filepath.Dir(file), path.Base(file))

	if err != nil {
		return err
	}

	err = t.temp.Execute(fh, t.data)
	if err != nil {
		return err
	}
	return fh.CloseAtomicallyReplace()
}
