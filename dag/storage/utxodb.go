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

package storage

import (
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/common/ptndb"
	"github.com/palletone/go-palletone/common/rlp"
	"github.com/palletone/go-palletone/dag/constants"
	"github.com/palletone/go-palletone/dag/errors"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/tokenengine"
)

type UtxoDb struct {
	db ptndb.Database
	//logger log.ILogger
}

func NewUtxoDb(db ptndb.Database) *UtxoDb {
	return &UtxoDb{db: db}
}

type IUtxoDb interface {
	GetPrefix(prefix []byte) map[string][]byte

	GetUtxoEntry(outpoint *modules.OutPoint) (*modules.Utxo, error)
	//GetUtxoPkScripHexByTxhash(txhash common.Hash, mindex, outindex uint32) (string, error)
	//GetAddrOutput(addr string) ([]modules.Output, error)
	GetAddrOutpoints(addr common.Address) ([]modules.OutPoint, error)
	GetAddrUtxos(addr common.Address) (map[modules.OutPoint]*modules.Utxo, error)
	GetAllUtxos() (map[modules.OutPoint]*modules.Utxo, error)
	SaveUtxoEntity(outpoint *modules.OutPoint, utxo *modules.Utxo) error
	SaveUtxoView(view map[modules.OutPoint]*modules.Utxo) error
	DeleteUtxo(outpoint *modules.OutPoint) error
}

// ###################### UTXO index for Address ######################
// key: outpoint_prefix + addr + outpoint's Bytes
func (utxodb *UtxoDb) saveUtxoOutpoint(address common.Address, outpoint *modules.OutPoint) error {
	key := append(constants.AddrOutPoint_Prefix, address.Bytes()...)
	key = append(key, outpoint.Bytes()...)

	// save outpoint tofind address
	out_key := append(constants.OutPointAddr_Prefix, outpoint.ToKey()...)
	StoreBytes(utxodb.db, out_key, address.String())
	return StoreBytes(utxodb.db, key, outpoint)
}
func (utxodb *UtxoDb) batchSaveUtxoOutpoint(batch ptndb.Batch, address common.Address, outpoint *modules.OutPoint) error {
	key := append(constants.AddrOutPoint_Prefix, address.Bytes()...)
	key = append(key, outpoint.Bytes()...)
	val, err := rlp.EncodeToBytes(outpoint)
	if err != nil {
		return err
	}
	return batch.Put(key, val)
}
func (utxodb *UtxoDb) deleteUtxoOutpoint(address common.Address, outpoint *modules.OutPoint) error {
	key := append(constants.AddrOutPoint_Prefix, address.Bytes()...)
	key = append(key, outpoint.Bytes()...)
	return utxodb.db.Delete(key)
}
func (db *UtxoDb) GetAddrOutpoints(address common.Address) ([]modules.OutPoint, error) {
	data := db.GetPrefix(append(constants.AddrOutPoint_Prefix, address.Bytes()...))
	outpoints := make([]modules.OutPoint, 0)
	for _, b := range data {
		out := new(modules.OutPoint)
		if err := rlp.DecodeBytes(b, out); err == nil {
			outpoints = append(outpoints, *out)
		}
	}
	return outpoints, nil
}

// ###################### UTXO Entity ######################
func (utxodb *UtxoDb) SaveUtxoEntity(outpoint *modules.OutPoint, utxo *modules.Utxo) error {
	key := outpoint.ToKey()
	log.Debug("Try to save utxo by key:", "outpoint_key", outpoint.String())
	err := StoreBytes(utxodb.db, key, utxo)
	if err != nil {
		return err
	}
	address, _ := tokenengine.GetAddressFromScript(utxo.PkScript[:])
	return utxodb.saveUtxoOutpoint(address, outpoint)
}

// SaveUtxoView to update the utxo set in the database based on the provided utxo view.
func (utxodb *UtxoDb) SaveUtxoView(view map[modules.OutPoint]*modules.Utxo) error {
	batch := utxodb.db.NewBatch()
	log.Debugf("Start batch save utxo, batch count:%d", len(view))
	for outpoint, utxo := range view {
		// No need to update the database if the utxo was not modified.
		if utxo == nil { // || utxo.IsModified()
			continue
		} else {
			key := outpoint.ToKey()
			val, err := rlp.EncodeToBytes(utxo)
			if err != nil {
				return err
			}
			batch.Put(key, val)
			address, _ := tokenengine.GetAddressFromScript(utxo.PkScript[:])
			// save utxoindex and  addr and key
			utxodb.batchSaveUtxoOutpoint(batch, address, &outpoint)

		}
	}

	return batch.Write()
}

