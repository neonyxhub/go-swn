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
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	neo_swn "go.neonyx.io/go-swn/internal/swn"
	api "go.neonyx.io/go-swn/pkg/bus/pb"
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

func TestEventToLocalListener(t *testing.T) {
	swn, err := newSWN(1)
	defer closeSWN(t, swn)
	require.NoError(t, err)

	swn.GrpcServer.Bus.HasListener = true
	done := make(chan bool, 1)
	completed := 0
	var mu sync.Mutex

	// imitating LocalFunnelEvents rpc
	var wg1 sync.WaitGroup
	wg1.Add(5)

	go func(done chan bool, wg *sync.WaitGroup) {
		for {
			select {
			case <-done:
				return
			case event := <-swn.GrpcServer.Bus.EventToLocal:
				require.True(t, strings.HasPrefix(event.Lexicon.Uri, "uri-"))
				mu.Lock()
				completed += 1
				mu.Unlock()
				wg.Done()
			}
		}
	}(done, &wg1)

	var wg2 sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg2.Add(1)
		go func(count int, wg *sync.WaitGroup) {
			defer wg.Done()
			resp, _, _ := mockEvent(count)
			err := swn.EventToLocalListener(resp)
			require.NoError(t, err)
		}(i, &wg2)
	}

	wg1.Wait()
	wg2.Wait()

	mu.Lock()
	require.Equal(t, completed, 5)
	mu.Unlock()
	done <- true

	// error case
	swn.GrpcServer.Bus.HasListener = false
	resp, _, _ := mockEvent(1)

	go func() {
		<-swn.GrpcServer.Bus.EventToLocal
	}()

	err = swn.EventToLocalListener(resp)
	require.Error(t, err, neo_swn.ErrNoLocalListener)
}

func TestStartEventListening(t *testing.T) {
	swn, err := newSWN(1)
	require.NoError(t, err)
	defer closeSWN(t, swn)

	evt, _, _ := mockEvent(1)
	swn.GrpcServer.Bus.EventFromLocal <- evt
	require.True(t, true, "event should be received via StartEventListening()")
}

func TestStopEventListening(t *testing.T) {
	swn, err := newSWN(1)
	require.NoError(t, err)
	defer closeSWN(t, swn)

	done := make(chan bool, 1)

	go func(done chan bool) {
		evt, _, _ := mockEvent(1)
		swn.GrpcServer.Bus.EventFromLocal <- evt
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

	for _, p := range getter.Peer.Host.Mux().Protocols() {
		log.Printf("protocol: %v\b", p)
	}

	err = sender.PassEventToNetwork(evt)
	require.NoError(t, err)

	evt2 := <-getter.GrpcServer.Bus.EventToLocal
	require.True(t, proto.Equal(evt, evt2))

	// invalid NewMultiaddrBytes()
	log.Println("checking NewMultiaddrBytes")

	maStr := "/ip4/0.0.0.0/tcp/65002/p2p/12D3KooWMGh4n7ra4WiQwCFrr3wCnreYV5KFePs3WBj65NHd4jfo"
	ma, err = multiaddr.NewMultiaddr(maStr)
	require.NoError(t, err)

	evt, _, _ = mockEvent(1)
	evt.Dest.Addr = []byte{}
	err = sender.PassEventToNetwork(evt)
	require.Error(t, err, "empty multiaddr")

	// invalid AddrInfoFromP2pAddr()
	evt, _, _ = mockEvent(1)
	evt.Dest.Addr = []byte{0xbe, 0xef}
	err = sender.PassEventToNetwork(evt)
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

	sender.Peer.EstablishConn(getter.Peer.Getp2pMA())

	conns := sender.Peer.Host.Network().ConnsToPeer(getter.ID())

	evt, _, _ := mockEvent(1)

	err = sender.ConnPassEvent(context.Background(), evt, conns[0])
	require.NoError(t, err)

	select {
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout: no event to getter is sent within 2 sec")
	case evt2 := <-getter.GrpcServer.Bus.EventToLocal:
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

	sender1.Peer.EstablishConn(getter.Peer.Getp2pMA())
	sender2.Peer.EstablishConn(getter.Peer.Getp2pMA())

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
			case evt := <-getter.GrpcServer.Bus.EventToLocal:
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
