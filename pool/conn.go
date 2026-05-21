package pool

import (
	"net"
	"sync/atomic"
	"time"
)

type pooledConn struct {
	conn       *net.TCPConn
	createdAt  time.Time
	lastUsedAt time.Time
	useCount   atomic.Uint64
	isClosed   atomic.Bool
}

func NewpooledConn(conn *net.TCPConn) *pooledConn {
	return &pooledConn{
		conn:       conn,
		createdAt:  time.Now(),
		lastUsedAt: time.Now(),
	}
}

func (pc *pooledConn) markUsed() {
	pc.lastUsedAt = time.Now()
	pc.useCount.Add(1)
}

func (pc *pooledConn) idleFor() time.Duration {
	return time.Since(pc.lastUsedAt)
}

func (pc *pooledConn) age() time.Duration {
	return time.Since(pc.createdAt)
}

func (pc *pooledConn) close() error {
	if pc.isClosed.CompareAndSwap(false, true) {
		return pc.conn.Close()
	}
	return nil
}
