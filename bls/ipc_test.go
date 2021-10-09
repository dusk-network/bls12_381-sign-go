package bls

import (
	"testing"
)

func TestSwitchToIPC(t *testing.T) {
	ipc.connect()
}

func TestSwitchToCgo(t *testing.T) {
	ipc.disconnect()
}

//
// func TestGenerateKeysIPC(t *testing.T) {
// 	SwitchToIPC()
// 	TestGenerateKeysIPC(t)
// }
//
// func TestSignVerifyIPC(t *testing.T) {
// 	SwitchToIPC()
// 	TestSignVerifyIPC(t)
// }
//
// func TestVerifyWrongKeyIPC(t *testing.T) {
// 	SwitchToIPC()
// 	TestVerifyWrongKeyIPC(t)
// }
//
// func TestAggregationIPC(t *testing.T) {
// 	SwitchToIPC()
// 	TestAggregationIPC(t)
// }
//
// func BenchmarkSignIPC(b *testing.B) {
// 	SwitchToIPC()
// 	BenchmarkSignIPC(b)
// }
//
// func BenchmarkVerifyIPC(b *testing.B) {
// 	SwitchToIPC()
// 	BenchmarkVerifyIPC(b)
// }
//
// func BenchmarkAggregatePkIPC(b *testing.B) {
// 	SwitchToIPC()
// 	BenchmarkAggregatePkIPC(b)
// }
//
// func BenchmarkAggregateSigIPC(b *testing.B) {
// 	SwitchToIPC()
// 	BenchmarkAggregateSigIPC(b)
// }
