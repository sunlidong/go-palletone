package jury

import (
	"time"
	"encoding/hex"
	"fmt"

	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/dag/errors"
	"github.com/palletone/go-palletone/core/vmContractPub/crypto"
)

func (p *Processor) ContractInstallReq(from, to common.Address, daoAmount, daoFee uint64, tplName, path, version string, local bool) (reqId []byte, TplId []byte, err error) {
	if from == (common.Address{}) || to == (common.Address{}) || tplName == "" || path == "" || version == "" {
		log.Error("ContractInstallReq", "param is error")
		return nil, nil, errors.New("ContractInstallReq request param is error")
	}

	log.Debug("ContractInstallReq", "enter, tplName ", tplName, "path", path, "version", version)
	msgReq := &modules.Message{
		App: modules.APP_CONTRACT_TPL_REQUEST,
		Payload: &modules.ContractInstallRequestPayload{
			TplName: tplName,
			Path:    path,
			Version: version,
		},
	}
	reqId, tx, err := p.createContractTxReq(from, to, daoAmount, daoFee, msgReq, true)
	if err != nil {
		return nil, nil, err
	}
	tpl, err := getContractTxContractInfo(tx, modules.APP_CONTRACT_TPL)
	if err != nil {
		errMsg := fmt.Sprintf("getContractTxContractInfo fail, tpl Name[%s]", tplName)
		return nil, nil, errors.New(errMsg)
	}
	templateId := tpl.(*modules.ContractTplPayload).TemplateId

	//broadcast
	go p.ptn.ContractBroadcast(ContractEvent{CType: CONTRACT_EVENT_COMMIT, Tx: tx})
	return reqId, templateId, nil
}

func (p *Processor) ContractDeployReq(from, to common.Address, daoAmount, daoFee uint64, templateId []byte, txid string, args [][]byte, timeout time.Duration) ([]byte, error) {
	if from == (common.Address{}) || to == (common.Address{}) || templateId == nil {
		log.Error("ContractDeployReq", "param is error")
		return nil, errors.New("ContractDeployReq request param is error")
	}
	randNum, err := crypto.GetRandomNonce()
	if err != nil {
		return nil, errors.New("GetRandomNonce error")
	}
	log.Debug("ContractDeployReq", "enter, templateId ", templateId, "txId", hex.EncodeToString(randNum))
	msgReq := &modules.Message{
		App: modules.APP_CONTRACT_DEPLOY_REQUEST,
		Payload: &modules.ContractDeployRequestPayload{
			TplId:   templateId,
			TxId:    hex.EncodeToString(randNum),
			Args:    args,
			Timeout: timeout,
		},
	}
	reqId, tx, err := p.createContractTxReq(from, to, daoAmount, daoFee, msgReq, false)
	if err != nil {
		return nil, err
	}
	//broadcast
	go p.ptn.ContractBroadcast(ContractEvent{CType: CONTRACT_EVENT_EXEC, Tx: tx})
	return reqId, nil
}

func (p *Processor) ContractInvokeReq(from, to common.Address, daoAmount, daoFee uint64, contractId common.Address, args [][]byte, timeout time.Duration) ([]byte, error) {
	if from == (common.Address{}) || to == (common.Address{}) || contractId == (common.Address{}) || args == nil {
		log.Error("ContractInvokeReq", "param is error")
		return nil, errors.New("ContractInvokeReq request param is error")
	}

	log.Debug("ContractInvokeReq", "enter, contractId ", contractId)
	msgReq := &modules.Message{
		App: modules.APP_CONTRACT_INVOKE_REQUEST,
		Payload: &modules.ContractInvokeRequestPayload{
			ContractId:   contractId.Bytes(),
			FunctionName: "",
			Args:         args,
			Timeout:      timeout,
		},
	}
	reqId, tx, err := p.createContractTxReq(from, to, daoAmount, daoFee, msgReq, false)
	if err != nil {
		return nil, err
	}
	//broadcast
	go p.ptn.ContractBroadcast(ContractEvent{CType: CONTRACT_EVENT_EXEC, Tx: tx})
	return reqId, nil
}

func (p *Processor) ContractStopReq(from, to common.Address, daoAmount, daoFee uint64, contractId common.Address, txid string, deleteImage bool) ([]byte, error) {
	if from == (common.Address{}) || to == (common.Address{}) || contractId == (common.Address{}) {
		log.Error("ContractStopReq", "param is error")
		return nil, errors.New("ContractStopReq request param is error")
	}

	log.Debug("ContractStopReq", "enter, contractId ", contractId)
	msgReq := &modules.Message{
		App: modules.APP_CONTRACT_STOP_REQUEST,
		Payload: &modules.ContractStopRequestPayload{
			ContractId:  contractId[:],
			Txid:        txid,
			DeleteImage: deleteImage,
		},
	}
	reqId, tx, err := p.createContractTxReq(from, to, daoAmount, daoFee, msgReq, false)
	if err != nil {
		return nil, err
	}
	//broadcast
	go p.ptn.ContractBroadcast(ContractEvent{CType: CONTRACT_EVENT_EXEC, Tx: tx})
	return reqId, nil
}
