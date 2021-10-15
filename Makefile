build:
	if [ -d bls/bls12_381-sign ]; \
	then \
		cd bls/bls12_381-sign \
		&& git pull; \
	else \
  		cd bls \
  		&& git clone https://github.com/dusk-network/bls12_381-sign; \
	fi;
	cd bls/bls12_381-sign \
	&& git checkout microservice \
	&& cd ../.. \
	&& protoc --proto_path=bls/bls12_381-sign/proto \
		bls/bls12_381-sign/proto/bls12381sig.proto \
		--go_opt=paths=source_relative --go_out=plugins=grpc:bls/proto  \
	&& cd bls/bls12_381-sign \
	&& cargo build --release \
	&& cd ../.. \
	&& cp bls/bls12_381-sign/target/release/bls12381svc bls12381svc_ubuntu-latest \
	&& cp bls/bls12_381-sign/target/release/libdusk_bls12_381_sign.a libdusk_bls12_381_sign_ubuntu-latest.a \
	&& go build ./...

test: build
	go test -v ./...

bench: build
	go test -v -bench=. ./...

clean:
	rm -fv /tmp/bls12381svc*
	rm -rfv bls/bls12_381-sign

installprotocubuntu: # like it says on the tin
	sudo apt install -y protobuf-compiler
	go install google.golang.org/grpc
	go install github.com/golang/protobuf/protoc-gen-go

memprofile:
	go test -run=. -bench=. -benchtime=5s -count 1 -benchmem -cpuprofile=cpu.out -memprofile=mem.out -trace=trace.out . | tee bench.txt
	go tool pprof -http :8081 mem.out