// Remove the utxo
func (utxodb *UtxoDb) DeleteUtxo(outpoint *modules.OutPoint) error {

	//1. get utxo
	utxo, err := utxodb.GetUtxoEntry(outpoint)
	if err != nil {
		log.Errorf("Try to soft delete an unknown utxo by key:%s", outpoint.String())
		return err
	}

	//2. soft delete utxo
	if utxo.IsSpent() {
		log.Errorf("Try to soft delete a deleted utxo by key:%s", outpoint.String())
		return errors.New("Try to soft delete a deleted utxo")
	}
	key := outpoint.ToKey()
	utxo.Spend()
	log.Debugf("Try to soft delete utxo by key:%s", outpoint.String())
	err = StoreBytes(utxodb.db, key, utxo)
	if err != nil {
		return err
	}
	//3. Remove index
	address, _ := tokenengine.GetAddressFromScript(utxo.PkScript[:])
	utxodb.deleteUtxoOutpoint(address, outpoint)
	//if err := utxodb.db.Delete(key); err != nil {
	//	return err
	//}

	return nil
}

//@Yiran
//func (utxodb *UtxoDb) SaveUtxoSnapshot(index *modules.ChainIndex) error {
//	//0. examine wrong calling
//	if index.Index%modules.TERMINTERVAL != 0 {
//		return errors.New("SaveUtxoSnapshot must wait until last term period end")
//	}
//	//1. get all utxo
//	utxos, err := utxodb.GetAllUtxos()
//	if err != nil {
//		return util.ErrorLogHandler(err, "utxodb.GetAllUtxos")
//	}
//	PTNutxos := make([]modules.Utxo, 0)
//	for _, utxo := range utxos {
//		if utxo.Asset.AssetId == modules.PTNCOIN {
//			PTNutxos = append(PTNutxos, *utxo)
//		}
//	}
//	//2. store utxo
//	key := util.KeyConnector([]byte(constants.UTXOSNAPSHOT_PREFIX), ConvertBytes(index))
//	return utxodb.SaveUtxoEntities(key, &PTNutxos)
//}

//func (utxodb *UtxoDb) GetUtxoSnapshot(index []byte) error {
//
//}

// ###################### SAVE IMPL END ######################

// ###################### GET IMPL START ######################
//  dbFetchUtxoEntry
func (utxodb *UtxoDb) GetUtxoEntry(outpoint *modules.OutPoint) (*modules.Utxo, error) {
	utxo := new(modules.Utxo)
	key := outpoint.ToKey()
	log.Debugf("Query utxo by outpoint:%s", outpoint.String())
	data, err := utxodb.db.Get(key)
	if err != nil {
		log.Error("get utxo entry failed", "error", err, "Query utxo by outpoint:%s", outpoint.String())
		if err.Error() == errors.ErrNotFound.Error() {
			return nil, errors.ErrUtxoNotFound
		}
		return nil, err
	}

	if err := rlp.DecodeBytes(data, &utxo); err != nil {
		return nil, err
	}
	return utxo, nil
}

//
//func (utxodb *UtxoDb) GetUtxoPkScripHexByTxhash(txhash common.Hash, mindex, outindex uint32) (string, error) {
//	outpoint := &modules.OutPoint{TxHash: txhash, MessageIndex: mindex, OutIndex: outindex}
//	utxo, err := utxodb.GetUtxoEntry(outpoint)
//	if err != nil {
//		return "", err
//	}
//	if utxo == nil {
//		return "", errors.New("get the pkscript is failed,the utxo is null.")
//	}
//	return hexutil.Encode(utxo.PkScript), nil
//}

func (utxodb *UtxoDb) GetUtxoByIndex(indexKey []byte) ([]byte, error) {
	return utxodb.db.Get(indexKey)
}

//func (db *UtxoDb) GetAddrOutput(addr string) ([]modules.Output, error) {
//
//	data := db.GetPrefix(append(constants.AddrOutput_Prefix, []byte(addr)...))
//	outputs := make([]modules.Output, 0)
//	for _, b := range data {
//		out := new(modules.Output)
//		if err := rlp.DecodeBytes(b, out); err == nil {
//			outputs = append(outputs, *out)
//		}
//	}
//	return outputs, nil
//}

func (db *UtxoDb) GetAddrUtxos(addr common.Address) (map[modules.OutPoint]*modules.Utxo, error) {
	allutxos := make(map[modules.OutPoint]*modules.Utxo, 0)
	outpoints, err := db.GetAddrOutpoints(addr)
	if err != nil {
		return nil, err
	}
	for _, out := range outpoints {
		if utxo, err := db.GetUtxoEntry(&out); err == nil {
			if !utxo.IsSpent() {
				allutxos[out] = utxo
			}
		}
	}
	return allutxos, nil
}
func (db *UtxoDb) GetAllUtxos() (map[modules.OutPoint]*modules.Utxo, error) {
	view := make(map[modules.OutPoint]*modules.Utxo, 0)

	items := db.GetPrefix(constants.UTXO_PREFIX)
	var err error
	for key, itme := range items {
		utxo := new(modules.Utxo)
		// outpint := new(modules.OutPoint)
		if err = rlp.DecodeBytes(itme, utxo); err == nil {
			if !utxo.IsSpent() {
				outpoint := modules.KeyToOutpoint([]byte(key))
				view[*outpoint] = utxo
			}
		}
	}

	return view, err
}

// get prefix: return maps
func (db *UtxoDb) GetPrefix(prefix []byte) map[string][]byte {
	return getprefix(db.db, prefix)
}

// ###################### GET IMPL END ######################
