package bls

import (
	"testing"
)

func TestSwitchToIPC(t *testing.T) {
	SwitchToIPC()
}

func TestSwitchToCgo(t *testing.T) {
	SwitchToCgo()
}

func TestIPC(t *testing.T) {
	SwitchToIPC()
	defer SwitchToCgo()
	TestGenerateKeys(t)
	TestSignVerify(t)
	TestVerifyWrongKey(t)
	TestAggregation(t)
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
