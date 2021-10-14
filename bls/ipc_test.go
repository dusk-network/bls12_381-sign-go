package bls

import (
	"crypto/rand"
	"testing"
)

// func TestSwitchToIPC(t *testing.T) {
// 	SwitchToIPC()
// 	SwitchToCgo()
// }

func TestIPC(t *testing.T) {
	SwitchToIPC()
	defer SwitchToCgo()
	TestGenerateKeys(t)
	TestSignVerify(t)
	TestVerifyWrongKey(t)
	//TestAggregation(t)
}

func TestAggregationIPC(t *testing.T) {
	SwitchToIPC()
	defer SwitchToCgo()
	sk, pk := GenerateKeys()
	msg := make([]byte, 100)
	rand.Read(msg)
	sig, err := Sign(sk, pk, msg)
	if err != nil {
		t.Fatal(err)
	}

	apk, err := CreateApk(pk)
	if err != nil {
		t.Fatal(err)
	}

	// Aggregating pk
	sk2, pk2 := GenerateKeys()
	sk3, pk3 := GenerateKeys()
	apk2, err := AggregatePk(apk, pk2, pk3)
	eprintln("aggregated pks")
	eprintln(apk)
	eprintln(pk)
	eprintln(pk2)
	eprintln(pk3)
	eprintln("new aggregated key from 3 above", apk2, err)
	if err != nil {
		t.Fatal(err)
	}

	// Aggregating sigs
	sig2, err := Sign(sk2, pk2, msg)
	if err != nil {
		t.Fatal(err)
	}

	sig3, err := Sign(sk3, pk3, msg)
	if err != nil {
		t.Fatal(err)
	}

	aggSig, err := AggregateSig(sig, sig2, sig3)
	eprintln("aggregated signature", aggSig)
	if err != nil {
		t.Fatal(err)
	}

	// Aggregated verification
	eprintln("verifying", apk2, aggSig, msg)
	err = Verify(apk2, aggSig, msg)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkSignIPC(b *testing.B) {
	SwitchToIPC()
	defer SwitchToCgo()
	BenchmarkSign(b)
}

func BenchmarkVerifyIPC(b *testing.B) {
	SwitchToIPC()
	defer SwitchToCgo()
	BenchmarkVerify(b)
}

func BenchmarkAggregatePkIPC(b *testing.B) {
	SwitchToIPC()
	defer SwitchToCgo()
	BenchmarkAggregatePk(b)
}

func BenchmarkAggregateSigIPC(b *testing.B) {
	SwitchToIPC()
	defer SwitchToCgo()
	BenchmarkAggregateSig(b)
}
