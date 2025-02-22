package network

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DataDog/datadog-agent/pkg/network/http"
	"github.com/DataDog/datadog-agent/pkg/process/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkStoreClosedConnection(b *testing.B) {
	conns := generateRandConnections(30000)
	for _, bench := range []struct {
		connCount int
	}{
		{
			connCount: 100,
		},
		{
			connCount: 1000,
		},
		{
			connCount: 10000,
		},
		{
			connCount: 30000,
		},
	} {
		b.Run(fmt.Sprintf("StoreClosedConnection-%d", bench.connCount), func(b *testing.B) {
			ns := newDefaultState()
			ns.GetDelta(DEBUGCLIENT, latestEpochTime(), nil, nil, nil) // Initial fetch to set up client

			b.ResetTimer()
			b.ReportAllocs()

			for n := 0; n < b.N; n++ {
				for _, c := range conns[:bench.connCount] {
					ns.StoreClosedConnection(&c)
				}
			}
		})
	}
}

func BenchmarkConnectionsGet(b *testing.B) {
	conns := generateRandConnections(30000)
	closed := generateRandConnections(30000)

	for _, bench := range []struct {
		connCount   int
		closedCount int
	}{
		{
			connCount:   100,
			closedCount: 0,
		},
		{
			connCount:   100,
			closedCount: 50,
		},
		{
			connCount:   100,
			closedCount: 100,
		},
		{
			connCount:   1000,
			closedCount: 0,
		},
		{
			connCount:   1000,
			closedCount: 500,
		},
		{
			connCount:   1000,
			closedCount: 1000,
		},
		{
			connCount:   10000,
			closedCount: 0,
		},
		{
			connCount:   10000,
			closedCount: 5000,
		},
		{
			connCount:   10000,
			closedCount: 10000,
		},
		{
			connCount:   30000,
			closedCount: 0,
		},
		{
			connCount:   30000,
			closedCount: 15000,
		},
		{
			connCount:   30000,
			closedCount: 30000,
		},
	} {
		b.Run(fmt.Sprintf("ConnectionsGet-%d-%d", bench.connCount, bench.closedCount), func(b *testing.B) {
			ns := newDefaultState()

			// Initial fetch to set up client
			ns.GetDelta(DEBUGCLIENT, latestTime, nil, nil, nil)

			for _, c := range closed[:bench.closedCount] {
				ns.StoreClosedConnection(&c)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for n := 0; n < b.N; n++ {
				ns.GetDelta(DEBUGCLIENT, latestTime, conns[:bench.connCount], nil, nil)
			}
		})
	}
}

func TestRemoveConnections(t *testing.T) {
	conn := ConnectionStats{
		Pid:                  123,
		Type:                 UDP,
		Family:               AFINET,
		Source:               util.AddressFromString("127.0.0.1"),
		Dest:                 util.AddressFromString("127.0.0.1"),
		SPort:                31890,
		DPort:                80,
		MonotonicSentBytes:   12345,
		LastSentBytes:        12345,
		MonotonicRecvBytes:   6789,
		LastRecvBytes:        6789,
		MonotonicRetransmits: 2,
		LastRetransmits:      2,
		IntraHost:            true,
	}

	buf := make([]byte, ConnectionByteKeyMaxLen)
	key, err := conn.ByteKey(buf)
	require.NoError(t, err)

	clientID := "1"
	state := newDefaultState().(*networkState)
	conns := state.GetDelta(clientID, latestEpochTime(), nil, nil, nil).Connections
	assert.Equal(t, 0, len(conns))

	conns = state.GetDelta(clientID, latestEpochTime(), []ConnectionStats{conn}, nil, nil).Connections
	assert.Equal(t, 1, len(conns))
	assert.Equal(t, conn, conns[0])

	client := state.clients[clientID]
	assert.Equal(t, 1, len(client.stats))

	state.RemoveConnections([]string{string(key)})
	assert.Equal(t, 0, len(client.stats))
}

func TestRetrieveClosedConnection(t *testing.T) {
	conn := ConnectionStats{
		Pid:                  123,
		Type:                 TCP,
		Family:               AFINET,
		Source:               util.AddressFromString("127.0.0.1"),
		Dest:                 util.AddressFromString("127.0.0.1"),
		SPort:                31890,
		DPort:                80,
		MonotonicSentBytes:   12345,
		LastSentBytes:        12345,
		MonotonicRecvBytes:   6789,
		LastRecvBytes:        6789,
		MonotonicRetransmits: 2,
		LastRetransmits:      2,
		IntraHost:            true,
	}

	clientID := "1"

	t.Run("without prior registration", func(t *testing.T) {
		state := newDefaultState()
		state.StoreClosedConnection(&conn)
		conns := state.GetDelta(clientID, latestEpochTime(), nil, nil, nil).Connections

		assert.Equal(t, 0, len(conns))
	})

	t.Run("with registration", func(t *testing.T) {
		state := newDefaultState()

		conns := state.GetDelta(clientID, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 0, len(conns))

		state.StoreClosedConnection(&conn)

		conns = state.GetDelta(clientID, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, conn, conns[0])

		// An other client that is not registered should not have the closed connection
		conns = state.GetDelta("2", latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 0, len(conns))

		// It should no more have connections stored
		conns = state.GetDelta(clientID, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 0, len(conns))
	})
}

func TestCleanupClient(t *testing.T) {
	clientID := "1"

	state := NewState(100*time.Millisecond, 50000, 75000, 75000, 75000, false)
	clients := state.(*networkState).getClients()
	assert.Equal(t, 0, len(clients))

	conns := state.GetDelta(clientID, latestEpochTime(), nil, nil, nil).Connections
	assert.Equal(t, 0, len(conns))

	// Should be a no op
	state.(*networkState).RemoveExpiredClients(time.Now())

	clients = state.(*networkState).getClients()
	assert.Equal(t, 1, len(clients))
	assert.Equal(t, "1", clients[0])

	// Should delete the client 1
	state.(*networkState).RemoveExpiredClients(time.Now().Add(150 * time.Millisecond))

	clients = state.(*networkState).getClients()
	assert.Equal(t, 0, len(clients))
}

func TestLastStats(t *testing.T) {
	client1 := "1"
	client2 := "2"
	state := newDefaultState()

	dSent := uint64(42)
	dRecv := uint64(133)
	dRetransmits := uint32(7)

	conn := ConnectionStats{
		Pid:                  123,
		Type:                 TCP,
		Family:               AFINET,
		Source:               util.AddressFromString("127.0.0.1"),
		Dest:                 util.AddressFromString("127.0.0.1"),
		SPort:                31890,
		DPort:                80,
		MonotonicSentBytes:   36,
		MonotonicRecvBytes:   24,
		MonotonicRetransmits: 2,
	}

	conn2 := conn
	conn2.MonotonicSentBytes += dSent
	conn2.MonotonicRecvBytes += dRecv
	conn2.MonotonicRetransmits += dRetransmits

	conn3 := conn2
	conn3.MonotonicSentBytes += dSent
	conn3.MonotonicRecvBytes += dRecv
	conn3.MonotonicRetransmits += dRetransmits

	// First get, we should not have any connections stored
	conns := state.GetDelta(client1, latestEpochTime(), nil, nil, nil).Connections
	assert.Equal(t, 0, len(conns))

	// Same for an other client
	conns = state.GetDelta(client2, latestEpochTime(), nil, nil, nil).Connections
	assert.Equal(t, 0, len(conns))

	// We should have only one connection but with last stats equal to monotonic
	conns = state.GetDelta(client1, latestEpochTime(), []ConnectionStats{conn}, nil, nil).Connections
	assert.Equal(t, 1, len(conns))
	assert.Equal(t, conn.MonotonicSentBytes, conns[0].LastSentBytes)
	assert.Equal(t, conn.MonotonicRecvBytes, conns[0].LastRecvBytes)
	assert.Equal(t, conn.MonotonicRetransmits, conns[0].LastRetransmits)
	assert.Equal(t, conn.MonotonicSentBytes, conns[0].MonotonicSentBytes)
	assert.Equal(t, conn.MonotonicRecvBytes, conns[0].MonotonicRecvBytes)
	assert.Equal(t, conn.MonotonicRetransmits, conns[0].MonotonicRetransmits)

	// This client didn't collect the first connection so last stats = monotonic
	conns = state.GetDelta(client2, latestEpochTime(), []ConnectionStats{conn2}, nil, nil).Connections
	assert.Equal(t, 1, len(conns))
	assert.Equal(t, conn2.MonotonicSentBytes, conns[0].LastSentBytes)
	assert.Equal(t, conn2.MonotonicRecvBytes, conns[0].LastRecvBytes)
	assert.Equal(t, conn2.MonotonicRetransmits, conns[0].LastRetransmits)
	assert.Equal(t, conn2.MonotonicSentBytes, conns[0].MonotonicSentBytes)
	assert.Equal(t, conn2.MonotonicRecvBytes, conns[0].MonotonicRecvBytes)
	assert.Equal(t, conn2.MonotonicRetransmits, conns[0].MonotonicRetransmits)

	// client 1 should have conn3 - conn1 since it did not collected conn2
	conns = state.GetDelta(client1, latestEpochTime(), []ConnectionStats{conn3}, nil, nil).Connections
	assert.Equal(t, 1, len(conns))
	assert.Equal(t, 2*dSent, conns[0].LastSentBytes)
	assert.Equal(t, 2*dRecv, conns[0].LastRecvBytes)
	assert.Equal(t, 2*dRetransmits, conns[0].LastRetransmits)
	assert.Equal(t, conn3.MonotonicSentBytes, conns[0].MonotonicSentBytes)
	assert.Equal(t, conn3.MonotonicRecvBytes, conns[0].MonotonicRecvBytes)
	assert.Equal(t, conn3.MonotonicRetransmits, conns[0].MonotonicRetransmits)

	// client 2 should have conn3 - conn2
	conns = state.GetDelta(client2, latestEpochTime(), []ConnectionStats{conn3}, nil, nil).Connections
	assert.Equal(t, 1, len(conns))
	assert.Equal(t, dSent, conns[0].LastSentBytes)
	assert.Equal(t, dRecv, conns[0].LastRecvBytes)
	assert.Equal(t, dRetransmits, conns[0].LastRetransmits)
	assert.Equal(t, conn3.MonotonicSentBytes, conns[0].MonotonicSentBytes)
	assert.Equal(t, conn3.MonotonicRecvBytes, conns[0].MonotonicRecvBytes)
	assert.Equal(t, conn3.MonotonicRetransmits, conns[0].MonotonicRetransmits)
}

func TestLastStatsForClosedConnection(t *testing.T) {
	clientID := "1"
	state := newDefaultState()

	dSent := uint64(42)
	dRecv := uint64(133)
	dRetransmits := uint32(0)

	conn := ConnectionStats{
		Pid:                  123,
		Type:                 TCP,
		Family:               AFINET,
		Source:               util.AddressFromString("127.0.0.1"),
		Dest:                 util.AddressFromString("127.0.0.1"),
		SPort:                31890,
		DPort:                80,
		MonotonicSentBytes:   36,
		MonotonicRecvBytes:   24,
		MonotonicRetransmits: 1,
	}

	conn2 := conn
	conn2.MonotonicSentBytes += dSent
	conn2.MonotonicRecvBytes += dRecv
	conn2.MonotonicRetransmits += dRetransmits

	// First get, we should not have any connections stored
	conns := state.GetDelta(clientID, latestEpochTime(), nil, nil, nil).Connections
	assert.Equal(t, 0, len(conns))

	// We should have one connection with last stats equal to monotonic stats
	conns = state.GetDelta(clientID, latestEpochTime(), []ConnectionStats{conn}, nil, nil).Connections
	assert.Equal(t, 1, len(conns))
	assert.Equal(t, conn.MonotonicSentBytes, conns[0].LastSentBytes)
	assert.Equal(t, conn.MonotonicRecvBytes, conns[0].LastRecvBytes)
	assert.Equal(t, conn.MonotonicRetransmits, conns[0].LastRetransmits)
	assert.Equal(t, conn.MonotonicSentBytes, conns[0].MonotonicSentBytes)
	assert.Equal(t, conn.MonotonicRecvBytes, conns[0].MonotonicRecvBytes)
	assert.Equal(t, conn.MonotonicRetransmits, conns[0].MonotonicRetransmits)

	state.StoreClosedConnection(&conn2)

	// We should have one connection with last stats
	conns = state.GetDelta(clientID, latestEpochTime(), nil, nil, nil).Connections

	assert.Equal(t, 1, len(conns))
	assert.Equal(t, dSent, conns[0].LastSentBytes)
	assert.Equal(t, dRecv, conns[0].LastRecvBytes)
	assert.Equal(t, dRetransmits, conns[0].LastRetransmits)
	assert.Equal(t, conn2.MonotonicSentBytes, conns[0].MonotonicSentBytes)
	assert.Equal(t, conn2.MonotonicRecvBytes, conns[0].MonotonicRecvBytes)
	assert.Equal(t, conn2.MonotonicRetransmits, conns[0].MonotonicRetransmits)
}

func TestRaceConditions(t *testing.T) {
	nClients := 10

	// Generate random conns
	genConns := func(n uint32) []ConnectionStats {
		conns := make([]ConnectionStats, 0, n)
		for i := uint32(0); i < n; i++ {
			conns = append(conns, ConnectionStats{
				Pid:                  1 + i,
				Type:                 TCP,
				Family:               AFINET,
				Source:               util.AddressFromString("127.0.0.1"),
				Dest:                 util.AddressFromString("127.0.0.1"),
				SPort:                uint16(rand.Int()),
				DPort:                uint16(rand.Int()),
				MonotonicSentBytes:   uint64(rand.Int()),
				MonotonicRecvBytes:   uint64(rand.Int()),
				MonotonicRetransmits: uint32(rand.Int()),
			})
		}
		return conns
	}

	state := newDefaultState()
	nConns := uint32(100)

	var wg sync.WaitGroup
	wg.Add(nClients)

	// Spawn multiple clients to get multiple times
	for i := 1; i <= nClients; i++ {
		go func(c string) {
			defer wg.Done()
			defer state.RemoveClient(c)
			timer := time.NewTimer(1 * time.Second)
			for {
				select {
				case <-timer.C:
					return
				default:
					state.GetDelta(c, latestEpochTime(), genConns(nConns), nil, nil)
				}
			}
		}(fmt.Sprintf("%d", i))
	}

	wg.Wait()
}

func TestSameKeyEdgeCases(t *testing.T) {
	// For this test all the connections have the same key
	// Each vertical bar represents a collection for a given client
	// Each horizontal bar represents a connection lifespan (from start to end with the number of sent bytes written on top of the line)

	client := "c"
	conn := ConnectionStats{
		Pid:                123,
		Type:               TCP,
		Family:             AFINET,
		Source:             util.AddressFromString("127.0.0.1"),
		Dest:               util.AddressFromString("127.0.0.1"),
		MonotonicSentBytes: 3,
	}

	t.Run("ShortlivedConnection", func(t *testing.T) {
		// +     3 bytes      +
		// |                  |
		// |   +---------+    |
		// |                  |
		// +                  +

		// c0                 c1

		// We expect:
		// c0: Nothing
		// c1: Monotonic: 3 bytes, Last seen: 3 bytes
		state := newDefaultState()

		// First get, we should have nothing
		conns := state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 0, len(conns))

		// Store the connection as closed
		state.StoreClosedConnection(&conn)

		// Second get, we should have monotonic and last stats = 3
		conns = state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 3, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 3, int(conns[0].LastSentBytes))
	})

	t.Run("TwoShortlivedConnections", func(t *testing.T) {
		//  +    3 bytes       5 bytes    +
		//  |                             |
		//  |    +-----+       +-----+    |
		//  |                             |
		//  +                             +

		//  c0                            c1

		// We expect:
		// c0: Nothing
		// c1: Monotonic: 8 bytes, Last seenL 8 bytes

		state := newDefaultState()

		// First get, we should have nothing
		conns := state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 0, len(conns))

		// Store the connection as closed
		state.StoreClosedConnection(&conn)

		conn2 := conn
		conn2.MonotonicSentBytes = 5
		conn2.LastUpdateEpoch++
		// Store the connection another time
		state.StoreClosedConnection(&conn2)

		// Second get, we should have monotonic and last stats = 8
		conns = state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 8, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 8, int(conns[0].LastSentBytes))
	})

	t.Run("TwoShortlivedConnectionsCrossing-1", func(t *testing.T) {
		// +    1 b  +  1 bytes    1 b +   1 b        +
		// |         |                 |              |
		// |    +-----------+      +------------+     |
		// |         |                 |              |
		// +         +                 +              +

		// c0        c1                c2             c3
		// We expect:

		// c0: Nothing
		// c1: Monotonic: 1 bytes, Last seen: 1 bytes
		// c2: Monotonic: 3 bytes, Last seen: 2 bytes
		// c3: Monotonic: 2 bytes, Last seen: 1 bytes

		state := newDefaultState()

		// First get for client c, we should have nothing
		conns := state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
		assert.Len(t, conns, 0)

		conn := ConnectionStats{
			Pid:                123,
			Type:               TCP,
			Family:             AFINET,
			Source:             util.AddressFromString("127.0.0.1"),
			Dest:               util.AddressFromString("127.0.0.1"),
			SPort:              9000,
			DPort:              1234,
			MonotonicSentBytes: 1,
		}

		// Simulate this connection starting
		conns = state.GetDelta(client, latestEpochTime(), []ConnectionStats{conn}, nil, nil).Connections
		require.Len(t, conns, 1)
		assert.EqualValues(t, 1, conns[0].LastSentBytes)
		assert.EqualValues(t, 1, conns[0].MonotonicSentBytes)

		// Store the connection as closed
		conn.MonotonicSentBytes++
		conn.LastUpdateEpoch = latestEpochTime()
		state.StoreClosedConnection(&conn)

		conn.MonotonicSentBytes = 1
		conn.LastUpdateEpoch = latestEpochTime()
		// Retrieve the connections
		conns = state.GetDelta(client, latestEpochTime(), []ConnectionStats{conn}, nil, nil).Connections
		require.Len(t, conns, 1)
		assert.EqualValues(t, 2, conns[0].LastSentBytes)
		assert.EqualValues(t, 3, conns[0].MonotonicSentBytes)

		conn.MonotonicSentBytes++
		conn.LastUpdateEpoch = latestEpochTime()
		// Store the connection as closed
		state.StoreClosedConnection(&conn)

		conns = state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
		require.Len(t, conns, 1)
		assert.EqualValues(t, 1, conns[0].LastSentBytes)
		assert.EqualValues(t, 2, conns[0].MonotonicSentBytes)
	})

	t.Run("TwoShortlivedConnectionsCrossing-2", func(t *testing.T) {
		// +    3 bytes    2 b  +  3 bytes    1 b +   2 b        +
		// |                    |                 |              |
		// |    +-----+    +-----------+      +------------+     |
		// |                    |                 |              |
		// +                    +                 +              +

		// c0                   c1                c2             c3
		// We expect:

		// c0: Nothing
		// c1: Monotonic: 5 bytes, Last seen: 5 bytes
		// c2: Monotonic: 6 bytes, Last seen: 4 bytes
		// c3: Monotonic: 3 bytes, Last seen: 2 bytes

		state := newDefaultState()

		// First get, we should have nothing
		conns := state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 0, len(conns))

		// Store the connection as closed
		state.StoreClosedConnection(&conn)

		conn2 := conn
		conn2.MonotonicSentBytes = 2
		conn2.LastUpdateEpoch++
		// Store the connection as an opened connection
		cs := []ConnectionStats{conn2}

		// Second get, we should have monotonic and last stats = 5
		conns = state.GetDelta(client, latestEpochTime(), cs, nil, nil).Connections
		require.Equal(t, 1, len(conns))
		assert.Equal(t, 5, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 5, int(conns[0].LastSentBytes))

		// Store the connection as closed
		conn2.MonotonicSentBytes += 3
		conn2.LastUpdateEpoch++
		state.StoreClosedConnection(&conn2)

		// Store the connection again
		conn3 := conn2
		conn3.MonotonicSentBytes = 1
		conn3.LastUpdateEpoch++
		cs = []ConnectionStats{conn3}

		// Third get, we should have monotonic = 6 and last stats = 4
		conns = state.GetDelta(client, latestEpochTime(), cs, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 6, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 4, int(conns[0].LastSentBytes))

		// Store the connection as closed
		conn3.MonotonicSentBytes += 2
		state.StoreClosedConnection(&conn3)

		// 4th get, we should have monotonic = 3 and last stats = 2
		conns = state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 3, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 2, int(conns[0].LastSentBytes))
	})

	t.Run("ConnectionCrossing", func(t *testing.T) {
		// 3 b  +  5 bytes        +
		//      |                 |
		// +-----------+          |
		//      |                 |
		//      +                 +

		//     c0                c1
		// We expect:

		// c0: Monotonic: 3 bytes, Last seen: 3 bytes
		// c1: Monotonic: 8 bytes, Last seen: 5 bytes

		state := newDefaultState()

		// this is to register we should not have anything
		conns := state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 0, len(conns))

		// Store the connection as opened
		cs := []ConnectionStats{conn}

		// First get, we should have monotonic = 3 and last seen = 3
		conns = state.GetDelta(client, latestEpochTime(), cs, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 3, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 3, int(conns[0].LastSentBytes))

		// Store the connection as closed
		conn2 := conn
		conn2.MonotonicSentBytes = 8
		state.StoreClosedConnection(&conn2)

		// Second get, we should have monotonic = 8 and last stats = 5
		conns = state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 8, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 5, int(conns[0].LastSentBytes))
	})

	t.Run("TwoShortlivedConnectionsCrossingWithTwoClients", func(t *testing.T) {
		//              +    3 bytes    2 b  +  3 bytes    1 b +   2 b        +
		//              |                    |                 |              |
		// client c     |    +-----+    +-----------+      +------------+     |
		//              |                    |                 |              |
		//              +                    +                 +              +
		//
		//              c0                   c1                c2             c3
		//
		//
		//              +    3 bytes  +  3 b    +  2 b      2 b     +  1 b         +
		//              |             |         |                   |              |
		// client d     |    +-----+  |  +----------+      +------------+          |
		//              |             |         |                   |              |
		//              +             +         +                   +              +
		//
		//              d0            d1        d2                  d3             d4

		// We expect:
		// c0: Nothing
		// d0: Nothing
		// d1: Monotonic: 3 bytes, Last seen: 3 bytes (this connection started after closed + collect, so we reset monotonic)
		// c1: Monotonic: 5 bytes, Last seen: 5 bytes
		// d2: Monotonic: 3 bytes, Last seen 3 bytes
		// c2: Monotonic: 6 bytes, Last seen: 4 bytes
		// d3: Monotonic: 7 bytes, Last seen 4 bytes
		// c3: Monotonic: 3 bytes, Last seen: 2 bytes
		// d4: Monotonic: 3 bytes, Last seen: 1 bytes

		clientD := "d"

		state := newDefaultState()

		// First get for client c, we should have nothing
		conns := state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 0, len(conns))

		// First get for client d, we should have nothing
		conns = state.GetDelta(clientD, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 0, len(conns))

		// Store the connection as closed
		state.StoreClosedConnection(&conn)

		// Second get for client d we should have monotonic and last stats = 3
		conns = state.GetDelta(clientD, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 3, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 3, int(conns[0].LastSentBytes))

		// Store the connection as an opened connection
		conn2 := conn
		conn2.MonotonicSentBytes = 2
		conn2.LastUpdateEpoch++
		cs := []ConnectionStats{conn2}

		// Second get, for client c we should have monotonic and last stats = 5
		conns = state.GetDelta(client, latestEpochTime(), cs, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 5, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 5, int(conns[0].LastSentBytes))

		// Store the connection as an opened connection
		conn2.MonotonicSentBytes++
		conn2.LastUpdateEpoch++
		cs = []ConnectionStats{conn2}

		// Third get, for client d we should have monotonic = 3 and last stats = 3
		conns = state.GetDelta(clientD, latestEpochTime(), cs, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 3, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 3, int(conns[0].LastSentBytes))

		// Store the connection as closed
		conn2.MonotonicSentBytes += 2
		conn2.LastUpdateEpoch++
		state.StoreClosedConnection(&conn2)

		// Store the connection again
		conn3 := conn2
		conn3.MonotonicSentBytes = 1
		conn3.LastUpdateEpoch++
		cs = []ConnectionStats{conn3}

		// Third get, for client c, we should have monotonic = 6 and last stats = 4
		conns = state.GetDelta(client, latestEpochTime(), cs, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 6, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 4, int(conns[0].LastSentBytes))

		// Store the connection again
		conn3.MonotonicSentBytes++
		conn3.LastUpdateEpoch++
		cs = []ConnectionStats{conn3}

		// 4th get, for client d, we should have monotonic = 7 and last stats = 4
		conns = state.GetDelta(clientD, latestEpochTime(), cs, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 7, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 4, int(conns[0].LastSentBytes))

		// Store the connection as closed
		conn3.MonotonicSentBytes++
		conn3.LastUpdateEpoch++
		state.StoreClosedConnection(&conn3)

		// 4th get, for client c we should have monotonic = 3 and last stats = 2
		conns = state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 3, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 2, int(conns[0].LastSentBytes))

		// 5th get, for client d we should have monotonic = 3 and last stats = 1
		conns = state.GetDelta(clientD, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 3, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 1, int(conns[0].LastSentBytes))
	})

	t.Run("ShortlivedConnectionCrossingWithThreeClients", func(t *testing.T) {
		//              +    3 bytes    2 b  +  3 bytes
		//              |                    |
		// client c     |    +-----+    +-----------+
		//              |                    |
		//              +                    +
		//
		//              c0                   c1
		//
		//
		//              +    3 bytes  +  3 b    +  2 b
		//              |             |         |
		// client d     |    +-----+  |  +----------+
		//              |             |         |
		//              +             +         +
		//
		//              d0            d1        d2
		//
		//
		//              +    2 b + 1b  +    5 bytes   +
		//              |        |     |              |
		// client e     |    +-----+   | +---------+  |
		//              |        |     |              |
		//              +        +     +              +
		//
		//              e0       e1    e2             e3

		// We expect:
		// c0, d0, e0: Nothing
		// e1: Monotonic: 2 bytes, Last seen 2 bytes
		// d1: Monotonic 3 bytes, Last seen: 3 bytes
		// e2: Monotonic: 3 bytes, Last seen: 1 bytes
		// c1: Monotonic: 5 bytes, Last seen: 5 bytes
		// d2: Monotonic: 3 bytes, Last seen 3 bytes
		// e3: Monotonic: 5 bytes, Last seen: 5 bytes

		clientD := "d"
		clientE := "e"

		state := newDefaultState()

		// First get for client c, we should have nothing
		conns := state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 0, len(conns))

		// First get for client d, we should have nothing
		conns = state.GetDelta(clientD, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 0, len(conns))

		// First get for client e, we should have nothing
		conns = state.GetDelta(clientE, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 0, len(conns))

		// Store the connection
		conn.MonotonicSentBytes = 2
		conn.LastUpdateEpoch++
		cs := []ConnectionStats{conn}

		// Second get for client e we should have monotonic and last stats = 2
		conns = state.GetDelta(clientE, latestEpochTime(), cs, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 2, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 2, int(conns[0].LastSentBytes))

		// Store the connection as closed
		conn.MonotonicSentBytes++
		conn.LastUpdateEpoch++
		state.StoreClosedConnection(&conn)

		// Second get for client d we should have monotonic and last stats = 3
		conns = state.GetDelta(clientD, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 3, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 3, int(conns[0].LastSentBytes))

		// Third get for client e we should have monotonic = 3and last stats = 1
		conns = state.GetDelta(clientE, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 3, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 1, int(conns[0].LastSentBytes))

		// Store the connection as an opened connection
		conn2 := conn
		conn2.MonotonicSentBytes = 2
		conn2.LastUpdateEpoch++
		cs = []ConnectionStats{conn2}

		// Second get, for client c we should have monotonic and last stats = 5
		conns = state.GetDelta(client, latestEpochTime(), cs, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 5, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 5, int(conns[0].LastSentBytes))

		// Store the connection as an opened connection
		conn2.MonotonicSentBytes++
		conn2.LastUpdateEpoch++
		cs = []ConnectionStats{conn2}

		// Third get, for client d we should have monotonic = 3 and last stats = 3
		conns = state.GetDelta(clientD, latestEpochTime(), cs, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 3, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 3, int(conns[0].LastSentBytes))

		// Store the connection as closed
		conn2.MonotonicSentBytes += 2
		conn2.LastUpdateEpoch++
		state.StoreClosedConnection(&conn2)

		// 4th get, for client e we should have monotonic = 5 and last stats = 5
		conns = state.GetDelta(clientE, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 1, len(conns))
		assert.Equal(t, 5, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 5, int(conns[0].LastSentBytes))
	})

	t.Run("LonglivedConnectionWithTwoClientsJoiningAtDifferentTimes", func(t *testing.T) {
		//              +      3 bytes       +  1 + 3 b        +   2 b
		//              |                    |                 |
		// client c     |    +------------------------------------------+
		//              |                    |                 |
		//              +                    +                 +
		//
		//              c0                   c1                c2
		//
		//                                                5 bytes
		//                                        +                      +
		//                                        |                      |
		// client d                               |---------------------+|
		//                                        |                      |
		//                                        +                      +
		//
		//                                       d0                      d1

		// We expect:
		// c0: Nothing
		// c1: Monotonic: 3 bytes, Last seen: 3 bytes
		// d0: Monotonic: 4 bytes, Last seen: 0 bytes
		// c2: Monotonic: 7 bytes, Last seen: 4 bytes
		// d1: Monotonic: 9 bytes, Last seen: 5 bytes

		clientD := "d"

		state := newDefaultState()

		// First get for client c, we should have nothing
		conns := state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
		assert.Equal(t, 0, len(conns))

		// Second get for client c we should have monotonic and last stats = 3
		conns = state.GetDelta(client, latestEpochTime(), []ConnectionStats{conn}, nil, nil).Connections
		assert.Len(t, conns, 1)
		assert.Equal(t, 3, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 3, int(conns[0].LastSentBytes))

		conn2 := conn
		conn2.MonotonicSentBytes++
		conn2.LastUpdateEpoch++

		// First get for client d we should have monotonic = 4 and last bytes = 4
		conns = state.GetDelta(clientD, latestEpochTime(), []ConnectionStats{conn2}, nil, nil).Connections
		assert.Len(t, conns, 1)
		assert.Equal(t, 4, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 0, int(conns[0].LastSentBytes))

		conn3 := conn2
		conn3.MonotonicSentBytes += 3
		conn3.LastUpdateEpoch++

		// Third get for client c we should have monotonic = 7 and last bytes = 4
		conns = state.GetDelta(client, latestEpochTime(), []ConnectionStats{conn3}, nil, nil).Connections
		assert.Len(t, conns, 1)
		assert.Equal(t, 7, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 4, int(conns[0].LastSentBytes))

		conn4 := conn3
		conn4.MonotonicSentBytes += 2
		conn4.LastUpdateEpoch++

		// Second get for client d we should have monotonic = 9 and last bytes = 5
		conns = state.GetDelta(clientD, latestEpochTime(), []ConnectionStats{conn4}, nil, nil).Connections
		assert.Len(t, conns, 1)
		assert.Equal(t, 9, int(conns[0].MonotonicSentBytes))
		assert.Equal(t, 5, int(conns[0].LastSentBytes))
	})
}

