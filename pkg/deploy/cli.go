package deploy

import (
	"os/exec"
	"fmt"
	"bytes"
	"time"
	"os"
	"syscall"
	"log"
	//"regexp"
)

func RunCliMode(inventory string, timeout int) {
	log.Println("Run in Cli mode")
	if inventory == "" {
		log.Println("Please specify inventory file!")
		os.Exit(1)
	}

	if _, err := os.Stat(inventory); os.IsNotExist(err) {
		log.Fatal("Cannot find deployment inventory file: %s!\n", inventory)
	}

	log.Println("Deploy in cli mode with deployment inventory file:", inventory)
	cmdStdout, cmdStderr, exitCode := RunAnsiblePlaybook(inventory, "setup_kube.yaml")
	if exitCode != 0 {
		if len(cmdStderr) > 0 {
			fmt.Printf(cmdStderr)
		}
		os.Exit(1)
	}
	log.Println(cmdStdout)
	//r, _ := regexp.Compile("kubeadm join.*")
	/**
	nodeJoinCmd = r.FindString(cmdStdout)
	if nodeJoinCmd == "" {
		log.Println("no kubeadm join command")
		return "", false
	}
	**/

}

func RunCommand(realCmd *exec.Cmd) (cmdStdout string, cmdStderr string, exitCode int) {
	bytesStdout := &bytes.Buffer{}
	bytesStderr := &bytes.Buffer{}

	realCmd.Stdout = bytesStdout
	realCmd.Stderr = bytesStderr

	// Start command asynchronously
	if err := realCmd.Start(); err != nil {
		log.Println(err)
		return "", "", 1
	}

	// Create a ticker that outputs elapsed time
	ticker := time.NewTicker(time.Second)
	go func(ticker *time.Ticker) {
		now := time.Now()
		for _ = range ticker.C {
			fmt.Printf("Elapsed time: %s\nOutput: %s\n", time.Since(now), string(bytesStdout.Bytes()))
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

	cmdStdout = string(bytesStdout.Bytes())
	cmdStderr = string(bytesStderr.Bytes())
	if waitStatus, ok := realCmd.ProcessState.Sys().(syscall.WaitStatus); !ok {
		exitCode = 1
	}else {
		exitCode = waitStatus.ExitStatus()
	}
	return
}

//RunAnsiblePlaybook
//ansible-playbook -i '<host>:<port>,' <playbook> -e "ansible_ssh_user=<user> ansible_ssh_pass=<pass>"
func RunAnsiblePlaybook(inventory string, playbook string)  (cmdStdout string, cmdStderr string, exitCode int){
	var cmd *exec.Cmd

	cmd = exec.Command("ansible-playbook", "-i", inventory, playbook)

	return RunCommand(cmd)
}

