// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/types/jsonpb"
	"github.com/stretchr/testify/assert"
)

func TestAllowExecName(t *testing.T) {
	//allow exec list
	old := AllowUserExec
	defer func() {
		AllowUserExec = old
	}()
	AllowUserExec = nil
	AllowUserExec = append(AllowUserExec, []byte("coins"))
	isok := IsAllowExecName([]byte("a"), []byte("a"))
	assert.Equal(t, isok, false)

	isok = IsAllowExecName([]byte("coins"), []byte("coins"))
	assert.Equal(t, isok, true)

	isok = IsAllowExecName([]byte("coins"), []byte("user.coins"))
	assert.Equal(t, isok, true)

	isok = IsAllowExecName([]byte("coins"), []byte("user.coinsx"))
	assert.Equal(t, isok, false)

	isok = IsAllowExecName([]byte("coins"), []byte("user.coins.evm2"))
	assert.Equal(t, isok, true)

	isok = IsAllowExecName([]byte("coins"), []byte("user.p.guodun.coins.evm2"))
	assert.Equal(t, isok, false)

	isok = IsAllowExecName([]byte("coins"), []byte("user.p.guodun.coins"))
	assert.Equal(t, isok, true)

	isok = IsAllowExecName([]byte("coins"), []byte("user.p.guodun.user.coins"))
	assert.Equal(t, isok, true)

	isok = IsAllowExecName([]byte("#coins"), []byte("user.p.guodun.user.coins"))
	assert.Equal(t, isok, false)

	isok = IsAllowExecName([]byte("coins-"), []byte("user.p.guodun.user.coins"))
	assert.Equal(t, isok, false)
}

func BenchmarkExecName(b *testing.B) {
	cfg := NewChain33Config(GetDefaultCfgstring())
	for i := 0; i < b.N; i++ {
		cfg.ExecName("hello")
	}
}

func BenchmarkG(b *testing.B) {
	cfg := NewChain33Config(GetDefaultCfgstring())
	for i := 0; i < b.N; i++ {
		cfg.G("TestNet")
	}
}

func BenchmarkS(b *testing.B) {
	cfg := NewChain33Config(GetDefaultCfgstring())
	for i := 0; i < b.N; i++ {
		cfg.S("helloword", true)
	}
}
func TestJsonNoName(t *testing.T) {
	flag := int32(1)
	params := struct {
		Flag int32
	}{
		Flag: flag,
	}
	data, err := json.Marshal(params)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, string(data), "{\"Flag\":1}")
}

func TestNil(t *testing.T) {
	v := reflect.ValueOf(nil)
	assert.Equal(t, v.IsValid(), false)
}

func TestProtoToJson(t *testing.T) {
	r := &Reply{}
	b, err := json.Marshal(r)
	assert.Nil(t, err)
	assert.Equal(t, b, []byte(`{}`))

	encode := &jsonpb.Marshaler{EmitDefaults: true}
	s, err := encode.MarshalToString(r)
	assert.Nil(t, err)
	assert.Equal(t, s, `{"isOk":false,"msg":null}`)
	var dr Reply
	err = jsonpb.UnmarshalString(`{"isOk":false,"msg":null}`, &dr)
	assert.Nil(t, err)
	assert.Nil(t, dr.Msg)
	encode2 := &jsonpb.Marshaler{EmitDefaults: false}
	s, err = encode2.MarshalToString(r)
	assert.Nil(t, err)
	assert.Equal(t, s, `{}`)

	r = &Reply{Msg: []byte("OK")}
	b, err = json.Marshal(r)
	assert.Nil(t, err)
	assert.Equal(t, b, []byte(`{"msg":"T0s="}`))

	encode = &jsonpb.Marshaler{EmitDefaults: true}
	s, err = encode.MarshalToString(r)
	assert.Nil(t, err)
	assert.Equal(t, s, `{"isOk":false,"msg":"0x4f4b"}`)

	err = jsonpb.UnmarshalString(`{"isOk":false,"msg":"0x4f4b"}`, &dr)
	assert.Nil(t, err)
	assert.Equal(t, dr.Msg, []byte("OK"))

	err = jsonpb.UnmarshalString(`{"isOk":false,"msg":"4f4b"}`, &dr)
	assert.Equal(t, err, jsonpb.ErrBytesFormat)

	err = jsonpb.UnmarshalString(`{"isOk":false,"msg":"0x"}`, &dr)
	assert.Nil(t, err)
	assert.Equal(t, dr.Msg, []byte(""))

	err = jsonpb.UnmarshalString(`{"isOk":false,"msg":"str://OK"}`, &dr)
	assert.Nil(t, err)
	assert.Equal(t, dr.Msg, []byte("OK"))

	err = jsonpb.UnmarshalString(`{"isOk":false,"msg":"str://0"}`, &dr)
	assert.Nil(t, err)
	assert.Equal(t, dr.Msg, []byte("0"))

	r = &Reply{Msg: []byte{}}
	b, err = json.Marshal(r)
	assert.Nil(t, err)
	assert.Equal(t, b, []byte(`{}`))

	encode = &jsonpb.Marshaler{EmitDefaults: true}
	s, err = encode.MarshalToString(r)
	assert.Nil(t, err)
	assert.Equal(t, s, `{"isOk":false,"msg":""}`)

	err = jsonpb.UnmarshalString(`{"isOk":false,"msg":""}`, &dr)
	assert.Nil(t, err)
	assert.Equal(t, dr.Msg, []byte{})
}

