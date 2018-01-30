package main

import (
	"path/filepath"
	"os"
	"fmt"

	"errors"
	"reflect"
	"net/http"
	"os/exec"
	"io"
	"bytes"
	"log"
	"time"
	//"syscall"
)

type ITest interface {
	ddd()
}

type tt struct {
	name string
}

type ttt struct {
	tt
}

func (t *ttt) ddd() {
	fmt.Println("in ttt's ddd", t.name)
}

func (t *tt) ddd() {
	fmt.Println("in tt's ddd", t.name)
}

func (t *tt) xxx() {
	//fmt.Println("in tt's xxx", t.name)
}

//getBasedir returns wegigo running dir
func getBasedir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Error get current execute path: %s", err)
		return ""
	}

	return dir
}

func test(t ITest) {
	t.ddd()

	pt, ok := t.(*tt)
	if !ok {
		fmt.Println("conversion failed")
		os.Exit(1)
	}
	pt.xxx()
}

func testfunc1(name string) ([]byte, error) {
	return nil, errors.New("testfunc1")
}

func testfunc2(name string) ([]string, error) {
	return nil, errors.New("testfunc2")
}

func testfunc3(name string) (os.FileInfo, error) {
	return nil, errors.New("testfunc2")
}

func getInterfaceType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// the specified command
	cmd := exec.Command("ls", "-ltR", "/")
	// sign: func Pipe() (*PipeReader, *PipeWriter)
	// A PipeReader is the read half of a pipe.
	// A PipeWriter is the write half of a pipe.
	pipeReader, pipeWriter := io.Pipe()
	// the pipe content is streamed to 2 io.Writers
	cmd.Stdout = pipeWriter
	cmd.Stderr = pipeWriter
	go writeCmdOutput(w, pipeReader)
	// Run starts the specified command and waits for it to complete.
	cmd.Run()
	// important to close any opened io.Writer
	pipeWriter.Close()
}

func writeCmdOutput(resp http.ResponseWriter, pipeReader *io.PipeReader) {
	buffer := make([]byte, 1024)
	for {
		n, err := pipeReader.Read(buffer)
		if err != nil {
			pipeReader.Close()
			break
		}
		data := buffer[0:n]
		resp.Write(data)
		// The Flusher interface is implemented by ResponseWriters
		// that allow an HTTP handler to flush buffered data to the client.
		// Its Flush() sends any buffered data to the client.
		if f, ok := resp.(http.Flusher); ok {
			f.Flush()
		}
		// reset buffer
		for i := 0; i < n; i++ {
			buffer[i] = 0
		}
	}
}

func handler1(w http.ResponseWriter, r *http.Request) {
	// the specified command
	cmd := exec.Command("watch", "-n", "1", "ls")

	RunCommand(cmd, w)
}

func RunCommand(realCmd *exec.Cmd, w http.ResponseWriter) {
	bytesStdout := &bytes.Buffer{}
	bytesStderr := &bytes.Buffer{}

	realCmd.Stdout = bytesStdout
	realCmd.Stderr = bytesStderr

	// Start command asynchronously
	if err := realCmd.Start(); err != nil {
		log.Println(err)
		w.Write([]byte("Error"))
	}

	// Create a ticker that outputs elapsed time
	ticker := time.NewTicker(time.Second)
	go func(ticker *time.Ticker) {
		//now := time.Now()
		for _ = range ticker.C {
			//fmt.Printf("Elapsed time: %s\nOutput: %s\n", time.Since(now), string(bytesStdout.Bytes()))
			w.Write(bytesStdout.Bytes())
		}
	}(ticker)

	// Create a timer that will kill the process
	//timer := time.NewTimer(time.Minute * time.Duration(timeout))
	timer := time.NewTimer(time.Minute * time.Duration(30))
	go func(timer *time.Timer, ticker *time.Ticker, cmd *exec.Cmd) {
		for _ = range timer.C {
			err := cmd.Process.Signal(os.Kill)
			log.Println(err)
			ticker.Stop()
		}
	}(timer, ticker, realCmd)

	// Only proceed once the process has finished
	realCmd.Wait()

	/**
	cmdStdout = string(bytesStdout.Bytes())
	cmdStderr = string(bytesStderr.Bytes())
	if waitStatus, ok := realCmd.ProcessState.Sys().(syscall.WaitStatus); !ok {
		exitCode = 1
	}else {
		exitCode = waitStatus.ExitStatus()
	}
	**/
	return
}

func routinePumpOutput(ch chan int) {
	log.Println("Enter output")
	time.Sleep(10 * time.Second)
	close(ch)
	log.Println("Exit output")
}

func routinePing(chPing chan int) {
	log.Println("Enter ping")
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	PINGLOOP:
	for {
		log.Println("in ping for loop")
		select {
		case <-ticker.C:
			log.Println("ticker ping")
		case <-chPing:
			log.Println("ping receive channel closed")
			break PINGLOOP
		}
	}

	close(chPing)
	log.Println("Exit ping")
}

func main() {
	type Foo struct {
		FirstName string `tag_name:"tag 1"`
		LastName  string `tag_name:"tag 2"`
		Age       int    `tag_name:"tag 3"`
	}

	f := &Foo{}

	a := "test"

	t := reflect.TypeOf(f)

	//v := reflect.ValueOf(f).Elem()

	fmt.Printf("%v", t.Elem().Name())
	fmt.Printf("%v", a.Elem().Name())
}

