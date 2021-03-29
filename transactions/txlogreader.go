package transactions

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mariotoffia/gocryptoadmin/transactions/txcommon"
)

type TxLogReaderImpl struct {
	readers   map[string]txcommon.TransactionLogReader
	dir       string
	recursive bool
	sort      bool
}

func NewTxLogReader() *TxLogReaderImpl {

	return &TxLogReaderImpl{
		readers:   map[string]txcommon.TransactionLogReader{},
		dir:       ".",
		recursive: false,
	}

}

func (lr *TxLogReaderImpl) IsRecursive() *TxLogReaderImpl {

	lr.recursive = true
	return lr

}

func (lr *TxLogReaderImpl) UseDir(dir string) *TxLogReaderImpl {

	lr.dir = dir
	return lr

}

func (lr *TxLogReaderImpl) SortRead() *TxLogReaderImpl {

	lr.sort = true
	return lr

}

func (lr *TxLogReaderImpl) RegisterReader(
	name string,
	reader txcommon.TransactionLogReader) *TxLogReaderImpl {

	lr.readers[name] = reader
	return lr

}

func (lr *TxLogReaderImpl) Read() []txcommon.Transaction {

	tx := lr.read(lr.dir, lr.recursive)

	if lr.sort {
		return SortLogs(tx)
	}

	return tx
}

func (lr *TxLogReaderImpl) read(directory string, recursive bool) []txcommon.Transaction {

	tx := []txcommon.Transaction{}

	if !filepath.IsAbs(directory) {

		var err error
		directory, err = filepath.Abs(directory)

		if err != nil {
			panic(err)
		}

	}

	files, err := ioutil.ReadDir(directory)

	if err != nil {
		panic(err)
	}

	for _, file := range files {

		if file.IsDir() {

			if !recursive {
				continue
			}

			tx = append(tx, lr.read(file.Name(), recursive)...)
		}

		if !strings.HasSuffix(file.Name(), ".csv") {
			continue
		}

		log := lr.logReaderFromFileName(file.Name())

		data, err := ioutil.ReadFile(filepath.Join(directory, file.Name()))

		if err != nil {
			panic(err)
		}

		tx = append(tx, log.Unmarshal(data)...)
	}

	return tx

}

func (lr *TxLogReaderImpl) logReaderFromFileName(name string) txcommon.TransactionLogReader {

	if lr, ok := lr.readers[logReaderNameFromFileName(name)]; ok {
		return lr
	}

	panic(
		fmt.Sprintf("Could not find logreader from file: %s, extracted lr name: %s",
			name, logReaderNameFromFileName(name),
		),
	)
}
func logReaderNameFromFileName(name string) string {
	return strings.SplitN(name, "_", 2)[0]
}
