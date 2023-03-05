package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host\n", os.Args[0])
		os.Exit(1)
	}

	host := os.Args[1]
	addr, err := net.ResolveIPAddr("ip4:icmp", host)
	if err != nil {
		fmt.Println("Error resolving IP address:", err)
		os.Exit(1)
	}

	conn, err := net.DialIP("ip4:icmp", nil, addr)
	if err != nil {
		fmt.Println("Error connecting to host:", err)
		os.Exit(1)
	}
	defer conn.Close()

	var seq int16 = 1
	var msg [512]byte
	msg[0] = 8 // ICMP echo request
	msg[1] = 0 // ICMP echo request
	msg[2] = 0 // checksum (2 bytes)
	msg[3] = 0 // checksum (2 bytes)
	msg[4] = byte(seq >> 8)
	msg[5] = byte(seq & 0xff)
	len := 8

	check := checkSum(msg[0:len])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 0xff)

	start := time.Now()
	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	_, err = conn.Write(msg[0:len])
	if err != nil {
		fmt.Println("Error sending ICMP request:", err)
		os.Exit(1)
	}

	recv := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conn.Read(recv)
	if err != nil {
		fmt.Println("Error receiving ICMP response:", err)
		os.Exit(1)
	}
	elapsed := time.Since(start)
	fmt.Printf("Received %d bytes from %s: icmp_seq=%d time=%v\n", n, addr.String(), seq, elapsed)
}

func checkSum(msg []byte) uint16 {
	sum := 0
	for n := 1; n < len(msg)-1; n += 2 {
		sum += int(msg[n])*256 + int(msg[n+1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)
	return uint16(^sum)
}