func TestJsonpbUTF8(t *testing.T) {
	r := &Reply{Msg: []byte("OK")}
	b, err := PBToJSONUTF8(r)
	assert.Nil(t, err)
	assert.Equal(t, b, []byte(`{"isOk":false,"msg":"OK"}`))

	var newreply Reply
	err = JSONToPBUTF8(b, &newreply)
	assert.Nil(t, err)
	assert.Equal(t, r, &newreply)
}

func TestJsonpb(t *testing.T) {
	r := &Reply{Msg: []byte("OK")}
	b, err := PBToJSON(r)
	assert.Nil(t, err)
	assert.Equal(t, b, []byte(`{"isOk":false,"msg":"0x4f4b"}`))

	var newreply Reply
	err = JSONToPB(b, &newreply)
	assert.Nil(t, err)
	assert.Equal(t, r, &newreply)
}

func TestHex(t *testing.T) {
	s := "0x4f4b"
	b, err := common.FromHex(s)
	assert.Nil(t, err)
	assert.Equal(t, b, []byte("OK"))
}

func TestGetLogName(t *testing.T) {
	name := GetLogName([]byte("xxx"), 0)
	assert.Equal(t, "LogReserved", name)
	assert.Equal(t, "LogErr", GetLogName([]byte("coins"), 1))
	assert.Equal(t, "LogFee", GetLogName([]byte("token"), 2))
	assert.Equal(t, "LogReserved", GetLogName([]byte("xxxx"), 100))
}

func TestDecodeLog(t *testing.T) {
	data, _ := common.FromHex("0x0a2b10c0c599b78c1d2222314c6d7952616a4e44686f735042746259586d694c466b5174623833673948795565122b1080ab8db78c1d2222314c6d7952616a4e44686f735042746259586d694c466b5174623833673948795565")
	l, err := DecodeLog([]byte("xxx"), 2, data)
	assert.Nil(t, err)
	j, err := json.Marshal(l)
	assert.Nil(t, err)
	assert.Equal(t, "{\"prev\":{\"balance\":999769400000,\"addr\":\"1LmyRajNDhosPBtbYXmiLFkQtb83g9HyUe\"},\"current\":{\"balance\":999769200000,\"addr\":\"1LmyRajNDhosPBtbYXmiLFkQtb83g9HyUe\"}}", string(j))
}

func TestGetRealExecName(t *testing.T) {
	a := []struct {
		key     string
		realkey string
	}{
		{"coins", "coins"},
		{"user.p.coins", "user.p.coins"},
		{"user.p.guodun.coins", "coins"},
		{"user.evm.hash", "evm"},
		{"user.p.para.evm.hash", "evm.hash"},
		{"user.p.para.user.evm.hash", "evm"},
		{"user.p.para.", "user.p.para."},
	}
	for _, v := range a {
		assert.Equal(t, string(GetRealExecName([]byte(v.key))), v.realkey)
	}
}

func genPrefixEdge(prefix []byte) (r []byte) {
	for j := 0; j < len(prefix); j++ {
		r = append(r, prefix[j])
	}

	i := len(prefix) - 1
	for i >= 0 {
		if r[i] < 0xff {
			r[i]++
			break
		} else {
			i--
		}
	}

	return r
}

func (t *StoreListReply) IterateCallBack(key, value []byte) bool {
	if t.Mode == 1 { //[start, end)
		if t.Num >= t.Count {
			t.NextKey = key
			return true
		}
		t.Num++
		t.Keys = append(t.Keys, cloneByte(key))
		t.Values = append(t.Values, cloneByte(value))
		return false
	} else if t.Mode == 2 { //prefix + suffix
		if len(key) > len(t.Suffix) {
			if string(key[len(key)-len(t.Suffix):]) == string(t.Suffix) {
				t.Num++
				t.Keys = append(t.Keys, cloneByte(key))
				t.Values = append(t.Values, cloneByte(value))
				if t.Num >= t.Count {
					t.NextKey = key
					return true
				}
			}
			return false
		}
		return false
	} else {
		fmt.Println("StoreListReply.IterateCallBack unsupported mode", "mode", t.Mode)
		return true
	}
}

