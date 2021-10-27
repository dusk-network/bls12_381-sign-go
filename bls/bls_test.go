package bls_test

import (
	"crypto/rand"
	"testing"

	"github.com/dusk-network/bls12_381-sign-go/bls"
	"github.com/stretchr/testify/assert"
)

func TestGenerateKeys(t *testing.T) {
	sk, pk := bls.GenerateKeys()
	assert.NotEqual(t, sk, make([]byte, 32))
	assert.NotEqual(t, pk, make([]byte, 96))
}

func TestSignVerify(t *testing.T) {
	sk, pk := bls.GenerateKeys()
	msg := make([]byte, 100)
	rand.Read(msg)
	sig, err := bls.Sign(sk, pk, msg)
	if err != nil {
		t.Fatal(err)
	}

	apk, err := bls.CreateApk(pk)
	if err != nil {
		t.Fatal(err)
	}

	err = bls.Verify(apk, sig, msg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVerifyWrongKey(t *testing.T) {
	sk, pk := bls.GenerateKeys()
	msg := make([]byte, 100)
	rand.Read(msg)
	sig, err := bls.Sign(sk, pk, msg)
	if err != nil {
		t.Fatal(err)
	}

	_, pk = bls.GenerateKeys()

	apk, err := bls.CreateApk(pk)
	if err != nil {
		t.Fatal(err)
	}

	assert.Error(t, bls.Verify(apk, sig, msg))
}

func TestAggregation(t *testing.T) {
	sk, pk := bls.GenerateKeys()
	msg := make([]byte, 100)
	rand.Read(msg)
	sig, err := bls.Sign(sk, pk, msg)
	if err != nil {
		t.Fatal(err)
	}

	apk, err := bls.CreateApk(pk)
	if err != nil {
		t.Fatal(err)
	}

	// Aggregating pk
	sk2, pk2 := bls.GenerateKeys()
	sk3, pk3 := bls.GenerateKeys()
	apk2, err := bls.AggregatePk(apk, pk2, pk3)
	if err != nil {
		t.Fatal(err)
	}

	// Aggregating sigs
	sig2, err := bls.Sign(sk2, pk2, msg)
	if err != nil {
		t.Fatal(err)
	}

	sig3, err := bls.Sign(sk3, pk3, msg)
	if err != nil {
		t.Fatal(err)
	}

	aggSig, err := bls.AggregateSig(sig, sig2, sig3)
	if err != nil {
		t.Fatal(err)
	}

	// Aggregated verification
	err = bls.Verify(apk2, aggSig, msg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAggregationRaw(t *testing.T) {
	sk, pk, pkraw := bls.GenerateKeysWithRaw()
	pkraw_conv, err := bls.PkToRaw(pk)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, pkraw, pkraw_conv)
	msg := make([]byte, 100)
	rand.Read(msg)
	sig, err := bls.Sign(sk, pk, msg)
	if err != nil {
		t.Fatal(err)
	}

	// Aggregating pk
	sk2, pk2, pk2raw := bls.GenerateKeysWithRaw()
	sk3, pk3, pk3raw := bls.GenerateKeysWithRaw()
	apk2, err := bls.AggregatePKsUnchecked(pkraw, pk2raw, pk3raw)
	if err != nil {
		t.Fatal(err)
	}

	// Aggregating sigs
	sig2, err := bls.Sign(sk2, pk2, msg)
	if err != nil {
		t.Fatal(err)
	}

	sig3, err := bls.Sign(sk3, pk3, msg)
	if err != nil {
		t.Fatal(err)
	}

	aggSig, err := bls.AggregateSig(sig, sig2, sig3)
	if err != nil {
		t.Fatal(err)
	}

	// Aggregated verification
	err = bls.Verify(apk2, aggSig, msg)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkSign(b *testing.B) {
	sk, pk := bls.GenerateKeys()
	msg := make([]byte, 100)
	rand.Read(msg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := bls.Sign(sk, pk, msg); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVerify(b *testing.B) {
	sk, pk := bls.GenerateKeys()
	msg := make([]byte, 100)
	rand.Read(msg)

	sig, err := bls.Sign(sk, pk, msg)
	if err != nil {
		b.Fatal(err)
	}

	apk, err := bls.CreateApk(pk)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := bls.Verify(apk, sig, msg); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAggregatePk(b *testing.B) {
	_, pk := bls.GenerateKeys()

	apk, err := bls.CreateApk(pk)
	if err != nil {
		b.Fatal(err)
	}

	pks := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		_, pk := bls.GenerateKeys()
		pks[i] = pk
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bls.AggregatePk(apk, pks[i])
	}
}
func BenchmarkAggregatePks(b *testing.B) {
	_, pk := bls.GenerateKeys()

	apk, err := bls.CreateApk(pk)
	if err != nil {
		b.Fatal(err)
	}

	pks := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		_, pk := bls.GenerateKeys()
		pks[i] = pk
	}

	b.ResetTimer()
	// println("Len pks ", len(pks))
	bls.AggregatePk(apk, pks...)
}

func BenchmarkAggregatePksRaw(b *testing.B) {
	_, _, firstpk := bls.GenerateKeysWithRaw()

	pks := make([][]byte, b.N)
	pks[0] = firstpk

	for i := 1; i < b.N; i++ {
		_, _, rawpk := bls.GenerateKeysWithRaw()
		pks[i] = rawpk
	}
	b.ResetTimer()
	// println("Len pksraw ", len(pks))
	// for i := 0; i < b.N; i++ {
	bls.AggregatePKsUnchecked(pks...)
	// }
}

func TestAggregatePk64(b *testing.T) {
	_, pk := bls.GenerateKeys()

	apk, err := bls.CreateApk(pk)
	if err != nil {
		b.Fatal(err)
	}

	pks := make([][]byte, 63)
	for i := 0; i < 63; i++ {
		_, pk := bls.GenerateKeys()
		pks[i] = pk
	}

	// b.ResetTimer()
	for i := 0; i < 63; i++ {
		bls.AggregatePk(apk, pks[i])
	}
}
func TestAggregatePk64s(b *testing.T) {
	_, pk := bls.GenerateKeys()

	apk, err := bls.CreateApk(pk)
	if err != nil {
		b.Fatal(err)
	}

	pks := make([][]byte, 63)
	for i := 0; i < 63; i++ {
		_, pk := bls.GenerateKeys()
		pks[i] = pk
	}

	bls.AggregatePk(apk, pks...)
}

func TestAggregatePk64sRaw(b *testing.T) {
	_, _, firstpk := bls.GenerateKeysWithRaw()

	pks := make([][]byte, 64)
	pks[0] = firstpk

	for i := 1; i < 64; i++ {
		_, _, rawpk := bls.GenerateKeysWithRaw()
		pks[i] = rawpk
	}
	// b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	bls.AggregatePKsUnchecked(pks...)
	// }
}

func BenchmarkAggregateSig(b *testing.B) {
	sk, pk := bls.GenerateKeys()
	msg := make([]byte, 100)
	rand.Read(msg)

	sig, err := bls.Sign(sk, pk, msg)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := bls.AggregateSig(sig, sig); err != nil {
			b.Fatal(err)
		}
	}
}
