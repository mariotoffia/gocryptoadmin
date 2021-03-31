package transactions

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mariotoffia/gocryptoadmin/transactions/txcommon"
)

type TxLogReaderImpl struct {
	readers       map[string]txcommon.TransactionLogReader
	dir           string
	recursive     bool
	ignoreUnknown bool
	secwindow     time.Duration
}

func NewTxLogReader() *TxLogReaderImpl {

	return &TxLogReaderImpl{
		readers:   map[string]txcommon.TransactionLogReader{},
		dir:       ".",
		recursive: false,
		secwindow: time.Duration(5 * 60),
	}

}

func (lr *TxLogReaderImpl) UseWindowSize(seconds int64) *TxLogReaderImpl {
	lr.secwindow = time.Duration(seconds)
	return lr
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
	reader txcommon.TransactionLogReader) *TxLogReaderImpl {

	lr.readers[name] = reader
	return lr

}

func (lr *TxLogReaderImpl) Read() []txcommon.Transaction {

	tx := lr.read(lr.dir, lr.recursive)
	return lr.preProcess(tx)
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

func (lr *TxLogReaderImpl) logReaderFromFileName(name string) txcommon.TransactionLogReader {

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

func (lr *TxLogReaderImpl) preProcess(logs []txcommon.Transaction) []txcommon.Transaction {

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].CreatedAt.Before(logs[j].CreatedAt)
	})

	type processor struct {
		next        time.Time
		groupId     int64
		intervalAcc txcommon.GroupAccumulate
	}

	groupId := int64(0)
	products := map[string]processor{}
	lastSide := map[string]txcommon.SideType{}

	for i, tx := range logs {

		if _, ok := products[tx.Exchange+tx.Asset+string(tx.Side)]; !ok {

			lastSide[tx.Exchange+tx.Asset] = tx.Side

			groupId++

			products[tx.Exchange+tx.Asset+string(tx.Side)] = processor{
				groupId: groupId,
				next:    tx.CreatedAt.Add(time.Second * lr.secwindow),
				intervalAcc: txcommon.GroupAccumulate{
					GrpSize:  tx.Size,
					GrpFee:   tx.Fee,
					GrpTotal: tx.Total,
				},
			}

			logs[i].GrpSize = tx.Size
			logs[i].GrpTotal = tx.Total
			logs[i].GrpFee = tx.Fee
			logs[i].GroupID = groupId

			continue
		}

		proc := products[tx.Exchange+tx.Asset+string(tx.Side)]

		proc.intervalAcc.GrpSize = proc.intervalAcc.GrpSize + tx.Size
		proc.intervalAcc.GrpTotal = proc.intervalAcc.GrpTotal + tx.Cost.Total
		proc.intervalAcc.GrpFee = proc.intervalAcc.GrpFee + tx.Cost.Fee

		if lastSide[tx.Exchange+tx.Asset] == tx.Side &&
			tx.CreatedAt.Before(proc.next) {

		} else {

			// Next slice
			groupId++
			proc.groupId = groupId
			proc.intervalAcc.GrpSize = tx.Size
			proc.intervalAcc.GrpFee = tx.Fee
			proc.intervalAcc.GrpTotal = tx.Total
			proc.next = tx.CreatedAt.Add(time.Second * lr.secwindow)

			lastSide[tx.Exchange+tx.Asset] = tx.Side // Might have changed

		}

		logs[i].GrpSize = proc.intervalAcc.GrpSize
		logs[i].GrpFee = proc.intervalAcc.GrpFee
		logs[i].GrpTotal = proc.intervalAcc.GrpTotal
		logs[i].GroupID = proc.groupId

		products[tx.Exchange+tx.Asset+string(tx.Side)] = proc

	}

	return logs
}