func cloneByte(v []byte) []byte {
	value := make([]byte, len(v))
	copy(value, v)
	return value
}

func TestIterateCallBack_PrefixWithoutExecAddr(t *testing.T) {
	key := "mavl-coins-bty-exec-16htvcBNSEA7fZhAdLJphDwQRQJaHpyHTp:1JmFaA6unrCFYEWPGRi7uuXY1KthTJxJEP"
	//prefix1 := "mavl-coins-bty-exec-16htvcBNSEA7fZhAdLJphDwQRQJaHpyHTp:"
	prefix2 := "mavl-coins-bty-exec-"
	//execAddr := "16htvcBNSEA7fZhAdLJphDwQRQJaHpyHTp"
	addr := "1JmFaA6unrCFYEWPGRi7uuXY1KthTJxJEP"

	var reply = &StoreListReply{
		Start:  []byte(prefix2),
		End:    genPrefixEdge([]byte(prefix2)),
		Suffix: []byte(addr),
		Mode:   int64(2),
		Count:  int64(100),
	}

	var acc = &Account{
		Currency: 0,
		Balance:  1,
		Frozen:   1,
		Addr:     addr,
	}

	value := Encode(acc)

	bRet := reply.IterateCallBack([]byte(key), value)
	assert.Equal(t, false, bRet)
	assert.Equal(t, 1, len(reply.Keys))
	assert.Equal(t, 1, len(reply.Values))
	assert.Equal(t, int64(1), reply.Num)
	assert.Equal(t, 0, len(reply.NextKey))

	bRet = reply.IterateCallBack([]byte(key), value)
	assert.Equal(t, false, bRet)
	assert.Equal(t, 2, len(reply.Keys))
	assert.Equal(t, 2, len(reply.Values))
	assert.Equal(t, int64(2), reply.Num)
	assert.Equal(t, 0, len(reply.NextKey))

	key2 := "mavl-coins-bty-exec-16htvcBNSEA7fZhAdLJphDwQRQJaHpyHTp:2JmFaA6unrCFYEWPGRi7uuXY1KthTJxJEP"
	bRet = reply.IterateCallBack([]byte(key2), value)
	assert.Equal(t, false, bRet)
	assert.Equal(t, 2, len(reply.Keys))
	assert.Equal(t, 2, len(reply.Values))
	assert.Equal(t, int64(2), reply.Num)
	assert.Equal(t, 0, len(reply.NextKey))

	key3 := "mavl-coins-bty-exec-26htvcBNSEA7fZhAdLJphDwQRQJaHpyHTp:1JmFaA6unrCFYEWPGRi7uuXY1KthTJxJEP"
	bRet = reply.IterateCallBack([]byte(key3), value)
	assert.Equal(t, false, bRet)
	assert.Equal(t, 3, len(reply.Keys))
	assert.Equal(t, 3, len(reply.Values))
	assert.Equal(t, int64(3), reply.Num)
	assert.Equal(t, 0, len(reply.NextKey))

	reply.Count = int64(4)

	bRet = reply.IterateCallBack([]byte(key3), value)
	assert.Equal(t, true, bRet)
	assert.Equal(t, 4, len(reply.Keys))
	assert.Equal(t, 4, len(reply.Values))
	assert.Equal(t, int64(4), reply.Num)
	assert.Equal(t, key3, string(reply.NextKey))
	fmt.Println(string(reply.NextKey))
}

