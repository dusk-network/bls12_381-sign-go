package bls

// #cgo darwin LDFLAGS: -L../ -ldusk_bls12_381_sign_darwin -ldl -lm
// #cgo linux LDFLAGS: -L../ -ldusk_bls12_381_sign_linux -ldl -lm
// #include "../libbls.h"
import "C"
import (
	"errors"
	"unsafe"
)

func GenerateKeys() ([]byte, []byte) {
	skBuf := make([]byte, C.SK_SIZE)
	pkBuf := make([]byte, C.PK_SIZE)

	sk_ptr := toPtr(skBuf)
	pk_ptr := toPtr(pkBuf)

	C.generate_keys(sk_ptr, pk_ptr)
	return skBuf, pkBuf
}

func Sign(sk, pk, msg []byte) ([]byte, error) {
	sk_ptr := toPtr(sk)
	pk_ptr := toPtr(pk)
	msg_ptr := toPtr(msg)
	sigBuf := make([]byte, C.SIG_SIZE)
	sig_ptr := toPtr(sigBuf)
	code := C.sign(sk_ptr, pk_ptr, msg_ptr, C.size_t(len(msg)), sig_ptr)
	return sigBuf, formatErr(code)
}

func Verify(apk, sig, msg []byte) error {
	apk_ptr := toPtr(apk)
	sig_ptr := toPtr(sig)
	msg_ptr := toPtr(msg)
	code := C.verify(apk_ptr, sig_ptr, msg_ptr, C.size_t(len(msg)))
	return formatErr(code)
}

func CreateApk(pk []byte) ([]byte, error) {
	pk_ptr := toPtr(pk)
	apkBuf := make([]byte, C.PK_SIZE)
	apk_ptr := toPtr(apkBuf)
	code := C.create_apk(pk_ptr, apk_ptr)
	return apkBuf, formatErr(code)
}

func AggregatePk(apk []byte, pks ...[]byte) ([]byte, error) {
	apk_ptr := toPtr(apk)
	pkBytes := make([]byte, 0)
	for _, pk := range pks {
		pkBytes = append(pkBytes, pk...)
	}

	pk_ptr := toPtr(pkBytes)
	retBuf := make([]byte, C.PK_SIZE)
	ret_ptr := toPtr(retBuf)
	code := C.aggregate_pk(apk_ptr, pk_ptr, C.size_t(len(pkBytes)), ret_ptr)
	return retBuf, formatErr(code)
}

func AggregateSig(sig []byte, sigs ...[]byte) ([]byte, error) {
	sig_ptr := toPtr(sig)
	sigBytes := make([]byte, 0)
	for _, sig := range sigs {
		sigBytes = append(sigBytes, sig...)
	}

	sigs_ptr := toPtr(sigBytes)
	retBuf := make([]byte, C.SIG_SIZE)
	ret_ptr := toPtr(retBuf)
	code := C.aggregate_sig(sig_ptr, sigs_ptr, C.size_t(len(sigBytes)), ret_ptr)
	return retBuf, formatErr(code)
}

func toPtr(data []byte) *C.uchar {
	return (*C.uchar)(unsafe.Pointer(&data[0]))
}

func formatErr(errCode C.bls_sign_status) error {
	switch errCode {
	case C.BLS_OK:
		return nil
	case C.BLS_INVALID_BYTES:
		return errors.New("invalid bytes provided")
	case C.BLS_VERIFICATION_FAILED:
		return errors.New("verification failed")
	default:
		return errors.New("encountered unknown exit code")
	}
}