func TestStatsResetOnUnderflow(t *testing.T) {
	conn := ConnectionStats{
		Pid:                123,
		Type:               TCP,
		Family:             AFINET,
		Source:             util.AddressFromString("127.0.0.1"),
		Dest:               util.AddressFromString("127.0.0.1"),
		MonotonicSentBytes: 3,
		IntraHost:          true,
	}

	client := "client"

	state := newDefaultState()

	// Register the client
	assert.Len(t, state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections, 0)

	// Get the connections once to register stats
	conns := state.GetDelta(client, latestEpochTime(), []ConnectionStats{conn}, nil, nil).Connections
	require.Len(t, conns, 1)

	// Expect LastStats to be 3
	conn.LastSentBytes = 3
	assert.Equal(t, conn, conns[0])

	// Get the connections again but by simulating an underflow
	conn.MonotonicSentBytes--

	conns = state.GetDelta(client, latestEpochTime(), []ConnectionStats{conn}, nil, nil).Connections
	require.Len(t, conns, 1)
	expected := conn
	expected.LastSentBytes = 2
	// We expect the LastStats to be 2
	assert.Equal(t, expected, conns[0])
}

func TestDoubleCloseOnTwoClients(t *testing.T) {
	conn := ConnectionStats{
		Pid:                123,
		Type:               TCP,
		Family:             AFINET,
		Source:             util.AddressFromString("127.0.0.1"),
		Dest:               util.AddressFromString("127.0.0.1"),
		MonotonicSentBytes: 3,
		LastSentBytes:      3,
		IntraHost:          true,
	}

	expectedConn := conn
	expectedConn.MonotonicSentBytes *= 2
	expectedConn.LastSentBytes *= 2

	client1 := "1"
	client2 := "2"

	state := newDefaultState()

	// Register the clients
	assert.Len(t, state.GetDelta(client1, latestEpochTime(), nil, nil, nil).Connections, 0)
	assert.Len(t, state.GetDelta(client2, latestEpochTime(), nil, nil, nil).Connections, 0)

	// Store the closed connection twice
	state.StoreClosedConnection(&conn)
	conn.LastUpdateEpoch++
	state.StoreClosedConnection(&conn)

	expectedConn.LastUpdateEpoch = conn.LastUpdateEpoch
	// Get the connections for client1 we should have only one with stats = 2*conn
	conns := state.GetDelta(client1, latestEpochTime(), nil, nil, nil).Connections
	require.Len(t, conns, 1)
	assert.Equal(t, expectedConn, conns[0])

	// Same for client2
	conns = state.GetDelta(client2, latestEpochTime(), nil, nil, nil).Connections
	require.Len(t, conns, 1)
	assert.Equal(t, expectedConn, conns[0])
}

