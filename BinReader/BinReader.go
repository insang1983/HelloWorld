package main

import (
	"fmt"
	"os"
	//"io"
	"io"
)

/*const (

)*/

type wtmp struct {
	uttype int16
	pid int32
	line [32]byte
	id [4]byte
	user [32]byte
	host [256]byte
	etermination int16
	eexit int16
	session int32
	tvsec int32
	tvusec int32
	addrv6 [4]int32
	unused [20]byte
}

/* A Simple function to verify error */
func CheckError(err error) {
	if err  != nil {
		fmt.Println("Error: " , err)
		os.Exit(0)
	}
}

func main() {
	f,err := os.Open("D:/Project/BinReader/wtmp")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	var wt wtmp

	for {
		n, err := f.Read(wt)
		n, err := != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return
		}
	}
}