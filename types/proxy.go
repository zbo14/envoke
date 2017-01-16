package types

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/tendermint/go-rpc/client"
	"github.com/tendermint/go-wire"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tndr "github.com/tendermint/tendermint/types"
)

const MinTxLength = 20
const MaxBlocks = 50

type Proxy struct {
	remote   string
	endpoint string

	rpc *rpcclient.ClientJSONRPC
	ws  *rpcclient.WSClient
}

func NewProxy(remote, endpoint string) *Proxy {
	return &Proxy{
		rpc:      rpcclient.NewClientJSONRPC(remote),
		ws:       rpcclient.NewWSClient(remote, endpoint),
		remote:   remote,
		endpoint: endpoint,
	}
}

// Txs

func (p *Proxy) BroadcastTx(mode string, tx tndr.Tx) (*ctypes.ResultBroadcastTx, error) {
	var err error
	result := new(ctypes.TMResult)
	switch mode {
	case "commit":
		_, err = p.rpc.Call("broadcast_tx_commit", []interface{}{tx}, result)
	case "sync":
		_, err = p.rpc.Call("broadcast_tx_sync", []interface{}{tx}, result)
	case "async":
		_, err = p.rpc.Call("broadcast_tx_async", []interface{}{tx}, result)
	default:
		err = errors.New("Unrecognized mode")
	}
	if err != nil {
		return nil, err
	}
	return (*result).(*ctypes.ResultBroadcastTx), nil
}

// Network

func (p *Proxy) GetStatus() (*ctypes.ResultStatus, error) {
	result := new(ctypes.TMResult)
	_, err := p.rpc.Call("status", []interface{}{}, result)
	if err != nil {
		return nil, err
	}
	return (*result).(*ctypes.ResultStatus), nil
}

// Consensus

func (p *Proxy) GetValidators() (*ctypes.ResultValidators, error) {
	result := new(ctypes.TMResult)
	_, err := p.rpc.Call("validators", nil, result)
	if err != nil {
		return nil, err
	}
	return (*result).(*ctypes.ResultValidators), nil
}

// Blockchain

func (p *Proxy) GetBlock(height int) (*ctypes.ResultBlock, error) {
	result := new(ctypes.TMResult)
	_, err := p.rpc.Call("block", []interface{}{height}, result)
	if err != nil {
		return nil, err
	}
	return (*result).(*ctypes.ResultBlock), nil
}

func (p *Proxy) GetChain(min, max int) (*ctypes.ResultBlockchainInfo, error) {
	result := new(ctypes.TMResult)
	if max-min < 0 {
		return nil, errors.New("maxHeight must be greater than minHeight")
	} else if max-min > MaxBlocks {
		return nil, errors.Errorf("you cannot query more than %d blocks at once", MaxBlocks)
	}
	_, err := p.rpc.Call("blockchain", []interface{}{min, max}, result)
	if err != nil {
		return nil, err
	}
	return (*result).(*ctypes.ResultBlockchainInfo), nil
}

// App query

func (p *Proxy) TMSPQuery(query []byte) (*ctypes.ResultTMSPQuery, error) {
	result := new(ctypes.TMResult)
	_, err := p.rpc.Call("tmsp_query", []interface{}{query}, result)
	if err != nil {
		return nil, err
	}
	return (*result).(*ctypes.ResultTMSPQuery), nil
}

// Websocket

func (p *Proxy) StartWS() (err error) {
	started, err := p.ws.Start()
	if err != nil {
		return err
	} else if !started {
		return errors.New("Failed to start websocket client")
	}
	return nil
}

func (p *Proxy) StopWS() error {
	stopped := p.ws.Stop()
	if !stopped {
		return errors.New("Failed to stop websocket")
	}
	return nil
}

func (p *Proxy) ReadResult(event string, evData tndr.TMEventData) (tndr.TMEventData, error) {
	select {
	case data := <-p.ws.ResultsCh:
		var items []interface{}
		err := json.Unmarshal(data, &items)
		if err != nil {
			return nil, err
		}
		var result = ctypes.ResultEvent{event, evData}
		wire.ReadJSONObject(&result, items[1], &err)
		if err != nil {
			return nil, err
		}
		if result.Name != event {
			return nil, errors.New("Wrong event type")
		}
		return result.Data, nil
	case err := <-p.ws.ErrorsCh:
		return nil, err
	}
}

func (p *Proxy) SubscribeNewBlock() error {
	eid := tndr.EventStringNewBlock()
	return p.ws.Subscribe(eid)
}

func (p *Proxy) UnsubscribeNewBlock() error {
	eid := tndr.EventStringNewBlock()
	return p.ws.Unsubscribe(eid)
}

func (p *Proxy) WriteWS(mode string, v interface{}) error {
	switch mode {
	case "json":
		return p.ws.WriteJSON(v)
	case "text":
		data := v.([]byte)
		return p.ws.WriteMessage(websocket.TextMessage, data)
	case "binary":
		data := v.([]byte)
		return p.ws.WriteMessage(websocket.BinaryMessage, data)
	default:
		return errors.New("Unrecognized write mode")
	}
}
