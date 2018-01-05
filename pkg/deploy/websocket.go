package deploy

import (
	"time"
	"log"
	"github.com/gorilla/websocket"
	"io"
	"bufio"
	"os/exec"
	"net/http"
	"syscall"
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

func routinePing(ws *websocket.Conn, ch chan int) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	PINGLOOP:
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
				break PINGLOOP
			}
			log.Println("ping success")
		//if we received signal from parent, then we need break the loop
		case <-ch:
			log.Println("ping routine receive exit signal")
			break PINGLOOP
		}
	}
	//close the channel
	close(ch)
}

func routinePumpOutput(ws *websocket.Conn, r io.Reader, ch chan struct{}) {
	s := bufio.NewScanner(r)

	SCANLOOP:
	for s.Scan() {
		select {
		//delay 100ms to write response to browser, so it will not hog browser
		case <-time.After(time.Millisecond*100):
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
				log.Println("pump output routine write message err:", err)
				break SCANLOOP
			}
		}
	}
	if s.Err() != nil {
		log.Println("scan:", s.Err())
	}

	//close channel
	close(ch)
}

func closeWebsocket(ws *websocket.Conn) {
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
}

func killProessGroup(cmd *exec.Cmd) {
	//get progress group id
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		log.Fatal("Error get command process group id")
	}
	log.Println("Found pgid:", pgid)
	//try send SIGKILL signal to whole process group
	if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
		log.Println("term:", err)
	}
}

//runPipecommand copied from: https://github.com/gorilla/websocket/blob/master/examples/command/main.go
func runServerCommand(ws *websocket.Conn, cmdName string, cmdArgs ...string) {
	//use io.pipe to elimate the buffer io
	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()
	defer pipeWriter.Close()

	//redirect command stdout and stderr writes to pipeWriter
	cmd := exec.Command(cmdName, cmdArgs...)
	//create a new process group
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Stdout = pipeWriter
	cmd.Stderr = pipeWriter

	if err := cmd.Start(); err != nil {
		log.Println("start", err)
		return
	}

	chPing := make(chan int)
	chOutput := make(chan struct{})
	//create a routine to read from pipeReader
	go routinePumpOutput(ws, pipeReader, chOutput)
	//create a routine to check if client side still alive or not
	go routinePing(ws, chPing)

	select {
	case <-chPing:
		log.Println("ping failed")
		killProessGroup(cmd)
	case <-chOutput:
		log.Println("output done")
		killProessGroup(cmd)
		//notify ping routine to exit
		chPing<-1
	}

	if err := cmd.Wait(); err != nil {
		log.Println("wait:", err)
	}

	log.Println("Exit run command")
}