func TestUnorderedCloseEvent(t *testing.T) {
	conn := ConnectionStats{
		Pid:                123,
		Type:               TCP,
		Family:             AFINET,
		Source:             util.AddressFromString("127.0.0.1"),
		Dest:               util.AddressFromString("127.0.0.1"),
		MonotonicSentBytes: 3,
	}

	client := "client"
	state := newDefaultState()

	// Register the client
	assert.Len(t, state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections, 0)

	// Simulate storing a closed connection while we were reading from the eBPF map
	// in this case the closed conn will have an earlier epoch
	conn.LastUpdateEpoch = latestEpochTime() + 1
	conn.MonotonicSentBytes++
	conn.MonotonicRecvBytes = 1
	state.StoreClosedConnection(&conn)

	conn.LastUpdateEpoch--
	conn.MonotonicSentBytes--
	conn.MonotonicRecvBytes = 0
	conns := state.GetDelta(client, latestEpochTime(), []ConnectionStats{conn}, nil, nil).Connections
	require.Len(t, conns, 1)
	assert.EqualValues(t, 4, conns[0].LastSentBytes)
	assert.EqualValues(t, 1, conns[0].LastRecvBytes)

	// Simulate some other gets
	assert.Len(t, state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections, 0)
	assert.Len(t, state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections, 0)
	assert.Len(t, state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections, 0)

	// Simulate having the connection getting active again
	conn.LastUpdateEpoch = latestEpochTime()
	conn.MonotonicSentBytes--
	state.StoreClosedConnection(&conn)

	conns = state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections
	require.Len(t, conns, 1)
	assert.EqualValues(t, 2, conns[0].LastSentBytes)
	assert.EqualValues(t, 0, conns[0].LastRecvBytes)

	// Ensure we don't have underflows / unordered conns
	assert.Zero(t, state.(*networkState).telemetry.statsResets)

	assert.Len(t, state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections, 0)
}

