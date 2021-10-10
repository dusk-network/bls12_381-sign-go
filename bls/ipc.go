package bls

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/dusk-network/bls12_381-sign-go/bls/proto"
	"google.golang.org/grpc"

	bls12381svc "github.com/dusk-network/bls12_381-sign-go"
)

type Bls12381Sign interface {
	GenerateKeys() (secret []byte, public []byte)
	Sign(sk, pk, msg []byte) (signature []byte, err error)
	Verify(apk, sig, msg []byte) (err error)
	CreateApk(pk []byte) (apk []byte, err error)
	AggregatePk(apk []byte, pks ...[]byte)
	AggregateSig(sig []byte, sigs ...[]byte)
}

func SwitchToCgo() {
	ipc.disconnect()
	GenerateKeys = cgo.GenerateKeys
	Sign = cgo.Sign
	Verify = cgo.Verify
	CreateApk = cgo.CreateApk
	AggregatePk = cgo.AggregatePk
	AggregateSig = cgo.AggregateSig
}

func SwitchToIPC() {
	ipc.connect()
	GenerateKeys = ipc.GenerateKeys
	Sign = ipc.Sign
	Verify = ipc.Verify
	CreateApk = ipc.CreateApk
	AggregatePk = ipc.AggregatePk
	AggregateSig = ipc.AggregateSig
}

type ipcState struct {
	connected bool
	cmd       *exec.Cmd
	proto.SignerClient
	*grpc.ClientConn
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
		eprintln("writing service binary to", ipcSvcBinPath, "...")
		if err := ioutil.WriteFile(
			ipcSvcBinPath, bls12381svc.Binary, 0o700,
		); err != nil {
			panic(err) // not sure what better to do just yet
		}
	}

	// spawn the IPC service
	s.cmd = exec.Command(ipcSvcBinPath)
	// command will print output to parent terminal
	s.cmd.Stdout = os.Stdout
	s.cmd.Stderr = os.Stderr
	if err := s.cmd.Start(); err != nil {
		panic(err)
	}

	// wait a second for the service to come up
	time.Sleep(time.Second)

	// connect the IPC
	dialer := func(ctx context.Context, path string) (net.Conn, error) {
		var d net.Dialer
		d.LocalAddr = nil
		raddr := net.UnixAddr{Name: ipcPath, Net: "unix"}
		conn, err := d.DialContext(ctx, "unix", raddr.String())
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
	var err error
	s.ClientConn, err = grpc.Dial(
		ipcPath,
		grpc.WithInsecure(),
		grpc.WithContextDialer(dialer),
	)
	if err != nil {
		panic(err)
	}
	s.SignerClient = proto.NewSignerClient(s.ClientConn)
	s.connected = true
}

func eprintln(args ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, args...)
}

func (s *ipcState) disconnect() {
	if !s.connected {
		return
	}
	//  mark that we are not connected so nobody tries to use this
	s.connected = false

	// disconnect the IPC
	if err := s.ClientConn.Close(); err != nil {
		eprintln(err)
	} else {
		eprintln("closed client connection")
	}

	// stop the IPC service. The service knows SIGINT means shut down so it will
	// stop and release its resources from this signal
	if err := s.cmd.Process.Signal(syscall.SIGINT); err != nil {
		panic(err)
	} else {
		eprintln("stopped process", ipcSvcBinPath)
	}

	// remove the socket file
	if err := os.Remove(ipcPath); err != nil {
		// panic(err)
		eprintln(err)
	} else {
		eprintln("removed socket", ipcPath)
	}

	// remove the service binary
	if err := os.Remove(ipcSvcBinPath); err != nil {
		eprintln(err)
	}
}

func (s *ipcState) GenerateKeys() (secret []byte, public []byte) {
	if !s.connected {
		eprintln("attempting to call API without being connected")
	}
	keys, err := s.SignerClient.GenerateKeys(context.Background(),
		&proto.GenerateKeysRequest{},
	)
	if err != nil {
		return nil, nil
	}
	return keys.GetSecretKey(), keys.GetPublicKey()
}

func (s *ipcState) Sign(sk, pk, msg []byte) (
	signature []byte,
	err error,
) {
	if !s.connected {
		eprintln("attempting to call API without being connected")
	}
	sig, err := s.SignerClient.Sign(context.Background(),
		&proto.SignRequest{
			SecretKey: sk,
			PublicKey: pk,
			Message:   msg,
		},
	)
	if err != nil {
		return nil, err
	}
	return sig.GetSignature(), nil
}

func (s *ipcState) Verify(apk, sig, msg []byte) (err error) {
	if !s.connected {
		eprintln("attempting to call API without being connected")
	}
	_, err = s.SignerClient.Verify(context.Background(),
		&proto.VerifyRequest{
			Apk:       apk,
			Signature: sig,
			Message:   msg,
		},
	)
	return err
}

func (s *ipcState) CreateApk(pk []byte) (apk []byte, err error) {
	if !s.connected {
		eprintln("attempting to call API without being connected")
	}
	var a *proto.CreateAPKResponse
	a, err = s.SignerClient.CreateAPK(context.Background(),
		&proto.CreateAPKRequest{
			PublicKey: pk,
		},
	)
	return a.GetAPK(), err
}

func (s *ipcState) AggregatePk(apk []byte, pks ...[]byte) (
	newApk []byte,
	err error,
) {
	if !s.connected {
		eprintln("attempting to call API without being connected")
	}
	var a *proto.AggregateResponse
	a, err = s.SignerClient.AggregatePK(context.Background(),
		&proto.AggregatePKRequest{
			APK:  apk,
			Keys: pks,
		},
	)
	return a.GetCode(), err
}

func (s *ipcState) AggregateSig(sig []byte, sigs ...[]byte) (
	aggregatedSig []byte,
	err error,
) {
	if !s.connected {
		eprintln("attempting to call API without being connected")
	}
	var a *proto.AggregateResponse
	a, err = s.SignerClient.AggregateSig(context.Background(),
		&proto.AggregateSigRequest{
			Signature:  sig,
			Signatures: sigs,
		},
	)
	return a.GetCode(), err
}
