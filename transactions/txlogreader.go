package transactions

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mariotoffia/gocryptoadmin/common"
)

type TxLogReaderImpl struct {
	readers       map[string]common.TransactionLogReader
	dir           string
	recursive     bool
	ignoreUnknown bool
}

func NewTxLogReader() *TxLogReaderImpl {

	return &TxLogReaderImpl{
		readers:   map[string]common.TransactionLogReader{},
		dir:       ".",
		recursive: false,
	}

}

func (lr *TxLogReaderImpl) IsRecursive() *TxLogReaderImpl {

	lr.recursive = true
	return lr

}

func (lr *TxLogReaderImpl) IgnoreUnknownFiles() *TxLogReaderImpl {

	lr.ignoreUnknown = true
	return lr

}

func (lr *TxLogReaderImpl) UseDir(dir string) *TxLogReaderImpl {

	lr.dir = dir
	return lr

}

func (lr *TxLogReaderImpl) RegisterReader(
	name string,
	reader common.TransactionLogReader) *TxLogReaderImpl {

	lr.readers[name] = reader
	return lr

}

func (lr *TxLogReaderImpl) Read() []common.Transaction {

	return lr.read(lr.dir, lr.recursive)

}

func (lr *TxLogReaderImpl) ReadBuffer(readerName string, data []byte) []common.Transaction {

	if lr, ok := lr.readers[readerName]; ok {
		return lr.Unmarshal(data)
	}

	if lr.ignoreUnknown {
		return []common.Transaction{}
	}

	panic(fmt.Sprintf("could not find reader named: %s", readerName))

}

func (lr *TxLogReaderImpl) read(directory string, recursive bool) []common.Transaction {

	tx := []common.Transaction{}

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

		if log == nil {
			continue
		}

		data, err := ioutil.ReadFile(filepath.Join(directory, file.Name()))

		if err != nil {
			panic(err)
		}

		tx = append(tx, log.Unmarshal(data)...)
	}

	return tx

}

func (lr *TxLogReaderImpl) logReaderFromFileName(name string) common.TransactionLogReader {

	if lr, ok := lr.readers[logReaderNameFromFileName(name)]; ok {
		return lr
	}

	if lr.ignoreUnknown {
		return nil
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
