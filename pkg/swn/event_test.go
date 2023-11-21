package swn_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	api "go.neonyx.io/go-swn/pkg/bus/pb"
	neo_swn "go.neonyx.io/go-swn/pkg/swn"

	"go.neonyx.io/go-swn/internal/grpcserver"
)

func mockEvent(i int) (*api.Event, []byte, error) {
	evt := &api.Event{
		Dest: &api.Destination{
			Addr: []byte(fmt.Sprintf("addr-%v", i)),
		},
		Lexicon: &api.LexiconUri{
			Uri: fmt.Sprintf("uri-%v", i),
		},
		Data: []byte(fmt.Sprintf("data-%v", i)),
	}

	rawEvt, err := proto.Marshal(evt)
	return evt, rawEvt, err
}

func TestProduceUpstream(t *testing.T) {
	sender, err := newSWN(1)
	defer closeSWN(t, sender)
	require.NoError(t, err)

	done := make(chan bool, 1)
	completed := 0
	var mu sync.Mutex

	N := 3

	// imitating LocalFunnelEvents rpc
	var wg1 sync.WaitGroup
	wg1.Add(N)

	// single consumer to receive events
	go func(done chan bool, wg *sync.WaitGroup) {
		for {
			select {
			case <-done:
				return
			case event := <-sender.EventIO.Upstream:
				require.True(t, strings.HasPrefix(event.Lexicon.Uri, "uri-"))
				mu.Lock()
				completed += 1
				mu.Unlock()
				wg.Done()
			}
		}
	}(done, &wg1)

	var wg2 sync.WaitGroup
	wg2.Add(N)
	for i := 0; i < N; i++ {
		go func(count int, wg *sync.WaitGroup) {
			defer wg.Done()
			resp, _, _ := mockEvent(count)
			err := sender.EventBus.SendUpstream(resp)
			require.NoError(t, err)
		}(i, &wg2)
	}

	wg1.Wait()
	wg2.Wait()

	mu.Lock()
	require.Equal(t, completed, N)
	mu.Unlock()
	done <- true

	// error case: no one listens to EventUpstream
	for i := 0; i < 2; i++ {
		resp, _, _ := mockEvent(1)
		err = sender.EventBus.SendUpstream(resp)
		require.Error(t, err, grpcserver.ErrNoLocalListener)
	}

	// 1 buffered upstream event should be flushed
	require.Equal(t, sender.EventIO.UpstreamBufCnt, 2)

	// bring listener back, and flush events
	go func() {
		<-sender.EventIO.Upstream
	}()
}

func TestStartEventListening(t *testing.T) {
	swn, err := newSWN(1)
	require.NoError(t, err)
	defer closeSWN(t, swn)

	evt, _, _ := mockEvent(1)
	swn.EventIO.Downstream <- evt
	require.True(t, true, "event should be received via StartEventListening()")
}

func TestStopEventListening(t *testing.T) {
	swn, err := newSWN(1)
	require.NoError(t, err)
	defer closeSWN(t, swn)

	done := make(chan bool, 1)

	go func(done chan bool) {
		evt, _, _ := mockEvent(1)
		swn.EventIO.Downstream <- evt
		done <- true
	}(done)

	swn.StopEventListening()

	select {
	case <-time.After(100 * time.Millisecond):
		require.True(t, true, "should timeout as no event listener is stopped")
	case <-done:
		t.Fatal("should not reached here")
	}
}