func TestAggregateClosedConnectionsTimestamp(t *testing.T) {
	conn := ConnectionStats{
		Pid:                123,
		Type:               TCP,
		Family:             AFINET,
		Source:             util.AddressFromString("127.0.0.1"),
		Dest:               util.AddressFromString("127.0.0.1"),
		MonotonicSentBytes: 3,
	}

	client := "client"
	state := newDefaultState()

	// Register the client
	assert.Len(t, state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections, 0)

	conn.LastUpdateEpoch = latestEpochTime()
	state.StoreClosedConnection(&conn)

	conn.LastUpdateEpoch = latestEpochTime()
	state.StoreClosedConnection(&conn)

	conn.LastUpdateEpoch = latestEpochTime()
	state.StoreClosedConnection(&conn)

	// Make sure the connections we get has the latest timestamp
	delta := state.GetDelta(client, latestEpochTime(), nil, nil, nil)
	assert.Equal(t, conn.LastUpdateEpoch, delta.Connections[0].LastUpdateEpoch)
}

func TestDNSStatsWithMultipleClients(t *testing.T) {
	c := ConnectionStats{
		Pid:    123,
		Type:   TCP,
		Family: AFINET,
		Source: util.AddressFromString("127.0.0.1"),
		Dest:   util.AddressFromString("127.0.0.1"),
		SPort:  1000,
		DPort:  53,
	}

	dKey := DNSKey{clientIP: c.Source, clientPort: c.SPort, serverIP: c.Dest, protocol: c.Type}

	getStats := func() map[DNSKey]map[string]map[QueryType]DNSStats {
		var d = "foo.com"
		statsByDomain := make(map[DNSKey]map[string]map[QueryType]DNSStats)
		stats := make(map[QueryType]DNSStats)
		countByRcode := make(map[uint32]uint32)
		countByRcode[uint32(DNSResponseCodeNoError)] = 1
		stats[DNSTypeA] = DNSStats{DNSCountByRcode: countByRcode}
		statsByDomain[dKey] = make(map[string]map[QueryType]DNSStats)
		statsByDomain[dKey][d] = stats
		return statsByDomain
	}

	client1 := "client1"
	client2 := "client2"
	client3 := "client3"
	state := newDefaultState()

	// Register the first two clients
	assert.Len(t, state.GetDelta(client1, latestEpochTime(), nil, nil, nil).Connections, 0)
	assert.Len(t, state.GetDelta(client2, latestEpochTime(), nil, nil, nil).Connections, 0)

	c.LastUpdateEpoch = latestEpochTime()
	state.StoreClosedConnection(&c)

	conns := state.GetDelta(client1, latestEpochTime(), nil, getStats(), nil).Connections
	require.Len(t, conns, 1)
	assert.EqualValues(t, 1, conns[0].DNSSuccessfulResponses)

	// Register the third client but also pass in dns stats
	conns = state.GetDelta(client3, latestEpochTime(), []ConnectionStats{c}, getStats(), nil).Connections
	require.Len(t, conns, 1)
	// DNS stats should be available for the new client
	assert.EqualValues(t, 1, conns[0].DNSSuccessfulResponses)

	conns = state.GetDelta(client2, latestEpochTime(), []ConnectionStats{c}, getStats(), nil).Connections
	require.Len(t, conns, 1)
	// 2nd client should get accumulated stats
	assert.EqualValues(t, 3, conns[0].DNSSuccessfulResponses)
}

