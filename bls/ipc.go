package bls

func SwitchToCgo() {
	ipc.disconnect()
	GenerateKeys = CgoGenerateKeys
	Sign = CgoSign
	Verify = CgoVerify
	CreateApk = CgoCreateApk
	AggregatePk = CgoAggregatePk
	AggregateSig = CgoAggregateSig
}

func SwitchToIPC() {
	ipc.connect()
	GenerateKeys = IPCGenerateKeys
	Sign = IPCSign
	Verify = IPCVerify
	CreateApk = IPCCreateApk
	AggregatePk = IPCAggregatePk
	AggregateSig = IPCAggregateSig
}

type ipcState struct{}

const ipcPath = "/tmp/bls12381svc.sock"

var ipc = new(ipcState)

func (s *ipcState) connect() {
	// connect the IPC
}

func (s *ipcState) disconnect() {
	// disconnect the IPC
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
