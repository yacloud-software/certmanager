.PHONY:	client server 
all:	client server
client:
	cd client && go install certmanager-client.go
cgi:
	cd cgi && go install certmanager-cgi.go `ls -1 *.go | grep -v certmanager-cgi.go`
server:
	cd server && go install certmanager-server.go `ls -1 *.go|grep -v certmanager-server.go`
