package main

import (
	"fmt"
	"os"
	"net"
	"time"
	//"bufio"
)

/* A Simple function to verify error */
func CheckError(err error) {
	if err  != nil {
		fmt.Println("Error: " , err)
		os.Exit(0)
	}
}

func main() {
	/* Lets prepare a address at any address at port 10001*/
	ServerAddr,err := net.ResolveUDPAddr("udp",":10001")
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)
	var file *os.File = nil
	var iCount int = 0
	prevTime := time.Now().Local()

	for {
		now := time.Now().Local()

		if now.Unix() - prevTime.Unix() >= 1 {
			fmt.Println(now.Format("20060102 150405"), "eps : ", iCount )
			iCount = 0
			prevTime = now
		}

		filename := now.Format("200601021504") + ".txt"
		if _, err := os.Stat(filename); os.IsNotExist(err){
			if file != nil {
				file.Close()
				file = nil
			}
			file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
			if err != nil {
				fmt.Println(err)
				continue 
			}
		} else {
			//있으면 열렸는지?
			//열렸으면 close 하고 open
			if file == nil {
				file, err = os.OpenFile(filename, os.O_WRONLY, os.FileMode(0644))
				if err != nil {
					fmt.Println(err)
					//return
					continue
				}
			} 
			file.Seek(0, os.SEEK_END)
		}

		ServerConn.SetReadDeadline(time.Now().Add( 1*time.Nanosecond))
		n, addr, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				//timeout
				continue;
			} else {
				fmt.Println("Error: ", err)
			}
		}

		iCount++ // 0인 경우를 위해 먼저 건수 증가

		if n > 0 {
			recv_msg := "["+ addr.IP.String()+ "]["+ now.Format("20060102 150405") + "]["+ string(buf[0:n])+ "]"
			file.WriteString(recv_msg + "\r\n")
		}

	}
}