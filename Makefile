build :
	rm -rf bls12_381-sign || true
	git clone https://github.com/dusk-network/bls12_381-sign
	cd bls12_381-sign && git checkout microservice && cargo build --release
	protoc --proto_path=bls12_381-sign/proto bls12_381-sign/proto/bls12381sig.proto --go_opt=paths=source_relative --go_out=plugins=grpc:bls/proto
	 go build ./...

test: build
	go test -v ./...

bench: build
	go test -v -bench=. ./...

clean:
	rm /tmp/bls12381svc*
