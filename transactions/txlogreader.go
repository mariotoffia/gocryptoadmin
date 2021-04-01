package transactions

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

type TxLogReaderImpl struct {
	readers       map[string]common.TransactionLogReader
	dir           string
	recursive     bool
	ignoreUnknown bool
	secwindow     time.Duration
}

func NewTxLogReader() *TxLogReaderImpl {

	return &TxLogReaderImpl{
		readers:   map[string]common.TransactionLogReader{},
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
	reader common.TransactionLogReader) *TxLogReaderImpl {

	lr.readers[name] = reader
	return lr

}

func (lr *TxLogReaderImpl) Read() []common.Transaction {

	tx := lr.read(lr.dir, lr.recursive)
	return lr.preProcess(tx)
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

func (lr *TxLogReaderImpl) preProcess(logs []common.Transaction) []common.Transaction {

	sort.Slice(logs, func(i, j int) bool {

		if logs[i].CreatedAt.Equal(logs[j].CreatedAt) {

			if logs[i].Exchange == logs[j].Exchange {

				if logs[i].Exchange == "coinbase-pro" &&
					logs[i].AssetPair == logs[j].AssetPair {

					li, err := strconv.ParseInt(logs[i].ID, 10, 64)
					if err != nil {
						panic(err)
					}

					lj, err := strconv.ParseInt(logs[j].ID, 10, 64)
					if err != nil {
						panic(err)
					}

					return li < lj

				}
			}

		}

		return logs[i].CreatedAt.Before(logs[j].CreatedAt)
	})

	type processor struct {
		next        time.Time
		groupId     int64
		intervalAcc common.GroupAccumulate
	}

	groupId := int64(0)
	products := map[string]processor{}
	lastSide := map[string]common.SideType{}

	for i, tx := range logs {

		if _, ok := products[tx.Exchange+tx.AssetPair+string(tx.Side)]; !ok {

			lastSide[tx.Exchange+tx.AssetPair] = tx.Side

			groupId++

			products[tx.Exchange+tx.AssetPair+string(tx.Side)] = processor{
				groupId: groupId,
				next:    tx.CreatedAt.Add(time.Second * lr.secwindow),
				intervalAcc: common.GroupAccumulate{
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

		proc := products[tx.Exchange+tx.AssetPair+string(tx.Side)]

		proc.intervalAcc.GrpSize = utils.ToFixed(proc.intervalAcc.GrpSize+tx.Size, 8)
		proc.intervalAcc.GrpTotal = utils.ToFixed(proc.intervalAcc.GrpTotal+tx.Cost.Total, 8)
		proc.intervalAcc.GrpFee = utils.ToFixed(proc.intervalAcc.GrpFee+tx.Cost.Fee, 8)

		if lastSide[tx.Exchange+tx.AssetPair] == tx.Side &&
			tx.CreatedAt.Before(proc.next) {

		} else {

			// Next slice
			groupId++
			proc.groupId = groupId
			proc.intervalAcc.GrpSize = tx.Size
			proc.intervalAcc.GrpFee = tx.Fee
			proc.intervalAcc.GrpTotal = tx.Total
			proc.next = tx.CreatedAt.Add(time.Second * lr.secwindow)

			lastSide[tx.Exchange+tx.AssetPair] = tx.Side // Might have changed

		}

		logs[i].GrpSize = proc.intervalAcc.GrpSize
		logs[i].GrpFee = proc.intervalAcc.GrpFee
		logs[i].GrpTotal = proc.intervalAcc.GrpTotal
		logs[i].GroupID = proc.groupId

		products[tx.Exchange+tx.AssetPair+string(tx.Side)] = proc

	}

	return logs
}