func TestIterateCallBack_PrefixWithExecAddr(t *testing.T) {
	key := "mavl-coins-bty-exec-16htvcBNSEA7fZhAdLJphDwQRQJaHpyHTp:1JmFaA6unrCFYEWPGRi7uuXY1KthTJxJEP"
	prefix1 := "mavl-coins-bty-exec-16htvcBNSEA7fZhAdLJphDwQRQJaHpyHTp:"
	//execAddr := "16htvcBNSEA7fZhAdLJphDwQRQJaHpyHTp"
	addr := "1JmFaA6unrCFYEWPGRi7uuXY1KthTJxJEP"

	var reply = &StoreListReply{
		Start:  []byte(prefix1),
		End:    genPrefixEdge([]byte(prefix1)),
		Suffix: []byte(addr),
		Mode:   int64(2),
		Count:  int64(1),
	}

	var acc = &Account{
		Currency: 0,
		Balance:  1,
		Frozen:   1,
		Addr:     addr,
	}

	value := Encode(acc)

	key2 := "mavl-coins-bty-exec-16htvcBNSEA7fZhAdLJphDwQRQJaHpyHTp:2JmFaA6unrCFYEWPGRi7uuXY1KthTJxJEP"
	bRet := reply.IterateCallBack([]byte(key2), value)
	assert.Equal(t, false, bRet)
	assert.Equal(t, 0, len(reply.Keys))
	assert.Equal(t, 0, len(reply.Values))
	assert.Equal(t, int64(0), reply.Num)
	assert.Equal(t, 0, len(reply.NextKey))

	bRet = reply.IterateCallBack([]byte(key), value)
	assert.Equal(t, true, bRet)
	assert.Equal(t, 1, len(reply.Keys))
	assert.Equal(t, 1, len(reply.Values))
	assert.Equal(t, int64(1), reply.Num)
	assert.Equal(t, len(key), len(reply.NextKey))

	//key2 := "mavl-coins-bty-exec-16htvcBNSEA7fZhAdLJphDwQRQJaHpyHTp:2JmFaA6unrCFYEWPGRi7uuXY1KthTJxJEP"
	reply.NextKey = nil
	reply.Count = int64(2)
	bRet = reply.IterateCallBack([]byte(key2), value)
	assert.Equal(t, false, bRet)
	assert.Equal(t, 1, len(reply.Keys))
	assert.Equal(t, 1, len(reply.Values))
	assert.Equal(t, int64(1), reply.Num)
	assert.Equal(t, 0, len(reply.NextKey))

	reply.NextKey = nil
	key3 := "mavl-coins-bty-exec-26htvcBNSEA7fZhAdLJphDwQRQJaHpyHTp:1JmFaA6unrCFYEWPGRi7uuXY1KthTJxJEP"
	bRet = reply.IterateCallBack([]byte(key3), value)
	assert.Equal(t, true, bRet)
	assert.Equal(t, 2, len(reply.Keys))
	assert.Equal(t, 2, len(reply.Values))
	assert.Equal(t, int64(2), reply.Num)
	assert.Equal(t, len(key3), len(reply.NextKey))

	bRet = reply.IterateCallBack([]byte(key), value)
	assert.Equal(t, true, bRet)
	assert.Equal(t, 3, len(reply.Keys))
	assert.Equal(t, 3, len(reply.Values))
	assert.Equal(t, int64(3), reply.Num)
	assert.Equal(t, len(key), len(reply.NextKey))
}

func TestJsonpbUTF8Tx(t *testing.T) {
	NewChain33Config(GetDefaultCfgstring())
	bdata, err := common.FromHex("0a05636f696e73121018010a0c108084af5f1a05310a320a3320e8b31b30b9b69483d7f9d3f04c3a22314b67453376617969715a4b6866684d66744e3776743267447639486f4d6b393431")
	assert.Nil(t, err)
	var r Transaction
	err = Decode(bdata, &r)
	assert.Nil(t, err)
	plType := LoadExecutorType("coins")
	var pl Message
	if plType != nil {
		pl, err = plType.DecodePayload(&r)
		if err != nil {
			pl = nil
		}
	}
	var pljson json.RawMessage
	assert.NotNil(t, pl)
	pljson, err = PBToJSONUTF8(pl)
	assert.Nil(t, err)
	assert.Equal(t, string(pljson), `{"transfer":{"cointoken":"","amount":"200000000","note":"1\n2\n3","to":""},"ty":1}`)
}

func TestSignatureClone(t *testing.T) {
	s1 := &Signature{Ty: 1, Pubkey: []byte("Pubkey1"), Signature: []byte("Signature1")}
	s2 := s1.Clone()
	s2.Pubkey = []byte("Pubkey2")
	assert.Equal(t, s1.Ty, s2.Ty)
	assert.Equal(t, s1.Signature, s2.Signature)
	assert.Equal(t, []byte("Pubkey1"), s1.Pubkey)
	assert.Equal(t, []byte("Pubkey2"), s2.Pubkey)
}

func TestTxClone(t *testing.T) {
	s1 := &Signature{Ty: 1, Pubkey: []byte("Pubkey1"), Signature: []byte("Signature1")}
	tx1 := &Transaction{Execer: []byte("Execer1"), Fee: 1, Signature: s1}
	tx2 := tx1.Clone()
	tx2.Signature.Pubkey = []byte("Pubkey2")
	tx2.Fee = 2
	assert.Equal(t, tx1.Execer, tx2.Execer)
	assert.Equal(t, int64(1), tx1.Fee)
	assert.Equal(t, tx1.Signature.Ty, tx2.Signature.Ty)
	assert.Equal(t, []byte("Pubkey1"), tx1.Signature.Pubkey)
	assert.Equal(t, []byte("Pubkey2"), tx2.Signature.Pubkey)

	tx2.Signature = nil
	assert.NotNil(t, tx1.Signature)
	assert.Nil(t, tx2.Signature)
}

