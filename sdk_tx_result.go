/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"context"
	"errors"
	"sync"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

type txResultDispatcher struct {
	nextBlockNum int64
	cc           *ChainClient
	stopC        chan struct{}
	txC          chan *common.Transaction

	mux             sync.Mutex // mux protect txRegistrations
	txRegistrations map[string]chan *common.Result
}

func newTxResultDispatcher(cc *ChainClient) *txResultDispatcher {
	return &txResultDispatcher{
		nextBlockNum:    -1,
		cc:              cc,
		stopC:           make(chan struct{}),
		txC:             make(chan *common.Transaction, 1),
		txRegistrations: make(map[string]chan *common.Result),
	}
}

// register registers for transaction result events.
// Note that unregister must be called when the registration is no longer needed.
// - txId is the transaction ID for which events are to be received
// - Returns the channel that is used to receive result. The channel
//   is closed when unregister is called.
func (d *txResultDispatcher) register(txId string) chan *common.Result {
	d.mux.Lock()
	defer d.mux.Unlock()
	if txResultC, exists := d.txRegistrations[txId]; exists {
		return txResultC
	}
	txResultC := make(chan *common.Result, 1)
	d.txRegistrations[txId] = txResultC
	return txResultC
}

// unregister removes the given registration and closes the event channel.
func (d *txResultDispatcher) unregister(txId string) {
	d.mux.Lock()
	defer d.mux.Unlock()
	if txResultC, exists := d.txRegistrations[txId]; exists {
		delete(d.txRegistrations, txId)
		close(txResultC)
	}
}

func (d *txResultDispatcher) start() {
	go d.autoSubscribe()
	for {
		select {
		case tx := <-d.txC:
			d.mux.Lock()
			if txResultC, exists := d.txRegistrations[tx.Payload.TxId]; exists {
				// non-blocking write to channel to ignore txResultC buffer is full in extreme cases
				select {
				case txResultC <- tx.Result:
				default:
				}
			}
			d.mux.Unlock()
		case <-d.stopC:
			return
		}
	}
}

func (d *txResultDispatcher) stop() {
	close(d.stopC)
}

func (d *txResultDispatcher) subscribe() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dataC, err := d.cc.SubscribeBlock(ctx, d.nextBlockNum, -1, false, false)
	if err != nil {
		return err
	}
	d.cc.logger.Debugf("txResultDispatcher subscribe success, block height %d", d.nextBlockNum)

	for {
		select {
		case block, ok := <-dataC:
			if !ok {
				return errors.New("chan is closed")
			}

			blockInfo := block.(*common.BlockInfo)
			d.cc.logger.Debugf("received block height: %d tx count: %d",
				blockInfo.Block.Header.BlockHeight, len(blockInfo.Block.Txs))
			for _, tx := range blockInfo.Block.Txs {
				d.cc.logger.Debugf("received tx %s", tx.Payload.TxId)
				d.txC <- tx
			}
			d.nextBlockNum = int64(blockInfo.Block.Header.BlockHeight) + 1
		case <-d.stopC:
			return nil
		}
	}
}

func (d *txResultDispatcher) autoSubscribe() {
	for {
		if err := d.subscribe(); err != nil {
			d.cc.logger.Debugf("txResultDispatcher subscribe failed, %s", err)
			d.cc.logger.Debug("txResultDispatcher will resubscribing after one second")
			time.Sleep(time.Second)
		} else {
			d.cc.logger.Debug("txResultDispatcher subscribe stopped")
			return
		}
	}
}