func TestPassEventToNetwork(t *testing.T) {
	getter, err := newSWN(1)
	require.NoError(t, err)
	defer closeSWN(t, getter)

	sender, err := newSWN(2)
	require.NoError(t, err)
	defer closeSWN(t, sender)

	evt, _, _ := mockEvent(1)
	ma := getter.Peer.Getp2pMA()
	evt.Dest.Addr = ma.Bytes()

	err = sender.PassEventToNetwork(evt)
	require.NoError(t, err)

	log.Println("waiting for event come to EventUpstream")
	evt2 := <-getter.EventIO.Upstream
	require.True(t, proto.Equal(evt, evt2))

	// invalid NewMultiaddrBytes()
	evt, _, _ = mockEvent(1)
	evt.Dest.Addr = []byte{}
	err = sender.PassEventToNetwork(evt)
	require.Error(t, err, "empty multiaddr")

	// invalid AddrInfoFromP2pAddr()
	evt, _, _ = mockEvent(1)
	evt.Dest.Addr = []byte{0xbe, 0xef}
	err = sender.PassEventToNetwork(evt)
	require.NoError(t, err)
	require.Error(t, peer.ErrInvalidAddr)

	// ErrNoExistingConnection
	// TODO: PassEventToNetwork() always connects to given multiaddr in Event
}

func TestConnPassEvent(t *testing.T) {
	getter, err := newSWN(1)
	require.NoError(t, err)
	defer closeSWN(t, getter)

	sender, err := newSWN(2)
	require.NoError(t, err)
	defer closeSWN(t, sender)

	err = sender.Peer.EstablishConn(context.Background(), getter.Peer.Getp2pMA())
	require.NoError(t, err)

	conns := sender.Peer.Host.Network().ConnsToPeer(getter.ID())
	require.Equal(t, 1, len(conns))

	// manually authorize
	getter.AuthDeviceMap[sender.ID().String()] = sender.Device.Id

	evt, _, _ := mockEvent(1)

	err = sender.ConnPassEvent(context.Background(), evt, conns[0])
	require.NoError(t, err)

	select {
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout: no event to getter is sent within 2 sec")
	case evt2 := <-getter.EventIO.Upstream:
		require.True(t, proto.Equal(evt, evt2))
	}
}

func TestMultipleSenders(t *testing.T) {
	getter, err := newSWN(1)
	require.NoError(t, err)
	defer closeSWN(t, getter)

	sender1, err := newSWN(2)
	require.NoError(t, err)
	defer closeSWN(t, sender1)

	sender2, err := newSWN(3)
	require.NoError(t, err)
	defer closeSWN(t, sender2)

	err = sender1.Peer.EstablishConn(context.Background(), getter.Peer.Getp2pMA())
	require.NoError(t, err)
	err = sender2.Peer.EstablishConn(context.Background(), getter.Peer.Getp2pMA())
	require.NoError(t, err)

	// manually authorize
	getter.AuthDeviceMap[sender1.ID().String()] = sender1.Device.Id
	getter.AuthDeviceMap[sender2.ID().String()] = sender2.Device.Id

	senders := []*neo_swn.SWN{sender1, sender2}
	done := make(chan bool, 1)
	completed := []time.Time{}
	var mu sync.Mutex
	var wg1 sync.WaitGroup
	wg1.Add(2)

	go func(done chan bool, wg *sync.WaitGroup) {
		for {
			select {
			case <-done:
				return
			case evt := <-getter.EventIO.Upstream:
				require.NotEmpty(t, evt)
				mu.Lock()
				completed = append(completed, time.Now())
				mu.Unlock()
				wg.Done()
			}
		}
	}(done, &wg1)

	var wg2 sync.WaitGroup

	for i, sender := range senders {
		wg2.Add(1)
		// send from each sender simultaneously
		go func(senderId int, sender *neo_swn.SWN, wg *sync.WaitGroup) {
			defer wg.Done()
			conns := sender.Peer.Host.Network().ConnsToPeer(getter.ID())
			require.Equal(t, len(conns), 1)

			// event for getter from senders
			evt, _, _ := mockEvent(senderId)
			ma := getter.Peer.Getp2pMA()
			evt.Dest.Addr = ma.Bytes()

			err = sender.ConnPassEvent(context.Background(), evt, conns[0])
			require.NoError(t, err)
		}(i, sender, &wg2)
	}

	wg1.Wait()
	wg2.Wait()

	done <- true
	require.Equal(t, len(completed), 2)
	diff := completed[0].Sub(completed[1])
	require.Equal(t, diff.Milliseconds(), int64(0))
}
