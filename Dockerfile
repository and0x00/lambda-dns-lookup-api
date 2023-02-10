FROM golang:latest
WORKDIR /build
ENTRYPOINT ["env","GOOS=darwin","GOARCH=arm64", "go", "build","-o","./bin"]