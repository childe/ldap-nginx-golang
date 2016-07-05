hash:=$(shell git rev-parse --short HEAD)

default: nginx-ldap-auth-daemon

all:
	echo $(hash)
	mkdir -p build/$(hash)

	GOOS=windows GOARCH=amd64 go build -o build/$(hash)/nginx-ldap-auth-daemon-windows-x64-$(hash).exe
	GOOS=windows GOARCH=386 go build -o build/$(hash)/nginx-ldap-auth-daemon-windows-386-$(hash).exe
	GOOS=linux GOARCH=amd64 go build -o build/$(hash)/nginx-ldap-auth-daemon-linux-x64-$(hash)
	GOOS=linux GOARCH=386 go build -o build/$(hash)/nginx-ldap-auth-daemon-linux-386-$(hash)
	GOOS=darwin GOARCH=amd64 go build -o build/$(hash)/nginx-ldap-auth-daemon-darwin-x64-$(hash)

nginx-ldap-auth-daemon:
	mkdir -p build/$(hash)
	go build -o build/nginx-ldap-auth-daemon
