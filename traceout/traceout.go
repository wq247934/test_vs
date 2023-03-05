package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const maxHops = 64
const timeout = 5 * time.Second

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host\n", os.Args[0])
		os.Exit(1)
	}

	host := os.Args[1]
	addr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		fmt.Println("Error resolving IP address:", err)
		os.Exit(1)
	}

	fmt.Printf("Tracing route to %s [%s]\n", host, addr.String())

	for i := 1; i <= maxHops; i++ {
		fmt.Printf("%2d ", i)

		conn, err := net.DialTimeout("ip4:icmp", addr.String(), timeout)
		if err != nil {
			fmt.Println("Error connecting to host:", err)
			continue
		}
		defer conn.Close()

		start := time.Now()
		msg := make([]byte, 64)
		msg[0] = 8       // ICMP echo request
		msg[1] = 0       // ICMP echo request
		msg[2] = 0       // checksum (2 bytes)
		msg[3] = 0       // checksum (2 bytes)
		msg[4] = 0       // identifier (2 bytes)
		msg[5] = 0       // identifier (2 bytes)
		msg[6] = byte(i) // sequence number (1 byte)
		msg[7] = 0       // sequence number (1 byte)

		check := checkSum(msg[0:8])
		msg[2] = byte(check >> 8)
		msg[3] = byte(check & 0xff)

		conn.SetDeadline(time.Now().Add(timeout))
		_, err = conn.Write(msg[0:8])
		if err != nil {
			fmt.Println("Error sending ICMP request:", err)
			continue
		}

		recv := make([]byte, 1024)
		n, err := conn.Read(recv)
		if err != nil {
			fmt.Println("Error receiving ICMP response:", err)
			continue
		}
		elapsed := time.Since(start)
		if n >= 20 {
			// ttl := int(recv[8])
			code := int(recv[20])
			fmt.Printf("%15s %4dms",
				conn.RemoteAddr().String(), elapsed/time.Millisecond)
			switch code {
			case 0:
				fmt.Printf("  Reached destination\n")
				return
			case 1:
				fmt.Printf("  Host unreachable\n")
			case 2:
				fmt.Printf("  Protocol unreachable\n")
			case 3:
				fmt.Printf("  Port unreachable\n")
			case 4:
				fmt.Printf("  Fragmentation needed but DF flag set\n")
			case 5:
				fmt.Printf("  Source route failed\n")
			case 6:
				fmt.Printf("  Destination network unknown\n")
			case 7:
				fmt.Printf("  Destination host unknown\n")
			case 8:
				fmt.Printf("  Source host isolated\n")
			case 9:
				fmt.Printf("  Communication with destination network administratively prohibited\n")
			case 10:
				fmt.Printf("  Communication with destination host administratively prohibited\n")
			case 11:
				fmt.Printf("  Network unreachable for Type Of Service\n")
			case 12:
				fmt.Printf("  Host unreachable for Type Of Service\n")
			case 13:
				fmt.Printf("  Communication administratively prohibited by filtering\n")
			case 14:
				fmt.Printf("  Host precedence violation\n")
			case 15:
				fmt.Printf("  Precedence cutoff in effect\n")
			default:
				fmt.Printf("  Unknown error code %d\n", code)
			}
		}
	}
}

func checkSum(msg []byte) uint16 {
	sum := uint32(0)
	for n := 1; n < len(msg)-1; n += 2 {
		sum += uint32(msg[n-1])<<8 + uint32(msg[n])
	}
	if len(msg)%2 == 1 {
		sum += uint32(msg[len(msg)-1]) << 8
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	return ^uint16(sum)
}
