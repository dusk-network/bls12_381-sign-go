package bls

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"time"

	bls12381svc "github.com/dusk-network/bls12_381-sign-go"
)

func SwitchToCgo() {
	// ipc.disconnect()
	// GenerateKeys = CgoGenerateKeys
	// Sign = CgoSign
	// Verify = CgoVerify
	// CreateApk = CgoCreateApk
	// AggregatePk = CgoAggregatePk
	// AggregateSig = CgoAggregateSig
}

func SwitchToIPC() {
	// ipc.connect()
	// GenerateKeys = IPCGenerateKeys
	// Sign = IPCSign
	// Verify = IPCVerify
	// CreateApk = IPCCreateApk
	// AggregatePk = IPCAggregatePk
	// AggregateSig = IPCAggregateSig
}

type ipcState struct {
	connected bool
	cmd       *exec.Cmd
}

const (
	ipcPath       = "/tmp/bls12381svc.sock"
	ipcSvcBinPath = "/tmp/bls12381svc"
)

var ipc = new(ipcState)

func (s *ipcState) connect() {
	if s.connected {
		return
	}
	if _, err := os.Stat(ipcSvcBinPath); os.IsNotExist(err) {
		// write the IPC service binary to disk
		fmt.Fprintln(os.Stderr, "writing service binary to", ipcSvcBinPath,
			"...",
		)
		if err := ioutil.WriteFile(
			ipcSvcBinPath, bls12381svc.Binary, 0o700,
		); err != nil {
			panic(err) // not sure what better to do just yet
		}
	}
	// spawn the IPC service
	s.cmd = exec.Command(ipcSvcBinPath)
	s.cmd.Stdout = os.Stdout
	s.cmd.Stdin = os.Stdin
	s.cmd.Stderr = os.Stderr
	if err := s.cmd.Start(); err != nil {
		panic(err)
	}
	// connect the IPC

	s.connected = true
	time.Sleep(time.Second * 2)
}

func (s *ipcState) disconnect() {
	if !s.connected {
		return
	}
	// disconnect the IPC

	// stop the IPC service
	if err := s.cmd.Process.Signal(syscall.SIGINT); err != nil {
		panic(err)
	}
	fmt.Fprintln(os.Stderr, "stopped process", ipcSvcBinPath)

	// remove the socket file
	if err := os.Remove(ipcPath); err != nil {
		// panic(err)
		fmt.Println(err)
	}
	fmt.Fprintln(os.Stderr, "removed socket", ipcPath)

	// // remove the service binary
	// if err := os.Remove(ipcSvcBinPath); err != nil {
	// 	// panic(err)
	// 	fmt.Println(err)
	// }
}

var IPCGenerateKeys = func() (secret []byte, public []byte) {
	return
}

var IPCSign = func(sk, pk, msg []byte) (
	signature []byte,
	err error,
) {
	return
}

var IPCVerify = func(apk, sig, msg []byte) (err error) {
	return
}

var IPCCreateApk = func(pk []byte) (apk []byte, err error) {
	return
}

var IPCAggregatePk = func(apk []byte, pks ...[]byte) (
	newApk []byte,
	err error,
) {
	return
}

var IPCAggregateSig = func(sig []byte, sigs ...[]byte) (
	aggregatedSig []byte,
	err error,
) {
	return
}
