package main

import (
	"fmt"
	"log"
	"net"
	"encoding/binary"
	"golang.org/x/crypto/ssh"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:22")

	if err != nil {
		log.Fatalf("Failed to listen on port 22 (%s)", err)
	}

	log.Printf("Listening on port 0.0.0.0:22")

	//configs
	serverVersion := "SSH-2.0-OpenSSH_7.6p1 Ubuntu-4ubuntu0.3"
	passwordCallback := func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
		//fmt.Printf("%+v\n", conn)
		//log.Printf("user: %s, password: %s", conn.User(), string(password))
		//log.Printf("SessionID: %d", binary.BigEndian.Uint64(conn.SessionID()))
		//log.Printf("ClientVersion: %s", conn.ClientVersion())
		//log.Printf("ServerVersion: %s", conn.ServerVersion())
		//log.Printf("RemoteAddr: %s", conn.RemoteAddr())

		log.Printf("%s;%s;%d;%s;%s", conn.User(), string(password), binary.BigEndian.Uint64(conn.SessionID()), conn.ClientVersion(), conn.RemoteAddr())
		
		//return &ssh.Permissions{}, nil
		return nil, fmt.Errorf("authentication failed")
	}

	//authentication
	sshConfig := ssh.ServerConfig{
		PasswordCallback: passwordCallback,
		ServerVersion: serverVersion,
	}

	hostKeyData := []byte(getHostKey())
	signer, err := ssh.ParsePrivateKey(hostKeyData)
	if err != nil {
		log.Fatalf("Failed to parse host key (%s)", err)
	}
	sshConfig.AddHostKey(signer)

	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection (%s)", err)
			continue
		}

		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, &sshConfig)
		if err != nil {
			//log.Printf("Failed handshake (%s)", err)
			continue
		}

		log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())

		go ssh.DiscardRequests(reqs)
		go handleChannels(sshConn, chans)
	}
}

func handleChannels(conn *ssh.ServerConn, chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		go handleChannel(conn, newChannel)
	}
}

func handleChannel(conn *ssh.ServerConn, newChannel ssh.NewChannel) {
	log.Printf("conn: %d", conn)
	log.Printf("channel: %d", newChannel)
}

func getHostKey() (string) {
	return "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAu3aX0lPqijyTj54AxZlG+y5AyigssfAKP+4DrefvD6NQY37H\nGZb0FNNU9nxHe1pFBQqd6+J1otTrwjnbreU3tLVr4MmFAzbO/qfKjsSFID9tfv0w\nCOp+b8PQwSEN1pJh2maC/TdqrEZvtaUpouNCAKILWlMuzCn6lKFmyDzGbLuBsldb\n4Qh6pTTB+Qn7bHh5fKscLBCqfnn//x0jDETVsboxu6kxMh7yD+BRC2dYt5vU1qRD\ngpJNvG65PiNmCbG6bg6GQWnAQvot71UJWc0V3y/7ueyQmyZi//BH3KCBknHbw9AJ\n2+NkpumiZLYCPukeM5tWamfuZn/lsvIc/OLmMwIDAQABAoIBABFUExCkJSgGFXXP\nGy8ozgDl86M8N3VzRN9H9xsaN2TwqbuoumrJI2LRbicisdDbNUoUAykM0+brW6em\nhYH7kDyqEIDE4AC+DkBH7ldoHw4uDscQTGJmmq1mImPX1FmjSlxP4YfamYe3MxhN\nXh3qd+1rDTWaPtcsgjc+/CtCQu1F91SDTk37zx1I+7yjRKo7Jk4c4n2xsOGwGfeo\nl3vFXqOBP28UUfoWMAzAAHk8QLPPgnddLFjw7wIJEQPvPdzEanW62KrZlbJq1w/s\nozHBLTAr5uMJeZvl9/th/PfG+Xi4c7S3hkRWG4h9J0H6YVZA8Tu+/8YLIVM9VJhC\nKSkOoaECgYEA59NChqTmqieHLzbe1U2v5c/sBAZJ7e/+MyG7ANMPtHHECKJZ0G8f\n6LFxD4gAkDfq4Yc41/m7WRC4xlEbi4pU9L7gVNLQNfNfiL0ysmZIIChtY7IGAGni\nRmiEq9/k+ztnaSOs7CE6FyFrk0USesySqZ5FWNbt3nhUzO01Q7cW32kCgYEAzwMO\n9354/Sdh+PRpzhNMD8akh9qMkKJjThm0UeAXXOGlaGNFss9UI4eB4d9+XarevyBx\nmhkjsuBG2Zz7m3KsVbGpJcgXh/HErAJlJRHt+UqIfIX71YW3Ol5fgjg8BXcpBk+K\nhcZPC30uSU66zM7Ja1eLbNcNO1sqrmI64xwOATsCgYEAh/Ty6AofqRzDgGIar1f/\nV7ToArg5dUyxdQVMKcCeTkIGKNYl/EKfoRUnbGdjhTD2FEv8f1VblXFkHBKHJ//5\nsQucfsKgD3PqzEPBTrUDibCL7tMCCA4RAR/c5vvIy7pb/GJK0LTv347fCyCQJOqC\n/OzwWJi8KiPB/+kBuvPOezkCgYAcrLnIApbTyj7B82ksiHPCw6tKvjU2W6gRy3G0\n3auezArTeNzQtfNbsIuHNCQW6XJNWzshM1ZEkth9kEcx8yJ4BFH/z8WiqRSrFvHX\nvrIOFArv5MdLfmgxB52HNi7qOuVN4Hq5qQyN9NsSgHtTn1k7Kzc+7lMA49H3sdei\nWeJ+vQKBgQCyb5QQDHfmXSpDFmr+jsBdZ+Z0FeW0yyp/FlWjLQwor1MdGoJG+bJM\nehsUbLdqXPZPlO3eC7aM1CREJ/0CiQ68DZ9PfkXcdEEAOhM5lpj9PXSrhRbGzkfn\n71T8RrKdeFt14JbD7BMSxXGilOIqw2qIchSgPIpItvM5w2ynxwdZdA==\n-----END RSA PRIVATE KEY-----"
}