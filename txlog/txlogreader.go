package txlog

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
	postProcessor common.TxLogProcessor
}

func NewTxLogReader(postProcessor common.TxLogProcessor) *TxLogReaderImpl {

	return &TxLogReaderImpl{
		readers:       map[string]common.TransactionLogReader{},
		dir:           ".",
		recursive:     false,
		postProcessor: postProcessor,
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

	reader.SetExchange(name)
	lr.readers[name] = reader

	return lr

}

func (lr *TxLogReaderImpl) Read() []common.TransactionLog {

	return lr.preProcess(lr.read(lr.dir, lr.recursive))

}

func (lr *TxLogReaderImpl) ReadBuffer(readerName string, data []byte) []common.TransactionLog {

	if log, ok := lr.readers[readerName]; ok {
		return lr.preProcess(log.Unmarshal(data))
	}

	if lr.ignoreUnknown {
		return []common.TransactionLog{}
	}

	panic(fmt.Sprintf("could not find reader named: %s", readerName))

}

func (lr *TxLogReaderImpl) ReadBufferAsExchange(
	readerName string,
	data []byte,
) []common.TransactionLog {

	if log, ok := lr.readers[readerName]; ok {

		tx := lr.preProcess(log.Unmarshal(data))

		for i := range tx {
			tx[i].Exchange = readerName
		}

		return tx
	}

	if lr.ignoreUnknown {
		return []common.TransactionLog{}
	}

	panic(fmt.Sprintf("could not find reader named: %s", readerName))

}

func (lr *TxLogReaderImpl) read(directory string, recursive bool) []common.TransactionLog {

	tx := []common.TransactionLog{}

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

func (lr *TxLogReaderImpl) preProcess(logs []common.TransactionLog) []common.TransactionLog {

	if lr.postProcessor == nil {
		return logs
	}

	lr.postProcessor.Reset()

	lr.postProcessor.ProcessMany(logs)

	return lr.postProcessor.Flush()
}