func TestBlockClone(t *testing.T) {
	b1 := getTestBlockDetail()
	b2 := b1.Clone()

	b2.Block.Signature.Ty = 22
	assert.NotEqual(t, b1.Block.Signature.Ty, b2.Block.Signature.Ty)
	assert.Equal(t, b1.Block.Signature.Signature, b2.Block.Signature.Signature)

	b2.Block.Txs[1].Execer = []byte("E22")
	assert.NotEqual(t, b1.Block.Txs[1].Execer, b2.Block.Txs[1].Execer)
	assert.Equal(t, b1.Block.Txs[1].Fee, b2.Block.Txs[1].Fee)

	b2.KV[1].Key = []byte("key22")
	assert.NotEqual(t, b1.KV[1].Key, b2.KV[1].Key)
	assert.Equal(t, b1.KV[1].Value, b2.KV[1].Value)

	b2.Receipts[1].Ty = 22
	assert.NotEqual(t, b1.Receipts[1].Ty, b2.Receipts[1].Ty)
	assert.Equal(t, b1.Receipts[1].Logs, b2.Receipts[1].Logs)

	b2.Block.Txs[0] = nil
	assert.NotNil(t, b1.Block.Txs[0])
}

func TestBlockBody(t *testing.T) {
	detail := getTestBlockDetail()
	b1 := BlockBody{
		Txs:        detail.Block.Txs,
		Receipts:   detail.Receipts,
		MainHash:   []byte("MainHash1"),
		MainHeight: 1,
		Hash:       []byte("Hash"),
		Height:     1,
	}
	b2 := b1.Clone()

	b2.Txs[1].Execer = []byte("E22")
	assert.NotEqual(t, b1.Txs[1].Execer, b2.Txs[1].Execer)
	assert.Equal(t, b1.Txs[1].Fee, b2.Txs[1].Fee)

	b2.Receipts[1].Ty = 22
	assert.NotEqual(t, b1.Receipts[1].Ty, b2.Receipts[1].Ty)
	assert.Equal(t, b1.Receipts[1].Logs, b2.Receipts[1].Logs)

	b2.Txs[0] = nil
	assert.NotNil(t, b1.Txs[0])

	b2.MainHash = []byte("MainHash2")
	assert.NotEqual(t, b1.MainHash, b2.MainHash)
	assert.Equal(t, b1.Height, b2.Height)
}

func getTestBlockDetail() *BlockDetail {
	s1 := &Signature{Ty: 1, Pubkey: []byte("Pubkey1"), Signature: []byte("Signature1")}
	s2 := &Signature{Ty: 2, Pubkey: []byte("Pubkey2"), Signature: []byte("Signature2")}
	tx1 := &Transaction{Execer: []byte("Execer1"), Fee: 1, Signature: s1}
	tx2 := &Transaction{Execer: []byte("Execer2"), Fee: 2, Signature: s2}

	sigBlock := &Signature{Ty: 1, Pubkey: []byte("BlockPubkey1"), Signature: []byte("BlockSignature1")}
	block := &Block{
		Version:    1,
		ParentHash: []byte("ParentHash"),
		TxHash:     []byte("TxHash"),
		StateHash:  []byte("TxHash"),
		Height:     1,
		BlockTime:  1,
		Difficulty: 1,
		MainHash:   []byte("MainHash"),
		MainHeight: 1,
		Signature:  sigBlock,
		Txs:        []*Transaction{tx1, tx2},
	}
	kv1 := &KeyValue{Key: []byte("key1"), Value: []byte("value1")}
	kv2 := &KeyValue{Key: []byte("key1"), Value: []byte("value1")}

	log1 := &ReceiptLog{Ty: 1, Log: []byte("log1")}
	log2 := &ReceiptLog{Ty: 2, Log: []byte("log2")}
	receipts := []*ReceiptData{
		{Ty: 11, Logs: []*ReceiptLog{log1}},
		{Ty: 12, Logs: []*ReceiptLog{log2}},
	}

	return &BlockDetail{
		Block:          block,
		Receipts:       receipts,
		KV:             []*KeyValue{kv1, kv2},
		PrevStatusHash: []byte("PrevStatusHash"),
	}
}
