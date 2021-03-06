/* This file is part of go-palletone.
   go-palletone is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-palletone is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.

   @author PalletOne core developers <dev@pallet.one>
   @date 2018
*/

package modules

import (
	"fmt"
	"strings"
	"time"
	"reflect"

	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/common/rlp"

	"bytes"
)

type MessageType byte

const (
	APP_PAYMENT MessageType = iota

	APP_CONTRACT_TPL
	APP_CONTRACT_DEPLOY
	APP_CONTRACT_INVOKE
	APP_CONTRACT_STOP
	APP_SIGNATURE

	APP_CONFIG
	APP_DATA
	APP_VOTE
	OP_MEDIATOR_CREATE

	APP_UNKNOW

	APP_CONTRACT_TPL_REQUEST    = 100
	APP_CONTRACT_DEPLOY_REQUEST = 101
	APP_CONTRACT_INVOKE_REQUEST = 102
	APP_CONTRACT_STOP_REQUEST   = 103
)

// key: message.UnitHash(message+timestamp)
type Message struct {
	App     MessageType `json:"app"`     // message type
	Payload interface{} `json:"payload"` // the true transaction data
}

// return message struct
func NewMessage(app MessageType, payload interface{}) *Message {
	m := new(Message)
	m.App = app
	m.Payload = payload
	return m
}

func (msg *Message) CopyMessages(cpyMsg *Message) *Message {
	msg.App = cpyMsg.App
	//msg.Payload = cpyMsg.Payload
	switch cpyMsg.App {
	// modified by albert·gou
	default:
		//case APP_PAYMENT, APP_CONTRACT_TPL, APP_DATA, APP_VOTE:
		msg.Payload = cpyMsg.Payload
	case APP_CONFIG:
		payload, _ := cpyMsg.Payload.(*ConfigPayload)
		newPayload := ConfigPayload{
			ConfigSet: []ContractWriteSet{},
		}
		for _, p := range payload.ConfigSet {
			newPayload.ConfigSet = append(newPayload.ConfigSet, ContractWriteSet{Key: p.Key, Value: p.Value})
		}
		msg.Payload = newPayload
	case APP_CONTRACT_DEPLOY:
		payload, _ := cpyMsg.Payload.(*ContractDeployPayload)
		newPayload := ContractDeployPayload{
			TemplateId: payload.TemplateId,
			ContractId: payload.ContractId,
			Args:       payload.Args,
			//ExecutionTime: payload.ExecutionTime,
		}
		readSet := []ContractReadSet{}
		for _, rs := range payload.ReadSet {
			readSet = append(readSet, ContractReadSet{Key: rs.Key, Version: &StateVersion{Height: rs.Version.Height, TxIndex: rs.Version.TxIndex}})
		}
		writeSet := []ContractWriteSet{}
		for _, ws := range payload.WriteSet {
			writeSet = append(writeSet, ContractWriteSet{Key: ws.Key, Value: ws.Value})
		}
		newPayload.ReadSet = readSet
		newPayload.WriteSet = writeSet
		msg.Payload = newPayload
	case APP_CONTRACT_INVOKE_REQUEST:
		payload, _ := cpyMsg.Payload.(*ContractInvokeRequestPayload)
		newPayload := ContractInvokeRequestPayload{
			ContractId: payload.ContractId,
			Args:       payload.Args,
			Timeout:    payload.Timeout,
		}
		msg.Payload = newPayload
	case APP_CONTRACT_INVOKE:
		payload, _ := cpyMsg.Payload.(*ContractInvokePayload)
		newPayload := ContractInvokePayload{
			ContractId: payload.ContractId,
			Args:       payload.Args,
			//ExecutionTime: payload.ExecutionTime,
		}
		readSet := []ContractReadSet{}
		for _, rs := range payload.ReadSet {
			readSet = append(readSet, ContractReadSet{Key: rs.Key, Version: &StateVersion{Height: rs.Version.Height, TxIndex: rs.Version.TxIndex}})
		}
		writeSet := []ContractWriteSet{}
		for _, ws := range payload.WriteSet {
			writeSet = append(writeSet, ContractWriteSet{Key: ws.Key, Value: ws.Value})
		}
		newPayload.ReadSet = readSet
		newPayload.WriteSet = writeSet
		msg.Payload = newPayload
	case APP_SIGNATURE:
		payload, _ := cpyMsg.Payload.(*SignaturePayload)
		newPayload := SignaturePayload{}
		newPayload.Signatures = payload.Signatures
		msg.Payload = newPayload
	}

	return msg
}

