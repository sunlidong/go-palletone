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

package common

import (
	"testing"
)

func TestAddressValidate(t *testing.T) {
	p2pkh := "P1Kp2hcLhGEP45Xgx7vmSrE37QXunJUd8gJ"
	addr, err := StringToAddress(p2pkh)

	if err != nil {
		t.Error(err)
	}
	t.Log(addr)
}
func TestAddressNotValidate(t *testing.T) {
	p2pkh := "P1Kp2hcLhGEP45Xgx7vmSrE37QXunJUd8gj"
	addr, err := StringToAddress(p2pkh)

	if err != nil {
		t.Log(addr)
		t.Log(err)
	} else {
		t.Error("It must invalid, but pass, please check your code")
	}

}
func TestHexToAddrString(t *testing.T) {
	addr := HexToAddress("0x00000000000000000000000000000000000000011C")
	t.Logf("0x1 contract address: %s", addr.String())
	newAddr, _ := StringToAddress(addr.String())
	t.Logf("contract hex is: %x", newAddr.Bytes())
	addr = HexToAddress("0x00000000000000000000000000000000000000021C")
	t.Logf("0x2 contract address: %s", addr.String())
	addr = HexToAddress("0x00000000000000000000000000000000000000031C")
	t.Logf("0x3 contract address: %s", addr.String())

	addr = HexToAddress("0x000000000000000000000000000000000000000100")
	t.Logf("0x1 user address: %s", addr.String())

	addr = HexToAddress("0x000000000000000000000000000000000000000105")
	t.Logf("0x1 p2sh address: %s", addr.String())
}