func TestDNSStatsWithMultipleClientsWithDomainCollectionEnabled(t *testing.T) {
	c := ConnectionStats{
		Pid:    123,
		Type:   TCP,
		Family: AFINET,
		Source: util.AddressFromString("127.0.0.1"),
		Dest:   util.AddressFromString("127.0.0.1"),
		SPort:  1000,
		DPort:  53,
	}

	dKey := DNSKey{clientIP: c.Source, clientPort: c.SPort, serverIP: c.Dest, protocol: c.Type}
	var d = "foo.com"
	getStats := func() map[DNSKey]map[string]map[QueryType]DNSStats {
		statsByDomain := make(map[DNSKey]map[string]map[QueryType]DNSStats)
		stats := make(map[QueryType]DNSStats)
		countByRcode := make(map[uint32]uint32)
		countByRcode[uint32(DNSResponseCodeNoError)] = 1
		stats[DNSTypeA] = DNSStats{DNSCountByRcode: countByRcode}
		statsByDomain[dKey] = make(map[string]map[QueryType]DNSStats)
		statsByDomain[dKey][d] = stats
		return statsByDomain
	}

	client1 := "client1"
	client2 := "client2"
	client3 := "client3"
	state := NewState(2*time.Minute, 50000, 75000, 75000, 7500, true)

	// Register the first two clients
	assert.Len(t, state.GetDelta(client1, latestEpochTime(), nil, nil, nil).Connections, 0)
	assert.Len(t, state.GetDelta(client2, latestEpochTime(), nil, nil, nil).Connections, 0)

	c.LastUpdateEpoch = latestEpochTime()
	state.StoreClosedConnection(&c)

	conns := state.GetDelta(client1, latestEpochTime(), nil, getStats(), nil).Connections
	require.Len(t, conns, 1)
	assert.EqualValues(t, 1, conns[0].DNSStatsByDomainByQueryType[d][DNSTypeA].DNSCountByRcode[DNSResponseCodeNoError])
	// domain agnostic stats should be 0
	assert.EqualValues(t, 0, conns[0].DNSSuccessfulResponses)

	// Register the third client but also pass in dns stats
	conns = state.GetDelta(client3, latestEpochTime(), []ConnectionStats{c}, getStats(), nil).Connections
	require.Len(t, conns, 1)
	// DNS stats should be available for the new client
	assert.EqualValues(t, 1, conns[0].DNSStatsByDomainByQueryType[d][DNSTypeA].DNSCountByRcode[DNSResponseCodeNoError])
	// domain agnostic stats should be 0
	assert.EqualValues(t, 0, conns[0].DNSSuccessfulResponses)

	conns = state.GetDelta(client2, latestEpochTime(), []ConnectionStats{c}, getStats(), nil).Connections
	require.Len(t, conns, 1)
	// 2nd client should get accumulated stats
	assert.EqualValues(t, 3, conns[0].DNSStatsByDomainByQueryType[d][DNSTypeA].DNSCountByRcode[DNSResponseCodeNoError])
	// domain agnostic stats should be 0
	assert.EqualValues(t, 0, conns[0].DNSSuccessfulResponses)
}

