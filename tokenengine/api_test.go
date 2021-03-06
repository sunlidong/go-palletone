/*
 *
 *    This file is part of go-palletone.
 *    go-palletone is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU General Public License as published by
 *    the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *    go-palletone is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU General Public License for more details.
 *    You should have received a copy of the GNU General Public License
 *    along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.
 * /
 *
 *  * @author PalletOne core developer <dev@pallet.one>
 *  * @date 2018
 *
 */

package tokenengine

import (
	"testing"

	//"github.com/palletone/go-palletone/tokenengine/btcd/chaincfg/chainhash"
	"encoding/hex"
	"fmt"
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/crypto"

	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/tokenengine/internal/txscript"
	"github.com/stretchr/testify/assert"
)

func TestGetAddressFromScript(t *testing.T) {
	addrStr := "P1JEStL6tb7TB8e6ZJSpJhQoqin2A6pabdA"
	addr, _ := common.StringToAddress(addrStr)
	p2pkhLock := GenerateP2PKHLockScript(addr.Bytes())
	t.Logf("P2PKH script:%x", p2pkhLock)
	getAddr, _ := GetAddressFromScript(p2pkhLock)
	t.Logf("Get Address:%s", getAddr.Str())
	assert.True(t, getAddr == addr, "Address parse error")

	addr2, _ := common.StringToAddress("P35SbSqXuXcHrtZuJKzbStpcqzwCg88jXfn")
	p2shLock := GenerateP2SHLockScript(addr2.Bytes())
	getAddr2, _ := GetAddressFromScript(p2shLock)
	t.Logf("Get Script Address:%s", getAddr2.Str())
	assert.True(t, getAddr2 == addr2, "Address parse error")

}

func TestGenerateP2CHLockScript(t *testing.T) {
	addrStr := "PCGTta3M4t3yXu8uRgkKvaWd2d8DR32W9vM"
	addr, err := common.StringToAddress(addrStr)
	assert.Nil(t, err)
	p2ch1 := GenerateLockScript(addr)
	p2ch1Str, _ := DisasmString(p2ch1)
	t.Logf("Pay to contract hash lock script:%x, string:%s", p2ch1, p2ch1Str)
	p2ch2 := GenerateP2CHLockScript(addr)
	assert.Equal(t, p2ch1, p2ch2)
	addr2, err := GetAddressFromScript(p2ch1)
	assert.Nil(t, err, "Err must nil")
	assert.Equal(t, addr2.String(), addrStr)
	t.Logf("get address:%s", addr2.String())
}

