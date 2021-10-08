buildfirst:
	rm -rf bls12_381-sign || true
	git clone https://github.com/dusk-network/bls12_381-sign
	cd bls12_381-sign && git checkout microservice && cargo build --release
	mv bls12_381-sign/target/release/libdusk_bls12_381_sign.a libdusk_bls12_381_sign.a
	# todo: also move bls12381svc binary, maybe embed?
	go build ./...

build:
	cd bls12_381-sign && git pull && cargo build --release
	mv bls12_381-sign/target/release/libdusk_bls12_381_sign.a libdusk_bls12_381_sign.a
	# todo: also move bls12381svc binary, maybe embed?
	go build ./...

test: build
	go test -v ./...

bench: build
	go test -v -bench=. ./...