func TestDNSStatsPIDCollisions(t *testing.T) {
	c := ConnectionStats{
		Pid:    123,
		Type:   TCP,
		Family: AFINET,
		Source: util.AddressFromString("127.0.0.1"),
		Dest:   util.AddressFromString("127.0.0.1"),
		SPort:  1000,
		DPort:  53,
	}

	var d = "foo.com"
	dKey := DNSKey{clientIP: c.Source, clientPort: c.SPort, serverIP: c.Dest, protocol: c.Type}
	statsByDomain := make(map[DNSKey]map[string]map[QueryType]DNSStats)
	stats := make(map[QueryType]DNSStats)
	countByRcode := make(map[uint32]uint32)
	countByRcode[uint32(DNSResponseCodeNoError)] = 1
	stats[DNSTypeA] = DNSStats{DNSCountByRcode: countByRcode}
	statsByDomain[dKey] = make(map[string]map[QueryType]DNSStats)
	statsByDomain[dKey][d] = stats

	client := "client"
	state := newDefaultState()

	// Register the client
	assert.Len(t, state.GetDelta(client, latestEpochTime(), nil, nil, nil).Connections, 0)

	c.LastUpdateEpoch = latestEpochTime()
	state.StoreClosedConnection(&c)

	// Store another connection with same DNSKey but different PID
	c.Pid++
	state.StoreClosedConnection(&c)

	conns := state.GetDelta(client, latestEpochTime(), nil, statsByDomain, nil).Connections
	require.Len(t, conns, 2)
	successes := 0
	for _, conn := range conns {
		successes += int(conn.DNSSuccessfulResponses)
	}
	assert.Equal(t, 1, successes)
	assert.Equal(t, int64(1), state.(*networkState).telemetry.dnsPidCollisions)
}