func (msg *Message) CompareMessages(inMsg *Message) bool {
	//return true //todo del

	if inMsg == nil || msg.App != inMsg.App {
		return false
	}
	switch msg.App {
	case APP_CONTRACT_TPL:
		payA, _ := msg.Payload.(*ContractTplPayload)
		payB, _ := inMsg.Payload.(*ContractTplPayload)
		return payA.Equal(payB)
	case APP_CONTRACT_DEPLOY:
		payA, _ := msg.Payload.(*ContractDeployPayload)
		payB, _ := inMsg.Payload.(*ContractDeployPayload)
		return payA.Equal(payB)
	case APP_CONTRACT_INVOKE:
		payA, _ := msg.Payload.(*ContractInvokePayload)
		payB, _ := inMsg.Payload.(*ContractInvokePayload)
		return payA.Equal(payB)
	case APP_CONTRACT_STOP:
		payA, _ := msg.Payload.(*ContractStopPayload)
		payB, _ := inMsg.Payload.(*ContractStopPayload)
		return payA.Equal(payB)
	case APP_SIGNATURE:
		//todo
		//payA, _ := msg.Payload.(*SignaturePayload)
		//payB, _ := inMsg.Payload.(*SignaturePayload)
		return true
	case APP_CONTRACT_TPL_REQUEST:
		payA, _ := msg.Payload.(*ContractInstallRequestPayload)
		payB, _ := inMsg.Payload.(*ContractInstallRequestPayload)
		return payA.Equal(payB)
	case APP_CONTRACT_DEPLOY_REQUEST:
		payA, _ := msg.Payload.(*ContractDeployRequestPayload)
		payB, _ := inMsg.Payload.(*ContractDeployRequestPayload)
		return reflect.DeepEqual(payA, payB)
	case APP_CONTRACT_INVOKE_REQUEST:
		payA, _ := msg.Payload.(*ContractInvokeRequestPayload)
		payB, _ := inMsg.Payload.(*ContractInvokeRequestPayload)
		return payA.Equal(payB)
	case APP_CONTRACT_STOP_REQUEST:
		payA, _ := msg.Payload.(*ContractStopRequestPayload)
		payB, _ := inMsg.Payload.(*ContractStopRequestPayload)
		return payA.Equal(payB)
	default:
		return false
	}

	return false
}

type ContractWriteSet struct {
	IsDelete bool
	Key      string
	Value    []byte
	//Value interface{}
}

func ToPayloadMapValueBytes(data interface{}) []byte {
	b, err := rlp.EncodeToBytes(data)
	if err != nil {
		return nil
	}
	return b
}

// Token exchange message and verify message
// App: payment
type PaymentPayload struct {
	Inputs   []*Input  `json:"inputs"`
	Outputs  []*Output `json:"outputs"`
	LockTime uint32    `json:"lock_time"`
}

func (pay *PaymentPayload) IsCoinbase() bool {
	if len(pay.Inputs) == 0 {
		return true
	}
	if pay.Inputs[0].PreviousOutPoint == nil {
		return true
	}
	return false
}

// NewTxOut returns a new bitcoin transaction output with the provided
// transaction value and public key script.
func NewTxOut(value uint64, pkScript []byte, asset *Asset) *Output {
	return &Output{
		Value:    value,
		PkScript: pkScript,
		Asset:    asset,
	}
}