func TestSignAndVerifyATx(t *testing.T) {

	privKeyBytes, _ := hex.DecodeString("2BE3B4B671FF5B8009E6876CCCC8808676C1C279EE824D0AB530294838DC1644")
	privKey, _ := crypto.ToECDSA(privKeyBytes)
	pubKey := privKey.PublicKey
	pubKeyBytes := crypto.CompressPubkey(&pubKey)
	pubKeyHash := crypto.Hash160(pubKeyBytes)
	t.Logf("Public Key:%x", pubKeyBytes)
	addr := crypto.PubkeyToAddress(&privKey.PublicKey)
	t.Logf("Addr:%s", addr.String())
	lockScript := GenerateP2PKHLockScript(pubKeyHash)
	t.Logf("UTXO lock script:%x", lockScript)

	tx := &modules.Transaction{
		TxMessages: make([]*modules.Message, 0),
	}
	payment := &modules.PaymentPayload{}
	utxoTxId, _ := common.NewHashFromStr("5651870aa8c894376dbd960a22171d0ad7be057a730e14d7103ed4a6dbb34873")
	outPoint := modules.NewOutPoint(utxoTxId, 0, 0)
	txIn := modules.NewTxIn(outPoint, []byte{})
	payment.AddTxIn(txIn)
	asset0 := &modules.Asset{}
	payment.AddTxOut(modules.NewTxOut(1, lockScript, asset0))
	payment2 := &modules.PaymentPayload{}
	utxoTxId2, _ := common.NewHashFromStr("1651870aa8c894376dbd960a22171d0ad7be057a730e14d7103ed4a6dbb34873")
	outPoint2 := modules.NewOutPoint(utxoTxId2, 1, 1)
	txIn2 := modules.NewTxIn(outPoint2, []byte{})
	payment2.AddTxIn(txIn2)
	asset1 := &modules.Asset{AssetId: modules.PTNCOIN}
	payment2.AddTxOut(modules.NewTxOut(1, lockScript, asset1))
	tx.TxMessages = append(tx.TxMessages, modules.NewMessage(modules.APP_PAYMENT, payment))
	tx.TxMessages = append(tx.TxMessages, modules.NewMessage(modules.APP_PAYMENT, payment2))

	tx.TxMessages = append(tx.TxMessages, modules.NewMessage(modules.APP_DATA, &modules.DataPayload{MainData: []byte("Hello PalletOne")}))

	//signResult, err := SignOnePaymentInput(tx, 0, 0, lockScript, privKey)
	//if err != nil {
	//	t.Errorf("Sign error:%s", err)
	//}
	//t.Logf("Sign Result:%x", signResult)
	//t.Logf("msg len:%d", len(tx.TxMessages))
	//tx.TxMessages[0].Payload.(*modules.PaymentPayload).Inputs[0].SignatureScript = signResult
	//
	//signResult2, err := SignOnePaymentInput(tx, 1, 0, lockScript, privKey)
	//tx.TxMessages[1].Payload.(*modules.PaymentPayload).Input[0].SignatureScript = signResult2
	lockScripts := map[modules.OutPoint][]byte{
		*outPoint:  lockScript[:],
		*outPoint2: GenerateP2PKHLockScript(pubKeyHash),
	}
	//privKeys := map[common.Address]*ecdsa.PrivateKey{
	//	addr: privKey,
	//}
	getPubKeyFn := func(common.Address) ([]byte, error) {
		return crypto.CompressPubkey(&privKey.PublicKey), nil
	}
	getSignFn := func(addr common.Address, hash []byte) ([]byte, error) {
		return crypto.Sign(hash, privKey)
	}
	var hashtype uint32
	hashtype = 1
	_, err := SignTxAllPaymentInput(tx, hashtype, lockScripts, nil, getPubKeyFn, getSignFn, 0)
	if err != nil {
		t.Logf("Sign error:%s", err)
	}
	err = ScriptValidate(lockScript, nil, tx, 0, 0)
	if err != nil {
		t.Logf("validate error:%s", err)
	}
	// textPayload :=tx.TxMessages[2].Payload.(*modules.DataPayload)
	//textPayload.Text=[]byte("Bad")
	//fmt.Printf("%s", tx.TxMessages[2].Payload.(*modules.DataPayload))

	err = ScriptValidate(lockScript, nil, tx, 1, 0)
	assert.Nil(t, err, fmt.Sprintf("validate error:%s", err))

}

