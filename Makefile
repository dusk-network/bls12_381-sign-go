OS := $(shell sh -c 'uname -s 2>/dev/null || echo linux' | tr "[:upper:]" "[:lower:]")

build: 
	rm -rf bls12_381-sign || true
	rm bls12_381-sign-go || true
	git clone https://github.com/dusk-network/bls12_381-sign
	cd bls12_381-sign && cargo build --release
	mv bls12_381-sign/target/release/libdusk_bls12_381_sign.a libdusk_bls12_381_sign_$(OS).a
	rm -rf bls12_381-sign
	go build ./...

test: build
	go test -v ./...

bench: build
	go test -v -bench=. ./...
