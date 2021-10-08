build :
	rm -rf bls12_381-sign || true
	git clone https://github.com/dusk-network/bls12_381-sign
	cd bls12_381-sign && git checkout microservice && cargo build --release
	cp bls12_381-sign/target/release/libdusk_bls12_381_sign.a libdusk_bls12_381_sign.a
	 go build ./...

test: build
	go test -v ./...

bench: build
	go test -v -bench=. ./...

