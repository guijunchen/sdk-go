/**
 * @Author: jasonruan
 * @Date:   2020-12-03 11:31:33
 **/
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"context"
)

func (cc ChainClient) SubscribeBlock(ctx context.Context, payloadBytes []byte) (<-chan interface{}, error) {
	panic("implement me")
}

func (cc ChainClient) SubscribeTx(ctx context.Context, payloadBytes []byte) (<-chan interface{}, error) {
	panic("implement me")
}

func (cc ChainClient) Subscribe(ctx context.Context, txType pb.TxType, payloadBytes []byte) (<-chan interface{}, error) {
	panic("implement me")
}