func TestHTTPStats(t *testing.T) {
	c := ConnectionStats{
		Source: util.AddressFromString("1.1.1.1"),
		Dest:   util.AddressFromString("0.0.0.0"),
		SPort:  1000,
		DPort:  80,
	}

	key := http.NewKey(c.Source, c.Dest, c.SPort, c.DPort, "/testpath", http.MethodGet)

	httpStats := make(map[http.Key]http.RequestStats)
	var rs http.RequestStats
	httpStats[key] = rs

	// Register client & pass in HTTP stats
	state := newDefaultState()
	delta := state.GetDelta("client", latestEpochTime(), []ConnectionStats{c}, nil, httpStats)

	// Verify connection has HTTP data embedded in it
	assert.Len(t, delta.HTTP, 1)

	// Verify HTTP data has been flushed
	delta = state.GetDelta("client", latestEpochTime(), []ConnectionStats{c}, nil, nil)
	assert.Len(t, delta.HTTP, 0)
}

func TestHTTPStatsWithMultipleClients(t *testing.T) {
	c := ConnectionStats{
		Source: util.AddressFromString("1.1.1.1"),
		Dest:   util.AddressFromString("0.0.0.0"),
		SPort:  1000,
		DPort:  80,
	}

	getStats := func(path string) map[http.Key]http.RequestStats {
		httpStats := make(map[http.Key]http.RequestStats)
		key := http.NewKey(c.Source, c.Dest, c.SPort, c.DPort, path, http.MethodGet)
		var rs http.RequestStats
		httpStats[key] = rs
		return httpStats
	}

	client1 := "client1"
	client2 := "client2"
	client3 := "client3"
	state := newDefaultState()

	// Register the first two clients
	assert.Len(t, state.GetDelta(client1, latestEpochTime(), nil, nil, nil).HTTP, 0)
	assert.Len(t, state.GetDelta(client2, latestEpochTime(), nil, nil, nil).HTTP, 0)

	// Store the connection to both clients & pass HTTP stats to the first client
	c.LastUpdateEpoch = latestEpochTime()
	state.StoreClosedConnection(&c)

	delta := state.GetDelta(client1, latestEpochTime(), nil, nil, getStats("/testpath"))
	assert.Len(t, delta.HTTP, 1)

	// Verify that the HTTP stats were also stored in the second client
	delta = state.GetDelta(client2, latestEpochTime(), nil, nil, nil)
	assert.Len(t, delta.HTTP, 1)

	// Register a third client & verify that it does not have the HTTP stats
	delta = state.GetDelta(client3, latestEpochTime(), []ConnectionStats{c}, nil, nil)
	assert.Len(t, delta.HTTP, 0)

	c.LastUpdateEpoch = latestEpochTime()
	state.StoreClosedConnection(&c)

	// Pass in new HTTP stats to the first client
	delta = state.GetDelta(client1, latestEpochTime(), nil, nil, getStats("/testpath2"))
	assert.Len(t, delta.HTTP, 1)

	// And the second client
	delta = state.GetDelta(client2, latestEpochTime(), nil, nil, getStats("/testpath3"))
	assert.Len(t, delta.HTTP, 2)

	// Verify that the third client also accumulated both new HTTP stats
	delta = state.GetDelta(client3, latestEpochTime(), nil, nil, nil)
	assert.Len(t, delta.HTTP, 2)
}

