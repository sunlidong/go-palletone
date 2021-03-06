/*
   This file is part of go-palletone.
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
*/
/*
 * @author PalletOne core developers <dev@pallet.one>
 * @date 2018
 */

package txspool

import (
	"container/heap"
	"math/big"

	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/dag/modules"
)

// priceHeap is a heap.Interface implementation over transactions for retrieving
// price-sorted transactions to discard when the pool fills up.
type priceHeap []*modules.TxPoolTransaction

func (h priceHeap) Len() int           { return len(h) }
func (h priceHeap) Less(i, j int) bool { return h[i].GetTxFee().Cmp(h[j].GetTxFee()) < 0 }
func (h priceHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index, h[j].Index = i, j
}

func (h *priceHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*modules.TxPoolTransaction)
	item.Index = n
	*h = append(*h, item)
}

// -1 标识该数据已经出了优先级队列了
func (h *priceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	x.Index = -1
	return x
}

func (h *priceHeap) Update(item *modules.TxPoolTransaction, priority float64) {
	item.Priority_lvl = priority
	heap.Fix(h, item.Index)
}

type priorityHeap []*modules.TxPoolTransaction

func (h priorityHeap) Len() int           { return len(h) }
func (h priorityHeap) Less(i, j int) bool { return h[i].Priority_lvl > h[j].Priority_lvl }
func (h priorityHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	//h[i].Index, h[j].Index = i, j
}

func (h *priorityHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*modules.TxPoolTransaction)
	item.Index = n
	*h = append(*h, item)
}

// -1 标识该数据已经出了优先级队列了 ,弹出优先级最高的
func (h *priorityHeap) Pop() interface{} {
	//old := *h
	//n := len(old)
	//x := old[0]
	//*h = old[1:n]
	//x.Index = -1
	//return x
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	x.Index = -1
	return x
}

func (h *priorityHeap) Update(item *modules.TxPoolTransaction, priority float64) {
	item.Priority_lvl = priority
	heap.Fix(h, item.Index)
}

// txPricedList is a price-sorted heap to allow operating on transactions pool
// contents in a price-incrementing way.
type txPricedList struct {
	all    *map[common.Hash]*modules.TxPoolTransaction // Pointer to the map of all transactions
	items  *priorityHeap                               // Heap of prices of all the stored transactions
	stales int                                         // Number of stale price points to (re-heap trigger)
}

// newTxPricedList creates a new price-sorted transaction heap.
func newTxPricedList(all *map[common.Hash]*modules.TxPoolTransaction) *txPricedList {
	return &txPricedList{
		all:   all,
		items: new(priorityHeap),
	}
}

// Put inserts a new transaction into the heap.
func (l *txPricedList) Put(tx *modules.TxPoolTransaction) *priorityHeap {
	heap.Push(l.items, tx)
	(*l.all)[tx.Tx.Hash()] = tx
	//sort.Sort(l.items)
	return l.items
}
func (l *txPricedList) Get() *modules.TxPoolTransaction {
	if l != nil {
		if l.items.Len() > 0 {
			tx, ok := heap.Pop(l.items).(*modules.TxPoolTransaction)
			if ok {
				if tx.Tx != nil {
					return tx
				}
			}
		}
	}
	return nil
}

// Removed notifies the prices transaction list that an old transaction dropped
// from the pool. The list will just keep a counter of stale objects and update
// the heap if a large enough ratio of transactions go stale.
func (l *txPricedList) Removed(hash common.Hash) {
	//// Bump the stale counter, but exit if still too low (< 25%)
	//l.stales++
	//if l.stales <= len(*l.items)/4 {
	//	return
	//}
	// Seems we've reached a critical number of stale transactions, reheap
	reheap := make(priorityHeap, 0, len(*l.all))

	l.stales, l.items = 0, &reheap

	for key, tx := range *l.all {
		if hash != key {
			*l.items = append(*l.items, tx)
		} else {
			l.stales--
		}
	}
	heap.Init(l.items)
}

// Cap finds all the transactions below the given price threshold, drops them
// from the priced list and returs them for further removal from the entire pool.
func (l *txPricedList) Cap(threshold *big.Int, local *utxoSet) modules.TxPoolTxs {
	drop := make(modules.TxPoolTxs, 0, 128) // Remote underpriced transactions to drop
	save := make(modules.TxPoolTxs, 0, 64)  // Local underpriced transactions to keep

	for len(*l.items) > 0 {
		// Discard stale transactions if found during cleanup
		tx := heap.Pop(l.items).(*modules.TxPoolTransaction)
		if _, ok := (*l.all)[tx.Tx.Hash()]; !ok {
			l.stales--
			continue
		}
		// Stop the discards if we've reached the threshold
		if tx.GetTxFee().Cmp(threshold) >= 0 {
			save = append(save, tx)
			break
		}
		// Non stale transaction found, discard unless local
		if local.containsTx(tx) {
			save = append(save, tx)
		} else {
			drop = append(drop, tx)
		}
	}
	for _, tx := range save {
		heap.Push(l.items, tx)
	}
	return drop
}

// Underpriced checks whether a transaction is cheaper than (or as cheap as) the
// lowest priced transaction currently being tracked.
func (l *txPricedList) Underpriced(tx *modules.TxPoolTransaction, local *utxoSet) bool {
	// Local transactions cannot be underpriced
	if local.containsTx(tx) {
		return false
	}
	// Discard stale price points if found at the heap start
	for len(*l.items) > 0 {
		head := []*modules.TxPoolTransaction(*l.items)[0]
		if _, ok := (*l.all)[head.Tx.Hash()]; !ok {
			l.stales--
			heap.Pop(l.items)
			continue
		}
		break
	}
	// Check if the transaction is underpriced or not
	if len(*l.items) == 0 {
		log.Error("Pricing query for empty pool") // This cannot happen, print to catch programming errors
		return false
	}
	cheapest := []*modules.TxPoolTransaction(*l.items)[0]
	return cheapest.Priority_lvl >= tx.Priority_lvl
}

// Discard finds a number of most underpriced transactions, removes them from the
// priced list and returns them for further removal from the entire pool.
func (l *txPricedList) Discard(count int, local *utxoSet) modules.TxPoolTxs {
	drop := make(modules.TxPoolTxs, 0, count) // Remote underpriced transactions to drop
	save := make(modules.TxPoolTxs, 0, 64)    // Local underpriced transactions to keep

	for len(*l.items) > 0 && count > 0 {
		// Discard stale transactions if found during cleanup
		tx := heap.Pop(l.items).(*modules.TxPoolTransaction)
		if _, ok := (*l.all)[tx.Tx.Hash()]; !ok {
			l.stales--
			continue
		}
		// Non stale transaction found, discard unless local
		if local.containsTx(tx) {
			save = append(save, tx)
		} else {
			drop = append(drop, tx)
			count--
		}
	}
	for _, tx := range save {
		heap.Push(l.items, tx)
	}
	return drop
}
