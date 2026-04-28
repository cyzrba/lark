package main

import (
	"fmt"
	"net"
	"testing"
	"time"

	pb "lark/apps/ws/proto"

	"github.com/gorilla/websocket"
)

func TestWSIdleClientKickedAfterThirtySeconds(t *testing.T) {
	port := freePort(t)
	srv := startWSServer(t, port)
	defer srv.stop(t)

	conn := dialWS(t, port, "idle-heartbeat-client")
	defer conn.Close()

	time.Sleep(31 * time.Second)

	if err := conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatalf("set read deadline failed: %v", err)
	}

	_, _, err := conn.ReadMessage()
	if err == nil {
		t.Fatal("expected idle client to be kicked after more than 30 seconds")
	}
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		t.Fatalf("expected server to close the idle connection, but it was still open after 31 seconds: %v", err)
	}
}

func TestWSHeartbeatPingKeepsClientAlivePastThirtySeconds(t *testing.T) {
	port := freePort(t)
	srv := startWSServer(t, port)
	defer srv.stop(t)

	conn := dialWS(t, port, "heartbeat-client")
	defer conn.Close()

	sendAndExpectHeartbeatPong(t, conn, "heartbeat-0")

	for i := 1; i <= 3; i++ {
		time.Sleep(10 * time.Second)
		sendAndExpectHeartbeatPong(t, conn, fmt.Sprintf("heartbeat-%d", i))
	}

	time.Sleep(2 * time.Second)
	sendAndExpectHeartbeatPong(t, conn, "heartbeat-after-32s")
}

func sendAndExpectHeartbeatPong(t *testing.T, conn *websocket.Conn, requestID string) {
	t.Helper()

	writeProtoMessage(t, conn, &pb.Message{
		RequestId: requestID,
		Payload: &pb.Message_HeartbeatPing{
			HeartbeatPing: &pb.HeartbeatPing{
				ClientUnixTime: time.Now().UnixMilli(),
			},
		},
	})

	got := readProtoMessage(t, conn)
	pong := got.GetHeartbeatPong()
	if pong == nil {
		t.Fatalf("expected heartbeat pong, got %#v", got.GetPayload())
	}
	if pong.GetAckRequestId() != requestID {
		t.Fatalf("unexpected heartbeat ack request id: got %q want %q", pong.GetAckRequestId(), requestID)
	}
	if pong.GetServerUnixTime() == 0 {
		t.Fatal("expected server_unix_time to be set in heartbeat pong")
	}
}