type StateVersion struct {
	Height  ChainIndex `json:"height"`
	TxIndex uint32     `json:"tx_index"`
}
type ContractStateValue struct {
	Value   []byte        `json:"value"`
	Version *StateVersion `json:"version"`
}

func (version *StateVersion) String() string {

	return fmt.Sprintf(
		"StateVersion[AssetId:{%#x}, Height:{%d},IsMain:%t,TxIdx:{%d}]",
		version.Height.AssetID,
		version.Height.Index,
		version.Height.IsMain,
		version.TxIndex)
}

func (version *StateVersion) ParseStringKey(key string) bool {
	ss := strings.Split(key, FIELD_SPLIT_STR)
	if len(ss) != 3 {
		return false
	}
	var v StateVersion
	if err := rlp.DecodeBytes([]byte(ss[2]), &v); err != nil {
		log.Error("State version parse string key", "error", err.Error())
		return false
	}
	if version == nil {
		version = &StateVersion{}
	}
	version.Height = v.Height
	version.TxIndex = v.TxIndex
	return true
}

//16+8+1+4=29
func (version *StateVersion) Bytes() []byte {
	idx := make([]byte, 8)
	littleEndian.PutUint64(idx, version.Height.Index)
	b := append(version.Height.AssetID.Bytes(), idx...)
	if version.Height.IsMain {
		b = append(b, byte(1))
	} else {
		b = append(b, byte(0))
	}
	txIdx := make([]byte, 4)
	littleEndian.PutUint32(txIdx, version.TxIndex)
	b = append(b, txIdx...)
	return b[:]
}
func (version *StateVersion) SetBytes(b []byte) {
	asset := IDType16{}
	asset.SetBytes(b[:15])
	heightIdx := littleEndian.Uint64(b[16:24])
	isMain := b[24]
	txIdx := littleEndian.Uint32(b[25:])
	cidx := ChainIndex{AssetID: asset, Index: heightIdx, IsMain: isMain == byte(1)}
	version.Height = cidx
	version.TxIndex = txIdx
}

const (
	FIELD_TPL_BYTECODE  = "TplBytecode"
	FIELD_TPL_NAME      = "TplName"
	FIELD_TPL_PATH      = "TplPath"
	FIELD_TPL_Memory    = "TplMemory"
	FIELD_SPLIT_STR     = "^*^"
	FIELD_GENESIS_ASSET = "GenesisAsset"
	FIELD_TPL_Version   = "TplVersion"
)

type DelContractState struct {
	IsDelete bool
}

func (delState DelContractState) Bytes() []byte {
	data, err := rlp.EncodeToBytes(delState)
	if err != nil {
		return nil
	}
	return data
}

func (delState DelContractState) SetBytes(b []byte) error {
	if err := rlp.DecodeBytes(b, &delState); err != nil {
		return err
	}
	return nil
}

type ContractReadSet struct {
	Key     string
	Version *StateVersion
	Value   []byte
}

//请求合约信息
type InvokeInfo struct {
	InvokeAddress string        `json:"invoke_address"` //请求地址
	InvokeTokens  *InvokeTokens `json:"invoke_tokens"`  //请求数量
	InvokeFees    *InvokeFees   `json:"invoke_fees"`    //请求交易费
}

//请求的数量
type InvokeTokens struct {
	Amount uint64 `json:"amount"` //数量
	Asset  *Asset `json:"asset"`  //资产
}

//申请提保证金
type Cashback struct {
	CashbackAddress string        `json:"cashback_address"` //请求地址
	CashbackTokens  *InvokeTokens `json:"cashback_tokens"`  //请求数量
	Role            string        `json:"role"`             //请求角色
	CashbackTime    int64         `json:"cashback_time"`    //请求时间
}

//申请提取保证金的列表
//type ListForCashback struct {
//	Cashbacks []*Cashback `json:"cashbacks"`
//}

