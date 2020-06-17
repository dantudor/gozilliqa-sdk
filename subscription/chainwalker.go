/*
 * Copyright (C) 2019 Zilliqa
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package subscription

import (
	"context"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/workpool"
	"github.com/ybbus/jsonrpc"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Walker struct {
	Provider     *provider.Provider
	FromBlock    uint64
	ToBlock      uint64
	CurrentBlock uint64
	Address      string
	EventLogs    map[uint64]Log
	WorkerNumber int64
	EventName    string
	Retry        int
	Interval     int64
}

type Log struct {
	Hash      string
	EventName string
	Address   string
	Logs      interface{}
}

func NewWalker(p *provider.Provider, from, to uint64, address string, workerNumber int64, eventName string, retry int, interval int64) *Walker {
	eventLogs := make(map[uint64]Log)
	return &Walker{
		Provider:     p,
		FromBlock:    from,
		ToBlock:      to,
		Address:      address,
		EventLogs:    eventLogs,
		WorkerNumber: workerNumber,
		EventName:    eventName,
		Retry:        retry,
		Interval:     interval,
	}
}

type GetEventReceiptTask struct {
	Provider *provider.Provider
	Id       string
	Complete *Complete
	Walker   *Walker
	BlockNum uint64
}

type Complete struct {
	sync.Mutex
	Number int
}

func (t GetEventReceiptTask) UUID() string {
	return t.Id
}

func NewGetReceiptTask(tx string, provider2 *provider.Provider, c *Complete, w *Walker, b uint64) GetEventReceiptTask {
	return GetEventReceiptTask{
		Id:       tx,
		Provider: provider2,
		Complete: c,
		Walker:   w,
		BlockNum: b,
	}
}

func (t GetEventReceiptTask) handleTxn(rsp *jsonrpc.RPCResponse) {
	resultMap := rsp.Result.(map[string]interface{})
	receipt := resultMap["receipt"].(map[string]interface{})
	addr := resultMap["toAddr"].(string)
	success, ok1 := receipt["success"]
	if success == nil || success.(bool) == false {
		return
	} else {
		eventLogs, ok := receipt["event_logs"]
		if ok {
			els := eventLogs.([]interface{})
			for _, el := range els {
				log := el.(map[string]interface{})
				eventName, ok2 := log["_eventname"]
				// important: currently we only compare contract address to toAddr
				if ok1 && ok2 && strings.Compare(strings.ToLower(addr), strings.ToLower(t.Walker.Address[2:])) == 0 && strings.Compare(eventName.(string), t.Walker.EventName) == 0 {
					logData := Log{
						Hash:      t.Id,
						EventName: eventName.(string),
						Address:   t.Walker.Address,
						Logs:      log,
					}
					t.Walker.EventLogs[t.BlockNum] = logData
				}
			}
		}
	}
}

func (t GetEventReceiptTask) Run() {
	t.Complete.Lock()
	defer t.Complete.Unlock()
	t.Complete.Number++
	//fmt.Println("start to get transaction " + t.Id)
	var rpcResult *jsonrpc.RPCResponse
	for i := 0; i < t.Walker.Retry; i++ {
		rsp, err := t.Provider.GetTransaction(t.Id)
		if err != nil || rsp.Error != nil {
			if err != nil {
				fmt.Printf("get transaction failed, id = %s, error = %s\n", t.Id, err.Error())
			} else {
				fmt.Printf("get transaction failed, id = %s, error = %s\n", t.Id, rsp.Error.Error())
			}
			fmt.Println("start to retry")
			time.Sleep(time.Millisecond * time.Duration(t.Walker.Interval))
		} else {
			rpcResult = rsp
			break
		}
	}

	if rpcResult == nil {
		fmt.Printf("still cannot get transaction = %s\n",t.Id)
	} else {
		t.handleTxn(rpcResult)
	}
}

func (w *Walker) StartTraversalBlock() {
	for i := w.FromBlock; i < w.ToBlock; i++ {
		var rpcResult *jsonrpc.RPCResponse
		for j := 0; j < w.Retry; j++ {
			rsp, err := w.Provider.GetTransactionsForTxBlock(strconv.FormatUint(i, 10))
			if err != nil || rsp.Error != nil {
				if err != nil {
					fmt.Printf("get block failed, number = %d, error = %s\n", i, err.Error())
				} else if rsp.Error.Message == "TxBlock has no transactions" {
					rsp = &jsonrpc.RPCResponse{
						JSONRPC: "empty",
					}
					rpcResult = rsp
					break
				} else {
					fmt.Printf("get block failed, number = %d, error = %s\n", i, rsp.Error.Error())
				}
				fmt.Println("start to retry")
				time.Sleep(time.Millisecond * time.Duration(w.Interval))
			} else {
				rpcResult = rsp
				break
			}
		}

		if rpcResult == nil {
			fmt.Printf("still cannot get result, block = %d\n", i)
			continue
		}

		if rpcResult.JSONRPC == "empty" {
			fmt.Printf("block is empty, block = %d\n", i)
			continue
		}

		txResult := rpcResult.Result.([]interface{})
		var txns []string

		// flat tx hash
		for _, txs := range txResult {
			if txs == nil {
				continue
			}
			txList := txs.([]interface{})
			if len(txList) > 0 {
				for _, tx := range txList {
					txns = append(txns, tx.(string))
				}
			} else {
				continue
			}
		}

		// get detail
		fmt.Println("tx for block ", i, " = ", txns)
		complete := &Complete{
			Number: 0,
		}

		wp := workpool.NewWorkPool(w.WorkerNumber)
		quit := make(chan int, 1)
		for _, tx := range txns {
			task := NewGetReceiptTask(tx, w.Provider, complete, w, i)
			wp.AddTask(task)
		}
		wp.Poll(context.TODO(), quit)
		select {
		case <-quit:
			break
		}

	}
}
