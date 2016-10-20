package main

import (
	"fmt"
	"os"
	"net"
	"time"
	"bufio"
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

	yyyymmddhh24mi := "200601021504"
	yyyymmddhh24miss := "20060102150405"
	yyyymmdd_hh24miss := "20060102 150405"

	buf := make([]byte, 4096)
	var recv_msg string
	var iCount int = 0

	now := time.Now().Local()
	before := now
	prevTime := now

	filename := now.Format(yyyymmddhh24mi) + ".txt"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0666))
	if err != nil {
		fmt.Println(now.Format(yyyymmdd_hh24miss),err)
		os.Exit(2)
	}
	fmt.Println(now.Format(yyyymmdd_hh24miss),filename, "is opend..")
	file_obj := bufio.NewWriter(file)
	file.Seek(0, os.SEEK_END)

	for {
		now = time.Now()

		if now.Local().Format(yyyymmddhh24miss) != prevTime.Local().Format(yyyymmddhh24miss) {
			fmt.Println(now.Format(yyyymmdd_hh24miss), "eps : ", iCount )
			iCount = 0
			prevTime = now
		}

		if now.Local().Format(yyyymmddhh24mi) != before.Local().Format(yyyymmddhh24mi) {
			before = now
			if file != nil {
				file.Close()
				file = nil
			}
			filename := now.Format(yyyymmddhh24mi) + ".txt"
			file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0666))
			if err != nil {
				fmt.Println(now.Format(yyyymmdd_hh24miss), err)
				continue 
			}
			file.Seek(0, os.SEEK_END)
			file_obj = bufio.NewWriter(file)
			fmt.Println(now.Format(yyyymmdd_hh24miss), filename, "is opend..")
		}

		ServerConn.SetReadDeadline(time.Now().Add(1*time.Second))
		n, addr, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue; //timeout
			} else {
				fmt.Println(now.Format(yyyymmdd_hh24miss), "Error: ", err)
			}
		}

		iCount++

		if n > 0 {
			recv_msg = "["+ addr.IP.String()+ "]["+ now.Format(yyyymmdd_hh24miss) + "]["+ string(buf[0:n])+ "]"
			file_obj.WriteString(recv_msg + "\r\n")
			//file_obj.Flush()
		}
	}
}