var (
	prvKey1, _  = crypto.FromWIF("KwN8TdhAMeU8b9UrEYTNTVEvDsy9CSyepwRVNEy2Fc9nbGqDZw4J") //"0454b0699a590b6fc8e66e81db1ca36e99d7c767cdfe44a217b6e105c5db97d5" //P1QJNzZhqGoxNL2igkdthNBQLNWdNGTWzQU
	pubKey1B, _ = hex.DecodeString("02f9286c44fe7ebff9788425d5025ad764cdf5aec3daef862d143c5f09134d75b0")
	address1, _ = common.StringToAddress("P1QJNzZhqGoxNL2igkdthNBQLNWdNGTWzQU")

	prvKey2, _  = crypto.FromWIF("Ky7gQF2rxXLjGSymCtCMa67N2YMt98fRgxyy5WfH92FpbWDxWVRM") //"3892859c02b1be2ce494e61c60181051d79ff21dca22fae1dc349887335b6676" //P1N4nEffoUskPrbnoEqBR69JQDX2vv9vYa8\
	pubKey2B, _ = hex.DecodeString("02a2ba6f2a6e1467334d032ec54ac862c655d7e8bd6bbbce36c771fcdc0ddfb01f")
	address2, _ = common.StringToAddress("P1N4nEffoUskPrbnoEqBR69JQDX2vv9vYa8")

	prvKey3, _  = crypto.FromWIF("KzRHTanikQgR5oqUts69JTrCXRuy9Zod5qXdnAbYwvUnuUDJ3Rro") //"5f7754e5407fc2a81f453645cbd92878a6341d30dbfe2e680fc81628d47e8023" //P1MzuBUT7ubGpkAFqUB6chqTSXmBThQv2HT
	pubKey3B, _ = hex.DecodeString("020945d0c9ed05cf5ca9fe38dde799d7b73f4a5bfb71fc1a3c1dca79e2d86462a7")
	address3, _ = common.StringToAddress("P1MzuBUT7ubGpkAFqUB6chqTSXmBThQv2HT")

	prvKey4, _  = crypto.FromWIF("L3nZf9ds5JG5Sq2WMCxP6QSfHK6WuSpnsU8Qk2ygfGD92h553xhx") //"c3ecda5c797ef8d7ded2d332eb1cb83198ef88ede1bf9de7b60910644b45f83f" //P1MzuBUT7ubGpkAFqUB6chqTSXmBThQv2HT
	pubKey4B, _ = hex.DecodeString("0342ccc3459303c6a24fd3382249af438763c7fab9ca57e919aec658f7d05eab68")
	address4, _ = common.StringToAddress("P1Lcf8CTxgUwmFamn2qM7SrAukNyezakAbK")
)

func build23Address() ([]byte, []byte, string) {

	redeemScript := GenerateRedeemScript(2, [][]byte{pubKey1B, pubKey2B, pubKey3B})
	lockScript := GenerateP2SHLockScript(crypto.Hash160(redeemScript))
	addressMulti, _ := GetAddressFromScript(lockScript)

	return lockScript, redeemScript, addressMulti.Str()
}

//构造一个2/3签名的地址和UTXO，然后用其中的2个私钥对其进行签名
func TestMultiSign1Step(t *testing.T) {
	lockScript, redeemScript, addressMulti := build23Address()
	t.Logf("MultiSign Address:%s\n", addressMulti)
	t.Logf("RedeemScript: %x\n", redeemScript)
	r, _ := DisasmString(redeemScript)
	t.Logf("RedeemScript: %s\n", r)
	tx := &modules.Transaction{
		TxMessages: make([]*modules.Message, 0),
	}
	asset0 := &modules.Asset{}
	payment := &modules.PaymentPayload{}
	utxoTxId, _ := common.NewHashFromStr("1111870aa8c894376dbd960a22171d0ad7be057a730e14d7103ed4a6dbb34873")
	outPoint := modules.NewOutPoint(utxoTxId, 0, 0)
	txIn := modules.NewTxIn(outPoint, []byte{})
	payment.AddTxIn(txIn)
	p1lockScript := GenerateP2PKHLockScript(crypto.Hash160(pubKey1B))
	payment.AddTxOut(modules.NewTxOut(1, p1lockScript, asset0))
	tx.TxMessages = append(tx.TxMessages, modules.NewMessage(modules.APP_PAYMENT, payment))
	//scriptCp:=make([]byte,len(lockScript))
	//copy(scriptCp,lockScript)
	//privKeys := map[common.Address]*ecdsa.PrivateKey{
	//	address1: prvKey1,
	//	address2: prvKey2,
	//}
	getPubKeyFn := func(addr common.Address) ([]byte, error) {
		if addr == address1 {
			return crypto.CompressPubkey(&prvKey1.PublicKey), nil
		}
		if addr == address2 {
			return crypto.CompressPubkey(&prvKey2.PublicKey), nil
		}
		return nil, nil
	}
	getSignFn := func(addr common.Address, hash []byte) ([]byte, error) {

		if addr == address1 {
			return crypto.Sign(hash, prvKey1)
		}
		if addr == address2 {
			return crypto.Sign(hash, prvKey2)
		}
		return nil, nil
	}
	sign12, err := MultiSignOnePaymentInput(tx, 0, 0, lockScript, redeemScript, getPubKeyFn, getSignFn, nil, 0)
	if err != nil {
		t.Logf("Sign error:%s", err)
	}
	t.Logf("PrvKey1&2 sign result:%x\n", sign12)
	pay1 := tx.TxMessages[0].Payload.(*modules.PaymentPayload)
	pay1.Inputs[0].SignatureScript = sign12
	str, _ := txscript.DisasmString(sign12)
	t.Logf("Signed script:{%s}", str)

	err = ScriptValidate(lockScript, nil, tx, 0, 0)
	assert.Nil(t, err, fmt.Sprintf("validate error:%s", err))
}

