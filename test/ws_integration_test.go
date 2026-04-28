package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	pb "lark/apps/ws/proto"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func TestWSConnect(t *testing.T) {
	port := freePort(t)
	srv := startWSServer(t, port)
	defer srv.stop(t)

	conn := dialWS(t, port, "alice")
	defer conn.Close()
}

func TestWSMessageForwardingViaChatRPC(t *testing.T) {
	port := freePort(t)
	srv := startWSServer(t, port)
	defer srv.stop(t)

	bob := dialWS(t, port, "bob")
	defer bob.Close()

	alice := dialWS(t, port, "alice")
	defer alice.Close()

	wantContent := "hello bob"
	outbound := &pb.Message{
		RequestId: "req-test-1",
		Payload: &pb.Message_Chat{
			Chat: &pb.ChatMessage{
				To:      "bob",
				Content: wantContent,
			},
		},
	}

	writeProtoMessage(t, alice, outbound)

	got := readProtoMessage(t, bob)
	chat := got.GetChat()
	if chat == nil {
		t.Fatalf("expected chat message, got %#v", got.GetPayload())
	}
	if chat.GetFrom() != "alice" {
		t.Fatalf("unexpected sender: got %q want %q", chat.GetFrom(), "alice")
	}
	if chat.GetTo() != "bob" {
		t.Fatalf("unexpected receiver: got %q want %q", chat.GetTo(), "bob")
	}
	if chat.GetContent() != wantContent {
		t.Fatalf("unexpected content: got %q want %q", chat.GetContent(), wantContent)
	}
	if got.GetRequestId() != outbound.GetRequestId() {
		t.Fatalf("unexpected request id: got %q want %q", got.GetRequestId(), outbound.GetRequestId())
	}
}

type processHandle struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

type wsServerProcess struct {
	ws   *processHandle
	chat *processHandle
}

func startWSServer(t *testing.T, port int) *wsServerProcess {
	t.Helper()

	root := repoRoot(t)
	chatPort := freePort(t)

	chatServer := startChatRPCServer(t, root, chatPort)
	wsServer := startWSAPIProcess(t, root, port, chatPort)

	return &wsServerProcess{
		ws:   wsServer,
		chat: chatServer,
	}
}

func startChatRPCServer(t *testing.T, root string, port int) *processHandle {
	t.Helper()

	cfgPath := filepath.Join(t.TempDir(), "chat.yaml")
	exePath := filepath.Join(t.TempDir(), "chat-rpc-test-server.exe")
	cfgContent := fmt.Sprintf("Name: chat.rpc\nListenOn: 127.0.0.1:%d\n", port)
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644); err != nil {
		t.Fatalf("write temp chat rpc config failed: %v", err)
	}

	buildGoBinary(t, root, "./apps/chat/rpc", exePath)
	cmd := startProcess(t, root, exePath, "-f", cfgPath)
	waitForTCP(t, port, 10*time.Second, cmd.stdout, cmd.stderr, cmd.cmd, "chat rpc")
	return cmd
}

func startWSAPIProcess(t *testing.T, root string, port, chatPort int) *processHandle {
	t.Helper()

	cfgPath := filepath.Join(t.TempDir(), "ws-api.yaml")
	exePath := filepath.Join(t.TempDir(), "ws-test-server.exe")
	cfgContent := fmt.Sprintf(
		"Name: ws-api\nHost: 127.0.0.1\nPort: %d\nChatRpc:\n  Endpoints:\n    - 127.0.0.1:%d\n",
		port,
		chatPort,
	)
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644); err != nil {
		t.Fatalf("write temp config failed: %v", err)
	}

	buildGoBinary(t, root, "./apps/ws", exePath)
	cmd := startProcess(t, root, exePath, "-f", cfgPath)
	waitForTCP(t, port, 10*time.Second, cmd.stdout, cmd.stderr, cmd.cmd, "ws api")
	return cmd
}

func buildGoBinary(t *testing.T, root, pkg, exePath string) {
	t.Helper()

	build := exec.Command("go", "build", "-o", exePath, pkg)
	build.Dir = root
	build.Env = append(os.Environ(),
		"GOCACHE="+filepath.Join(root, ".gocache"),
	)
	buildOutput, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build %s failed: %v\noutput:\n%s", pkg, err, string(buildOutput))
	}
}

func startProcess(t *testing.T, root, exePath string, args ...string) *processHandle {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, exePath, args...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(),
		"GOCACHE="+filepath.Join(root, ".gocache"),
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		cancel()
		t.Fatalf("start process %s failed: %v", exePath, err)
	}

	return &processHandle{
		cmd:    cmd,
		cancel: cancel,
		stdout: &stdout,
		stderr: &stderr,
	}
}

func (p *wsServerProcess) stop(t *testing.T) {
	t.Helper()

	stopProcess(t, p.ws)
	stopProcess(t, p.chat)
}

func stopProcess(t *testing.T, p *processHandle) {
	t.Helper()

	if p == nil {
		return
	}

	p.cancel()
	if p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
	}

	done := make(chan error, 1)
	go func() {
		done <- p.cmd.Wait()
	}()

	select {
	case <-time.After(3 * time.Second):
		t.Logf("process did not exit in time")
	case <-done:
	}
}

func dialWS(t *testing.T, port int, name string) *websocket.Conn {
	t.Helper()

	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("127.0.0.1:%d", port),
		Path:   "/ws/" + name,
	}

	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		if resp != nil {
			t.Fatalf("dial websocket failed: %v, status=%s", err, resp.Status)
		}
		t.Fatalf("dial websocket failed: %v", err)
	}

	return conn
}

func writeProtoMessage(t *testing.T, conn *websocket.Conn, message *pb.Message) {
	t.Helper()

	data, err := proto.Marshal(message)
	if err != nil {
		t.Fatalf("marshal proto failed: %v", err)
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		t.Fatalf("write websocket message failed: %v", err)
	}
}

func readProtoMessage(t *testing.T, conn *websocket.Conn) *pb.Message {
	t.Helper()

	if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatalf("set read deadline failed: %v", err)
	}

	messageType, payload, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read websocket message failed: %v", err)
	}
	if messageType != websocket.BinaryMessage {
		t.Fatalf("unexpected websocket message type: got %d want %d", messageType, websocket.BinaryMessage)
	}

	var message pb.Message
	if err := proto.Unmarshal(payload, &message); err != nil {
		t.Fatalf("unmarshal proto failed: %v", err)
	}

	return &message
}

func freePort(t *testing.T) int {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("allocate free port failed: %v", err)
	}
	defer ln.Close()

	return ln.Addr().(*net.TCPAddr).Port
}

func waitForTCP(t *testing.T, port int, timeout time.Duration, stdout, stderr *bytes.Buffer, cmd *exec.Cmd, serviceName string) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	address := fmt.Sprintf("127.0.0.1:%d", port)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", address, 300*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return
		}

		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			t.Fatalf("%s exited early\nstdout:\n%s\nstderr:\n%s", serviceName, stdout.String(), stderr.String())
		}

		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("%s was not ready in time\nstdout:\n%s\nstderr:\n%s", serviceName, stdout.String(), stderr.String())
}

func repoRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory failed: %v", err)
	}

	if strings.EqualFold(filepath.Base(dir), "test") {
		return filepath.Dir(dir)
	}

	return dir
}
