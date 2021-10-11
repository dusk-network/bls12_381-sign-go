package bls

import (
	"testing"
)

// func TestSwitchToIPC(t *testing.T) {
// 	SwitchToIPC()
// 	SwitchToCgo()
// }

func TestIPC(t *testing.T) {
	SwitchToIPC()
	TestGenerateKeys(t)
	// TestSignVerify(t)
	// TestVerifyWrongKey(t)
	// TestAggregation(t)
	SwitchToCgo()
}

func BenchmarkSignIPC(b *testing.B) {
	SwitchToIPC()
	BenchmarkSign(b)
	SwitchToCgo()
}

func BenchmarkVerifyIPC(b *testing.B) {
	SwitchToIPC()
	BenchmarkVerify(b)
	SwitchToCgo()
}

func BenchmarkAggregatePkIPC(b *testing.B) {
	SwitchToIPC()
	BenchmarkAggregatePk(b)
	SwitchToCgo()
}

func BenchmarkAggregateSigIPC(b *testing.B) {
	SwitchToIPC()
	BenchmarkAggregateSig(b)
	SwitchToCgo()
}
