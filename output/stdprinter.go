package output

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/output/functions"
)

//go:embed templates/*.gtpl
var embeddedTemplates embed.FS

// StdPrinter implements `common.TxEntryProcessor` _interface_ and outputs
// the `common.TransactionEntry` instances onto the console for read consumption.
type StdPrinter struct {
	headerTemplate string
	lineTemplate   string
	fullTemplate   string
	header         *template.Template
	output         *template.Template
	w              io.Writer
	entries        []common.TransactionEntry
}

// NewStdPrinter creates a new printer with specified _w_ as `io.Writer`, if
// `nil` it will set `os.Stdout` as _w_.
func NewStdPrinter(w io.Writer) *StdPrinter {

	if w == nil {
		w = os.Stdout
	}

	return &StdPrinter{
		entries: []common.TransactionEntry{},
		w:       w,
	}
}

// NewStdPrinterDefaults is the same as `NewStdPrinter` but it uses the built-in
// templates.
//
// If supplying two templates, the first is considered to be header
// template and the second is a line template. Otherwise if one argument it is
// considered to be a full template.
//
// A variation on two argument is when first argument is empty string and the second
// is a valid template, hence it will have no header but line template.
func NewStdPrinterDefaults(w io.Writer, templates ...string) *StdPrinter {

	if len(templates) > 2 || len(templates) < 1 {
		panic("must supply with one or two template arguments")
	}

	scp := NewStdPrinter(w)

	err := fs.WalkDir(
		embeddedTemplates, "templates",
		func(path string, d fs.DirEntry, err error) error {

			if err != nil || d.IsDir() {
				return err
			}

			file := filepath.Base(path)

			for i, t := range templates {

				if file == fmt.Sprintf("%s.gtpl", t) {

					data, err := fs.ReadFile(embeddedTemplates, path)
					if err != nil {
						panic(err)
					}

					if i == 0 {

						if len(templates) == 1 {
							scp.FullTemplate(string(data))
						} else {
							scp.HeaderTemplate(string(data))
						}

					} else {

						scp.LineTemplate(string(data))

					}

					break
				}

			}

			return nil
		})

	if err != nil {
		panic(err)
	}

	return scp
}

// LineTemplate will set the line template that will be used for each `Process`
// invocation.
//
// A line template will get one `common.TransactionEntry` at the time in it's context
// and hence may call any function on it.
func (scp *StdPrinter) LineTemplate(lineTemplate string) *StdPrinter {

	scp.lineTemplate = lineTemplate

	scp.output = template.Must(
		template.New("lineTemplate").
			Funcs(functions.Templatefuncs).
			Parse(lineTemplate),
	)

	return scp
}

// FullTemplate will set the full template that is used when all `common.TransactionEntry`
// instances has been collected (when `Flush` is invoked).
//
// Hence, the full template receives an array of `common.TransactionEntry` instead of
// one at the time (as with _lineTemplate_).
func (scp *StdPrinter) FullTemplate(fullTemplate string) *StdPrinter {

	scp.fullTemplate = fullTemplate

	scp.output = template.Must(
		template.New("fullTemplate").
			Funcs(functions.Templatefuncs).
			Parse(fullTemplate),
	)

	return scp
}

// HeaderTemplate sets the header template that is emitted before the first
// entry is emitted using the line template.
//
// If it is set and full template is used, it will be outputted before the
// full template when `Flush`. The template will receive the first available
// `common.TransactionEntry` to perform it's logic.
func (scp *StdPrinter) HeaderTemplate(headerTemplate string) *StdPrinter {

	scp.headerTemplate = headerTemplate

	scp.header = template.Must(
		template.New("fullTemplate").
			Funcs(functions.Templatefuncs).
			Parse(headerTemplate),
	)

	return scp
}

func (scp *StdPrinter) ProcessMany(tx []common.TransactionEntry) {

	for i := range tx {
		scp.Process(tx[i])
	}

}

func (scp *StdPrinter) Process(tx common.TransactionEntry) {

	scp.entries = append(scp.entries, tx.Clone())

	if scp.fullTemplate != "" {
		return
	}

	if scp.lineTemplate == "" {

		panic(
			"both full and line template is empty, need to have one set",
		)

	}

	if len(scp.entries) == 1 && scp.header != nil {

		if err := scp.header.Execute(scp.w, tx); err != nil {
			panic(err)
		}

	}

	if err := scp.output.Execute(scp.w, tx); err != nil {
		panic(err)
	}

}

func (scp *StdPrinter) Reset() {
	scp.entries = []common.TransactionEntry{}
}

func (scp *StdPrinter) Flush() []common.TransactionEntry {

	entries := scp.entries
	scp.Reset()

	if scp.fullTemplate != "" {

		if len(entries) > 1 && scp.header != nil {

			if err := scp.header.Execute(scp.w, entries[0]); err != nil {
				panic(err)
			}

		}

		if err := scp.output.Execute(scp.w, entries); err != nil {
			panic(err)
		}

	}

	return entries
}