//申请没收保证金
type Forfeiture struct {
	ApplyAddress      string        `json:"apply_address"`      //谁发起的
	ForfeitureAddress string        `json:"forfeiture_address"` //没收节点地址
	ApplyTokens       *InvokeTokens `json:"apply_tokens"`       //没收数量
	ForfeitureRole    string        `json:"forfeiture_role"`    //没收角色
	//Extra             string        `json:"extra"`              //备注
	ApplyTime int64 `json:"apply_time"` //请求时间
}

//申请没收保证金的列表
//type ListForForfeiture struct {
//	Forfeitures []*Forfeiture `json:"forfeitures"`
//}

//请求合约利息
type InvokeFees struct {
	Amount uint64 `json:"amount"`
	Asset  *Asset `json:"asset"`
}

//申请成为Mediator
type MediatorRegisterInfo struct {
	Address string `json:"address"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

//Mediator 申请列表
//type BecomeMediatorApplyList struct {
//	MediatorList []*MediatorInfo `json:"mediator_list"`
//}

//同意 Mediator 列表
//type AgreeForBecomeMediatorList struct {
//	MediatorList []*MediatorInfo `json:"mediator_list"`
//}

//Mediator 退出列表
//type QuitMediatorApplyList struct {
//	MediatorList []*MediatorInfo `json:"mediator_list"`
//}

//数量及资产类型
//type AmountAsset struct {
//	Amount uint64 `json:"amount"`
//	Asset  *Asset `json:"asset"`
//}

//节点状态数据库保存值
type DepositBalance struct {
	TotalAmount      uint64        `json:"total_amount"`      //保证金总量
	LastModifyTime   int64         `json:"last_modify_time"`  //最后一次改变，主要来计算币龄收益
	EnterTime        int64         `json:"enter_time"`        //这是加入列表时的时间
	PayValues        []*PayValue   `json:"pay_values"`        //交付的历史记录
	CashbackValues   []*Cashback   `json:"cashback_values"`   //退款的历史记录
	ForfeitureValues []*Forfeiture `json:"forfeiture_values"` //被没收的历史记录
}

//交易的内容
type PayValue struct {
	PayTokens *InvokeTokens `json:"pay_tokens"` //数量和资产
	PayTime   int64         `json:"pay_time"`   //发生时间
	//PayExtra  string        `json:"pay_extra"`  //额外内容
}

type TokenPayOut struct {
	Asset    *Asset
	Amount   uint64
	PayTo    common.Address
	LockTime uint32
}

// Contract template deploy message
// App: contract_template
type ContractTplPayload struct {
	TemplateId []byte `json:"template_id"` // contract template id
	Name       string `json:"name"`        // contract template name
	Path       string `json:"path"`        // contract template execute path
	Version    string `json:"version"`     // contract template version
	Memory     uint16 `json:"memory"`      // contract template bytecode memory size(Byte), use to compute transaction fee
	Bytecode   []byte `json:"bytecode"`    // contract bytecode
}

// App: contract_deploy
type ContractDeployPayload struct {
	TemplateId []byte             `json:"template_id"` // contract template id
	ContractId []byte             `json:"contract_id"` // contract id
	Name       string             `json:"name"`        // the name for contract
	Args       [][]byte           `json:"args"`        // contract arguments list
	Jury       []common.Address   `json:"jury"`        // contract jurors list
	ReadSet    []ContractReadSet  `json:"read_set"`    // the set data of read, and value could be any type
	WriteSet   []ContractWriteSet `json:"write_set"`   // the set data of write, and value could be any type
}

// Contract invoke message
// App: contract_invoke
//如果是用户想修改自己的State信息，那么ContractId可以为空或者0字节
type ContractInvokePayload struct {
	ContractId   []byte             `json:"contract_id"` // contract id
	FunctionName string             `json:"function_name"`
	Args         [][]byte           `json:"args"`      // contract arguments list
	ReadSet      []ContractReadSet  `json:"read_set"`  // the set data of read, and value could be any type
	WriteSet     []ContractWriteSet `json:"write_set"` // the set data of write, and value could be any type
	Payload      []byte             `json:"payload"`   // the contract execution result
}

// App: contract_deploy
type ContractStopPayload struct {
	ContractId []byte             `json:"contract_id"` // contract id
	Jury       []common.Address   `json:"jury"`        // contract jurors list
	ReadSet    []ContractReadSet  `json:"read_set"`    // the set data of read, and value could be any type
	WriteSet   []ContractWriteSet `json:"write_set"`   // the set data of write, and value could be any type
}

//contract invoke result
type ContractInvokeResult struct {
	ContractId   []byte             `json:"contract_id"` // contract id
	RequestId    common.Hash        `json:"request_id"`
	FunctionName string             `json:"function_name"`
	Args         [][]byte           `json:"args"`         // contract arguments list
	ReadSet      []ContractReadSet  `json:"read_set"`     // the set data of read, and value could be any type
	WriteSet     []ContractWriteSet `json:"write_set"`    // the set data of write, and value could be any type
	Payload      []byte             `json:"payload"`      // the contract execution result
	TokenPayOut  []*TokenPayOut     `json:"token_payout"` //从合约地址付出Token
	TokenSupply  []*TokenSupply     `json:"token_supply"` //增发Token请求产生的结果
	TokenDefine  *TokenDefine       `json:"token_define"` //定义新Token
}

//用户钱包发起的合约调用申请
type ContractInstallRequestPayload struct {
	TplName string `json:"tpl_name"`
	Path    string `json:"install_path"`
	Version string `json:"tpl_version"`
}

type ContractDeployRequestPayload struct {
	TplId   []byte        `json:"tpl_name"`
	TxId    string        `json:"transaction_id"` //todo
	Args    [][]byte      `json:"args"`
	Timeout time.Duration `json:"timeout"`
}

type ContractInvokeRequestPayload struct {
	ContractId   []byte        `json:"contract_id"` // contract id
	FunctionName string        `json:"function_name"`
	Args         [][]byte      `json:"args"` // contract arguments list
	Timeout      time.Duration `json:"timeout"`
}

type ContractStopRequestPayload struct {
	ContractId  []byte `json:"contract_id"`
	Txid        string `json:"transaction_id"`
	DeleteImage bool   `json:"delete_image"`
}

// Token exchange message and verify message
// App: config	// update global config
type ConfigPayload struct {
	ConfigSet []ContractWriteSet `json:"config_set"` // the array of global config
}
type SignaturePayload struct {
	Signatures []SignatureSet `json:"signature_set"` // the array of signature
}
type SignatureSet struct {
	PubKey    []byte //compress public key
	Signature []byte //
}

// Token exchange message and verify message
// App: text
type DataPayload struct {
	MainData  []byte `json:"main_data"`
	ExtraData []byte `json:"extra_data"`
}
type FileInfo struct {
	UnitHash    common.Hash   `json:"unit_hash"`
	UintHeight  uint64        `json:"unit_index"`
	ParentsHash []common.Hash   `json:"parents_hash"`
	Txid        common.Hash   `json:"txid"`
	Timestamp   int64 `json:"timestamp"`
	MainData  []byte `json:"main_data"`
	ExtraData []byte `json:"extra_data"`
}

func NewPaymentPayload(inputs []*Input, outputs []*Output) *PaymentPayload {
	return &PaymentPayload{
		Inputs:   inputs,
		Outputs:  outputs,
		LockTime: defaultTxInOutAlloc,
	}
}

func NewContractTplPayload(templateId []byte, name string, path string, version string, memory uint16, bytecode []byte) *ContractTplPayload {
	return &ContractTplPayload{
		TemplateId: templateId,
		Name:       name,
		Path:       path,
		Version:    version,
		Memory:     memory,
		Bytecode:   bytecode,
	}
}

func NewContractDeployPayload(templateid []byte, contractid []byte, name string, args [][]byte, excutiontime time.Duration,
	jury []common.Address, readset []ContractReadSet, writeset []ContractWriteSet) *ContractDeployPayload {
	return &ContractDeployPayload{
		TemplateId: templateid,
		ContractId: contractid,
		Name:       name,
		Args:       args,
		//ExecutionTime: excutiontime,
		Jury:     jury,
		ReadSet:  readset,
		WriteSet: writeset,
	}
}

//TokenPayOut   []*modules.TokenPayOut     `json:"token_payout"`   //从合约地址付出Token
//	TokenSupply   []*modules.TokenSupply     `json:"token_supply"`   //增发Token请求产生的结果
//	TokenDefine   *modules.TokenDefine       `json:"token_define"`   //定义新Token
func NewContractInvokePayload(contractid []byte, funcName string, args [][]byte, excutiontime time.Duration,
	readset []ContractReadSet, writeset []ContractWriteSet, payload []byte) *ContractInvokePayload {
	return &ContractInvokePayload{
		ContractId:   contractid,
		FunctionName: funcName,
		Args:         args,
		//ExecutionTime: excutiontime,
		ReadSet:  readset,
		WriteSet: writeset,
		Payload:  payload,
		//TokenPayOut:   tokenPayOut,
		//TokenSupply:   tokenSupply,
		//TokenDefine:   tokenDefine,
	}
}

func (a *ContractReadSet) Equal(b *ContractReadSet) bool {
	if b == nil {
		return false
	}
	if !strings.EqualFold(a.Key, b.Key) || !bytes.Equal(a.Value, b.Value) {
		return false
	}
	if a.Version != nil && b.Version != nil {
		if a.Version.TxIndex != b.Version.TxIndex || a.Version.Height != b.Version.Height {
			return false
		}
	} else if a.Version != b.Version {
		return false
	}

	return true
}

func (a *ContractWriteSet) Equal(b *ContractWriteSet) bool {
	if b == nil {
		return false
	}
	if !(a.IsDelete == b.IsDelete) || !strings.EqualFold(a.Key, b.Key) || !bytes.Equal(a.Value, b.Value) {
		return false
	}
	return true
}

func (a *ContractTplPayload) Equal(b *ContractTplPayload) bool {
	if b == nil {
		return false
	}
	if bytes.Equal(a.TemplateId, b.TemplateId) && strings.EqualFold(a.Name, b.Name) && strings.EqualFold(a.Path, b.Path) &&
		strings.EqualFold(a.Version, b.Version) && a.Memory == b.Memory && bytes.Equal(a.Bytecode, b.Bytecode) {
		return true
	}
	return false
}

func (a *ContractDeployPayload) Equal(b *ContractDeployPayload) bool {
	if b == nil {
		return false
	}
	if !bytes.Equal(a.TemplateId, b.TemplateId) || !bytes.Equal(a.ContractId, b.ContractId) || !strings.EqualFold(a.Name, b.Name) {
		return false
	}
	if len(a.Args) == len(b.Args) {
		for i := 0; i < len(a.Args); i++ {
			if !bytes.Equal(a.Args[i], b.Args[i]) {
				return false
			}
		}
	} else {
		return false
	}
	if len(a.Jury) == len(b.Jury) {
		for i := 0; i < len(a.Jury); i++ {
			if !a.Jury[i].Equal(b.Jury[i]) {
				return false
			}
		}
	} else {
		return false
	}
	if len(a.ReadSet) == len(b.ReadSet) {
		for i := 0; i < len(a.ReadSet); i++ {
			a.ReadSet[i].Equal(&b.ReadSet[i])
		}
	} else {
		return false
	}
	if len(a.WriteSet) == len(b.WriteSet) {
		for i := 0; i < len(a.WriteSet); i++ {
			a.WriteSet[i].Equal(&b.WriteSet[i])
		}
	} else {
		return false
	}
	return true
}

func (a *ContractInvokePayload) Equal(b *ContractInvokePayload) bool {
	if b == nil {
		return false
	}
	if !bytes.Equal(a.ContractId, b.ContractId) || !strings.EqualFold(a.FunctionName, b.FunctionName) || !bytes.Equal(a.Payload, b.Payload) {
		return false
	}
	if len(a.Args) == len(b.Args) {
		for i := 0; i < len(a.Args); i++ {
			if !bytes.Equal(a.Args[i], b.Args[i]) {
				return false
			}
		}
	} else {
		return false
	}
	if len(a.ReadSet) == len(b.ReadSet) {
		for i := 0; i < len(a.ReadSet); i++ {
			a.ReadSet[i].Equal(&b.ReadSet[i])
		}
	} else {
		return false
	}
	if len(a.WriteSet) == len(b.WriteSet) {
		for i := 0; i < len(a.WriteSet); i++ {
			a.WriteSet[i].Equal(&b.WriteSet[i])
		}
	} else {
		return false
	}
	return true
}

func (a *ContractStopPayload) Equal(b *ContractStopPayload) bool {
	if b == nil {
		return false
	}
	if !bytes.Equal(a.ContractId, b.ContractId) {
		return false
	}
	if len(a.Jury) == len(b.Jury) {
		for i := 0; i < len(a.Jury); i++ {
			if !a.Jury[i].Equal(b.Jury[i]) {
				return false
			}
		}
	} else {
		return false
	}
	if len(a.ReadSet) == len(b.ReadSet) {
		for i := 0; i < len(a.ReadSet); i++ {
			a.ReadSet[i].Equal(&b.ReadSet[i])
		}
	} else {
		return false
	}
	if len(a.WriteSet) == len(b.WriteSet) {
		for i := 0; i < len(a.WriteSet); i++ {
			a.WriteSet[i].Equal(&b.WriteSet[i])
		}
	} else {
		return false
	}
	return true
}

func (a *ContractInstallRequestPayload) Equal(b *ContractInstallRequestPayload) bool {
	if b == nil {
		return false
	}
	if !strings.EqualFold(a.TplName, b.TplName) || !strings.EqualFold(a.Path, b.Path) || !strings.EqualFold(a.Version, b.Version) {
		return false
	}
	return true
}

func (a *ContractDeployRequestPayload) Equal(b *ContractDeployRequestPayload) bool {
	if b == nil {
		return false
	}
	if !bytes.Equal(a.TplId, b.TplId) || !strings.EqualFold(a.TxId, b.TxId) || a.Timeout != b.Timeout {
		return false
	}
	if len(a.Args) == len(b.Args) {
		for i := 0; i < len(a.Args); i++ {
			if !bytes.Equal(a.Args[i], b.Args[i]) {
				return false
			}
		}
	} else {
		return false
	}
	return true
}

func (a *ContractInvokeRequestPayload) Equal(b *ContractInvokeRequestPayload) bool {
	if b == nil {
		return false
	}
	if !bytes.Equal(a.ContractId, b.ContractId) || !strings.EqualFold(a.FunctionName, b.FunctionName) || a.Timeout != b.Timeout {
		return false
	}
	if len(a.Args) == len(b.Args) {
		for i := 0; i < len(a.Args); i++ {
			if !bytes.Equal(a.Args[i], b.Args[i]) {
				return false
			}
		}
	} else {
		return false
	}
	return true
}

func (a *ContractStopRequestPayload) Equal(b *ContractStopRequestPayload) bool {
	if b == nil {
		return false
	}
	if !bytes.Equal(a.ContractId, b.ContractId) || !strings.EqualFold(a.Txid, b.Txid) || a.DeleteImage != b.DeleteImage {
		return false
	}
	return true
}