func TestDetermineConnectionIntraHost(t *testing.T) {
	tests := []struct {
		name      string
		conn      ConnectionStats
		intraHost bool
	}{
		{
			name: "equal source/dest",
			conn: ConnectionStats{
				Source: util.AddressFromString("1.1.1.1"),
				Dest:   util.AddressFromString("1.1.1.1"),
				SPort:  123,
				DPort:  456,
			},
			intraHost: true,
		},
		{
			name: "source/dest loopback",
			conn: ConnectionStats{
				Source: util.AddressFromString("127.0.0.1"),
				Dest:   util.AddressFromString("127.0.0.1"),
				SPort:  123,
				DPort:  456,
			},
			intraHost: true,
		},
		{
			name: "dest nat'ed to loopback",
			conn: ConnectionStats{
				Source: util.AddressFromString("1.1.1.1"),
				Dest:   util.AddressFromString("2.2.2.2"),
				SPort:  123,
				DPort:  456,
				IPTranslation: &IPTranslation{
					ReplSrcIP:   util.AddressFromString("127.0.0.1"),
					ReplDstIP:   util.AddressFromString("1.1.1.1"),
					ReplSrcPort: 456,
					ReplDstPort: 123,
				},
			},
			intraHost: true,
		},
		{
			name: "local connection with nat on both sides",
			conn: ConnectionStats{
				Source:    util.AddressFromString("1.1.1.1"),
				Dest:      util.AddressFromString("169.254.169.254"),
				SPort:     12345,
				DPort:     80,
				Direction: OUTGOING,
				NetNS:     1212,
				IPTranslation: &IPTranslation{
					ReplSrcIP:   util.AddressFromString("127.0.0.1"),
					ReplDstIP:   util.AddressFromString("1.1.1.1"),
					ReplSrcPort: 8181,
					ReplDstPort: 12345,
				},
			},
			intraHost: true,
		},
		{
			name: "local connection with nat on both sides",
			conn: ConnectionStats{
				Source:    util.AddressFromString("127.0.0.1"),
				Dest:      util.AddressFromString("1.1.1.1"),
				SPort:     8181,
				DPort:     12345,
				Direction: INCOMING,
				NetNS:     1233,
				IPTranslation: &IPTranslation{
					ReplSrcIP:   util.AddressFromString("1.1.1.1"),
					ReplDstIP:   util.AddressFromString("169.254.169.254"),
					ReplSrcPort: 12345,
					ReplDstPort: 80,
				},
			},
			intraHost: true,
		},
		{
			name: "remote connection with source translation (redirect)",
			conn: ConnectionStats{
				Source:    util.AddressFromString("4.4.4.4"),
				Dest:      util.AddressFromString("2.2.2.2"),
				SPort:     12345,
				DPort:     80,
				Direction: INCOMING,
				NetNS:     2,
				IPTranslation: &IPTranslation{
					ReplSrcIP:   util.AddressFromString("2.2.2.2"),
					ReplDstIP:   util.AddressFromString("127.0.0.1"),
					ReplSrcPort: 12345,
					ReplDstPort: 15006,
				},
			},
			intraHost: false,
		},
		{
			name: "local connection, same network ns",
			conn: ConnectionStats{
				Source:    util.AddressFromString("1.1.1.1"),
				Dest:      util.AddressFromString("2.2.2.2"),
				SPort:     12345,
				DPort:     80,
				Direction: OUTGOING,
				NetNS:     1,
			},
			intraHost: true,
		},
		{
			name: "local connection, same network ns",
			conn: ConnectionStats{
				Source:    util.AddressFromString("2.2.2.2"),
				Dest:      util.AddressFromString("1.1.1.1"),
				SPort:     80,
				DPort:     12345,
				Direction: INCOMING,
				NetNS:     1,
			},
			intraHost: true,
		},
		{
			name: "local connection, different network ns",
			conn: ConnectionStats{
				Source:    util.AddressFromString("1.1.1.1"),
				Dest:      util.AddressFromString("2.2.2.2"),
				SPort:     12345,
				DPort:     80,
				Direction: OUTGOING,
				NetNS:     1,
			},
			intraHost: true,
		},
		{
			name: "local connection, different network ns",
			conn: ConnectionStats{
				Source:    util.AddressFromString("2.2.2.2"),
				Dest:      util.AddressFromString("1.1.1.1"),
				SPort:     80,
				DPort:     12345,
				Direction: INCOMING,
				NetNS:     2,
			},
			intraHost: true,
		},
		{
			name: "remote connection",
			conn: ConnectionStats{
				Source:    util.AddressFromString("1.1.1.1"),
				Dest:      util.AddressFromString("3.3.3.3"),
				SPort:     12345,
				DPort:     80,
				Direction: OUTGOING,
				NetNS:     1,
			},
			intraHost: false,
		},
	}

	var conns []ConnectionStats
	for _, te := range tests {
		conns = append(conns, te.conn)
	}
	state := newDefaultState().(*networkState)
	state.determineConnectionIntraHost(conns)
	for i, te := range tests {
		assert.Equal(t, te.intraHost, conns[i].IntraHost, "name: %s, conn: %+v", te.name, conns[i])
		if conns[i].Direction == INCOMING {
			if conns[i].IntraHost {
				assert.Nil(t, conns[i].IPTranslation, "name: %s, conn: %+v", te.name, conns[i])
			} else {
				assert.NotNil(t, conns[i].IPTranslation, "name: %s, conn: %+v", te.name, conns[i])
			}
		}
	}
}

func generateRandConnections(n int) []ConnectionStats {
	cs := make([]ConnectionStats, 0, n)
	for i := 0; i < n; i++ {
		cs = append(cs, ConnectionStats{
			Pid:                  123,
			Type:                 TCP,
			Family:               AFINET,
			Source:               util.AddressFromString("127.0.0.1"),
			Dest:                 util.AddressFromString("127.0.0.1"),
			SPort:                uint16(rand.Intn(math.MaxUint16)),
			DPort:                uint16(rand.Intn(math.MaxUint16)),
			MonotonicRecvBytes:   rand.Uint64(),
			MonotonicSentBytes:   rand.Uint64(),
			MonotonicRetransmits: rand.Uint32(),
		})
	}
	return cs
}

var latestTime uint64

func latestEpochTime() uint64 {
	return atomic.AddUint64(&latestTime, 1)
}

func newDefaultState() State {
	// Using values from ebpf.NewConfig()
	return NewState(2*time.Minute, 50000, 75000, 75000, 7500, false)
}
