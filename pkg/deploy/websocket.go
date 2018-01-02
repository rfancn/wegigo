package deploy

import (
	"time"
	"log"
	"github.com/gorilla/websocket"
	"io"
	"bufio"
	"os/exec"
	//"os"
	"net/http"
	"os"
	//"fmt"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 5 * time.Second

	// Time to wait before force close on connection.
	closeGracePeriod = 5 * time.Second
)

func newWebsocket(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	var wsUpgrader = websocket.Upgrader{
		WriteBufferSize: 1024,
	}
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Errot new websocket:", err)
		return nil
	}
	return ws
}

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
				close(done)
			}
		case <-done:
			log.Println("ping done received")
			return
		}
	}
}

func pumpOutput(ws *websocket.Conn, r io.Reader, done chan struct{}) {
	continueScan := true
	s := bufio.NewScanner(r)

	for continueScan {
		//check Scan return value to decide if continue to scan or not
		continueScan = s.Scan()

		select {
		//delay to be 100ms, so it will not hog browser
		case <-time.After(time.Millisecond*100):
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
				log.Println("write message err:", err)
				continueScan = false
				break
			}
		}
	}
	if s.Err() != nil {
		log.Println("scan:", s.Err())
	}
	close(done)

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}

//runPipecommand copied from: https://github.com/gorilla/websocket/blob/master/examples/command/main.go
func runServerCommand(ws *websocket.Conn, cmdName string, cmdArgs ...string) {
	defer ws.Close()

	//use io.pipe to elimate the buffer io
	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()
	defer pipeWriter.Close()

	//redirect command stdout and stderr writes to pipeWriter
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdout = pipeWriter
	cmd.Stderr = pipeWriter

	if err := cmd.Start(); err != nil {
		log.Println("start", err)
		return
	}

	//create a routine to read from pipeReader
	stdoutDone := make(chan struct{})
	go pumpOutput(ws, pipeReader, stdoutDone)

	<-stdoutDone
	log.Println("command output ended")
	//try kill the command is not any A bigger bonk on the head.
	if err := cmd.Process.Signal(os.Kill); err != nil {
		log.Println("term:", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Println("wait:", err)
	}
	log.Println("Exit run command")
}