//构造一个2/3签名的地址和UTXO，然后用其中的2个私钥分两步对其进行签名
func TestMultiSign2Step(t *testing.T) {
	lockScript, redeemScript, addressMulti := build23Address()
	t.Logf("MultiSign Address:%s\n", addressMulti)
	t.Logf("RedeemScript: %x\n", redeemScript)
	t.Logf("RedeemScript: %d\n", redeemScript)
	tx := &modules.Transaction{
		TxMessages: make([]*modules.Message, 0),
	}
	payment := &modules.PaymentPayload{}
	utxoTxId, _ := common.NewHashFromStr("1111870aa8c894376dbd960a22171d0ad7be057a730e14d7103ed4a6dbb34873")
	outPoint := modules.NewOutPoint(utxoTxId, 0, 0)
	txIn := modules.NewTxIn(outPoint, []byte{})
	payment.AddTxIn(txIn)
	p1lockScript := GenerateP2PKHLockScript(crypto.Hash160(pubKey1B))
	asset0 := &modules.Asset{}
	payment.AddTxOut(modules.NewTxOut(1, p1lockScript, asset0))
	tx.TxMessages = append(tx.TxMessages, modules.NewMessage(modules.APP_PAYMENT, payment))
	//scriptCp:=make([]byte,len(lockScript))
	//copy(scriptCp,lockScript)
	//privKeys := map[common.Address]*ecdsa.PrivateKey{
	//	address1: prvKey1,
	//}
	getPubKeyFn := func(addr common.Address) ([]byte, error) {
		if addr == address1 {
			return crypto.CompressPubkey(&prvKey1.PublicKey), nil
		}
		if addr == address2 {
			return crypto.CompressPubkey(&prvKey2.PublicKey), nil
		}
		return nil, nil
	}
	getSignFn := func(addr common.Address, hash []byte) ([]byte, error) {

		if addr == address1 {
			return crypto.Sign(hash, prvKey1)
		}
		if addr == address2 {
			return crypto.Sign(hash, prvKey2)
		}
		return nil, nil
	}
	sign1, err := MultiSignOnePaymentInput(tx, 0, 0, lockScript, redeemScript, getPubKeyFn, getSignFn, nil, 0)
	if err != nil {
		t.Logf("Sign error:%s", err)
	}
	t.Logf("PrvKey1 sign result:%x\n", sign1)
	pay1 := tx.TxMessages[0].Payload.(*modules.PaymentPayload)
	pay1.Inputs[0].SignatureScript = sign1

	//privKeys2 := map[common.Address]*ecdsa.PrivateKey{
	//	address2: prvKey2,
	//}
	//scriptCp2:=make([]byte,len(lockScript))
	//copy(scriptCp2,lockScript)
	sign2, err := MultiSignOnePaymentInput(tx, 0, 0, lockScript, redeemScript, getPubKeyFn, getSignFn, sign1, 0)
	if err != nil {
		t.Logf("Sign error:%s", err)
	}
	t.Logf("PrvKey2 sign result:%x\n", sign2)

	pay1 = tx.TxMessages[0].Payload.(*modules.PaymentPayload)
	pay1.Inputs[0].SignatureScript = sign2
	str, _ := txscript.DisasmString(sign2)
	t.Logf("Signed script:{%s}", str)

	err = ScriptValidate(lockScript, nil, tx, 0, 0)
	assert.Nil(t, err, fmt.Sprintf("validate error:%s", err))
}
func mockPickupJuryRedeemScript(addr common.Address, ver int) ([]byte, error) {
	return GenerateRedeemScript(2, [][]byte{pubKey1B, pubKey2B, pubKey3B, pubKey4B}), nil
}
func TestContractPayout(t *testing.T) {
	tx := &modules.Transaction{
		TxMessages: make([]*modules.Message, 0),
	}
	asset0 := &modules.Asset{}
	payment := &modules.PaymentPayload{}
	utxoTxId, _ := common.NewHashFromStr("1111870aa8c894376dbd960a22171d0ad7be057a730e14d7103ed4a6dbb34873")
	outPoint := modules.NewOutPoint(utxoTxId, 0, 0)
	txIn := modules.NewTxIn(outPoint, []byte{})
	payment.AddTxIn(txIn)
	p1lockScript := GenerateP2PKHLockScript(crypto.Hash160(pubKey1B))
	payment.AddTxOut(modules.NewTxOut(1, p1lockScript, asset0))
	tx.TxMessages = append(tx.TxMessages, modules.NewMessage(modules.APP_PAYMENT, payment))
	//scriptCp:=make([]byte,len(lockScript))
	//copy(scriptCp,lockScript)
	contractAddr, _ := common.StringToAddress("PCGTta3M4t3yXu8uRgkKvaWd2d8DR32W9vM")
	lockScript := GenerateP2CHLockScript(contractAddr) //Token 锁定到保证金合约中
	l, _ := txscript.DisasmString(lockScript)
	t.Logf("Lock Script:%s", l)
	//privKeys := map[common.Address]*ecdsa.PrivateKey{
	//	address1: prvKey1,
	//	address2: prvKey2,
	//}
	getPubKeyFn := func(addr common.Address) ([]byte, error) {
		if addr == address1 {
			return crypto.CompressPubkey(&prvKey1.PublicKey), nil
		}
		if addr == address2 {
			return crypto.CompressPubkey(&prvKey2.PublicKey), nil
		}
		if addr == address3 {
			return crypto.CompressPubkey(&prvKey3.PublicKey), nil
		}
		if addr == address4 {
			return crypto.CompressPubkey(&prvKey4.PublicKey), nil
		}
		return nil, nil
	}
	getSignFn := func(addr common.Address, hash []byte) ([]byte, error) {

		if addr == address1 {
			return crypto.Sign(hash, prvKey1)
		}
		if addr == address2 {
			return crypto.Sign(hash, prvKey2)
		}
		if addr == address3 {
			return crypto.Sign(hash, prvKey3)
		}
		if addr == address4 {
			return crypto.Sign(hash, prvKey4)
		}
		return nil, nil
	}
	redeemScript, _ := mockPickupJuryRedeemScript(contractAddr, 1)
	r, _ := txscript.DisasmString(redeemScript)
	t.Logf("RedeemScript:%s", r)
	sign12, err := MultiSignOnePaymentInput(tx, 0, 0, lockScript, redeemScript, getPubKeyFn, getSignFn, nil, 1)
	if err != nil {
		t.Logf("Sign error:%s", err)
	}
	t.Logf("PrvKey1&2 sign result:%x\n", sign12)
	pay1 := tx.TxMessages[0].Payload.(*modules.PaymentPayload)
	pay1.Inputs[0].SignatureScript = sign12
	str, _ := txscript.DisasmString(sign12)
	t.Logf("Signed script:{%s}", str)

	err = ScriptValidate(lockScript, mockPickupJuryRedeemScript, tx, 0, 0)
	assert.Nil(t, err, fmt.Sprintf("validate error:%s", err))
}
