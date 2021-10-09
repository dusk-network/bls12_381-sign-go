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
	TestGenerateKeys(t)
	TestSignVerify(t)
	TestVerifyWrongKey(t)
	TestAggregation(t)
}

func BenchmarkSignIPC(b *testing.B) {
	SwitchToIPC()
	BenchmarkSign(b)
}

func BenchmarkVerifyIPC(b *testing.B) {
	SwitchToIPC()
	BenchmarkVerify(b)
}

func BenchmarkAggregatePkIPC(b *testing.B) {
	SwitchToIPC()
	BenchmarkAggregatePk(b)
}

func BenchmarkAggregateSigIPC(b *testing.B) {
	SwitchToIPC()
	BenchmarkAggregateSig(b)
}
