// +build !linux

package kcp

import (
	"sync/atomic"
)

func (s *UDPSession) txLoop() {
	for {
		select {
		case txqueue := <-s.chTxQueue:
			nbytes := 0
			npkts := 0
			for k := range txqueue {
				if n, err := s.conn.WriteTo(txqueue[k], s.remote); err == nil {
					nbytes += n
					npkts++
				} else {
					s.notifyWriteError(err)
				}

				xmitBuf.Put(txqueue[k])
			}
			atomic.AddUint64(&DefaultSnmp.OutPkts, uint64(npkts))
			atomic.AddUint64(&DefaultSnmp.OutBytes, uint64(nbytes))
		case <-s.die:
			return
		}
	}
}