package transactions

import (
	"sort"
	"time"

	"github.com/mariotoffia/gocryptoadmin/transactions/txcommon"
)

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

		size := tx.Size
		if tx.Side == txcommon.SideTypeSell {
			size = -tx.Size
		}

		if _, ok := products[tx.Exchange+tx.Asset+string(tx.Side)]; !ok {

			lastSide[tx.Exchange+tx.Asset] = tx.Side

			groupId++

			products[tx.Exchange+tx.Asset+string(tx.Side)] = processor{
				groupId: groupId,
				next:    tx.CreatedAt.Add(time.Second * lr.secwindow),
				intervalAcc: txcommon.GroupAccumulate{
					GrpSize:  size,
					GrpFee:   tx.Fee,
					GrpTotal: tx.Total,
				},
			}

			logs[i].GrpSize = size
			logs[i].GrpTotal = tx.Total
			logs[i].GrpFee = tx.Fee
			logs[i].GroupID = groupId

			continue
		}

		proc := products[tx.Exchange+tx.Asset+string(tx.Side)]

		proc.intervalAcc.GrpSize = proc.intervalAcc.GrpSize + size
		proc.intervalAcc.GrpTotal = proc.intervalAcc.GrpTotal + tx.Cost.Total
		proc.intervalAcc.GrpFee = proc.intervalAcc.GrpFee + tx.Cost.Fee

		if lastSide[tx.Exchange+tx.Asset] == tx.Side &&
			tx.CreatedAt.Before(proc.next) {

		} else {

			// Next slice
			groupId++
			proc.groupId = groupId
			proc.intervalAcc.GrpSize = size
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
