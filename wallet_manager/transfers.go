package wallet_manager

import (
	"database/sql"
	"runtime"
	"sort"
	"sync"

	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
)

type GetTransfersParams struct {
	In                       sql.NullBool
	Out                      sql.NullBool
	Coinbase                 sql.NullBool
	Sender                   sql.NullString
	Receiver                 sql.NullString
	BurnGreaterOrEqualThan   sql.NullInt64
	AmountGreaterOrEqualThan sql.NullInt64
	TXID                     sql.NullString
	BlockHash                sql.NullString
	SC_ID                    sql.NullString
	SC_Entrypoint            sql.NullString
	Offset                   sql.NullInt64
	Limit                    sql.NullInt64
}

func filterEntries(allEntries []rpc.Entry, start, end int, params GetTransfersParams, entryChan chan<- rpc.Entry, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := start; i < end; i++ {
		e := allEntries[i]
		add := true

		if params.Coinbase.Valid {
			add = e.Coinbase == params.Coinbase.Bool
		}

		if params.In.Valid {
			add = e.Incoming == params.In.Bool
		}

		if params.Out.Valid {
			add = (!e.Incoming && !e.Coinbase) == params.Out.Bool
		}

		if params.Sender.Valid {
			add = e.Sender == params.Sender.String
		}

		if params.Receiver.Valid {
			add = e.Destination == params.Receiver.String
		}

		if params.AmountGreaterOrEqualThan.Valid {
			add = e.Amount >= uint64(params.AmountGreaterOrEqualThan.Int64)
		}

		if params.BurnGreaterOrEqualThan.Valid {
			add = e.Burn >= uint64(params.BurnGreaterOrEqualThan.Int64)
		}

		if params.TXID.Valid {
			add = e.TXID == params.TXID.String
		}

		if params.BlockHash.Valid {
			add = e.BlockHash == params.BlockHash.String
		}

		if params.SC_ID.Valid {
			add = false
			for _, arg := range e.SCDATA {
				if arg.Name == "SC_ID" {
					hash, ok := arg.Value.(crypto.Hash)
					if ok && hash.String() == params.SC_ID.String {
						add = true
						break
					}
				}
			}
		}

		if add && params.SC_Entrypoint.Valid {
			add = false
			for _, arg := range e.SCDATA {
				if arg.Name == "entrypoint" {
					entrypoint, ok := arg.Value.(string)
					if ok && entrypoint == params.SC_Entrypoint.String {
						add = true
						break
					}
				}
			}
		}

		if add {
			entryChan <- e
		}
	}
}

func (w *Wallet) GetTransfers(scId string, params GetTransfersParams) []rpc.Entry {
	w.Memory.Lock()
	defer w.Memory.Unlock()

	account := w.Memory.GetAccount()
	allEntries := account.EntriesNative[crypto.HashHexToHash(scId)]
	totalEntries := len(allEntries)
	if allEntries == nil || totalEntries < 1 {
		return allEntries
	}

	workers := runtime.NumCPU()
	var wg sync.WaitGroup
	entryChan := make(chan rpc.Entry)

	chunkSize := totalEntries / workers
	if chunkSize < 50 {
		chunkSize = totalEntries
		workers = 1
	}

	var entries []rpc.Entry
	done := make(chan bool)
	go func() {
		for e := range entryChan {
			entries = append(entries, e)
		}

		done <- true
	}()

	for i := 0; i < workers; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if i == workers-1 {
			end = totalEntries
		}

		wg.Add(1)
		go filterEntries(allEntries, start, end, params, entryChan, &wg)
	}

	wg.Wait()
	close(entryChan)
	<-done

	sort.Slice(entries, func(a, b int) bool {
		return entries[a].Time.Unix() > entries[b].Time.Unix()
	})

	if params.Offset.Valid {
		offset := params.Offset.Int64
		if len(entries) > int(offset) {
			entries = entries[offset:]
		}
	}

	if params.Limit.Valid {
		limit := params.Limit.Int64
		if len(entries) > int(limit) {
			entries = entries[:limit]
		}
	}

	return entries
}
