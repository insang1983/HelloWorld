package main

import (
	"fmt"
	"os"
	"io"
	"encoding/binary"
	"math/big"
	"net"
	//"flag"
	//"time"
	"time"
	"unsafe"
)

type wtmp struct {
	Uttype int32
	Pid int32
	Line [32]byte
	Id [4]byte
	User [32]byte
	Host [256]byte
	Etermination int16
	Eexit int16
	Session int32
	Tvsec int32
	Tvusec int32
	Addrv6 [4]int32
	Unused [20]byte
}

/* A Simple function to verify error
func CheckError(err error) {
	if err  != nil {
		fmt.Println("Error: " , err)
		os.Exit(0)
	}
}
*/

func ConvByteToString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

func IP6toInt(IPv6Address net.IP) *big.Int {
	IPv6Int := big.NewInt(0)

	// from http://golang.org/pkg/net/#pkg-constants
	// IPv6len = 16
	IPv6Int.SetBytes(IPv6Address.To16())
	return IPv6Int
}

func main() {
	Argc := len(os.Args)
	var filename string
	var tailflag int = 0

	i := 0
	for i < len(os.Args) {
		fmt.Printf("[%d]%s\n", i, os.Args[i])
		i++
	}

	if Argc == 2 {
		filename = os.Args[1]
	} else if Argc == 3 {
		//tailflag = flag.Int("t", 0, "tail option")
		if os.Args[1] == "-t" {
			tailflag = 1
		}
		filename = os.Args[2]
		//flag.Parse()
	} else {
		return
	}

	f,err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	wt := new(wtmp)

	wtsize := unsafe.Sizeof(wt)
	fmt.Printf("wsize : %d\n", wtsize)

	for {

		if tailflag == 1 {
			f.Seek(-384*5, os.SEEK_END)
			tailflag = 2
		}

		err = binary.Read(f, binary.LittleEndian, wt)
		if err == io.EOF {
			if tailflag == 2 {
				continue
			} else {
				break
			}
			//return
		}

		i++

		/*
		Mwtmp := map[string]interface{}{}
		Mwtmp,_ := json.Marshal(wt)

		Mwtmp["type"] = wt.Uttype
		Mwtmp["pid"] = wt.Pid
		Mwtmp["line"] = fmt.Sprintf("<%s>", bytes.Trim(wt.Line, "0x00"))
		Mwtmp["id"] = wt.Id
		Mwtmp["user"] = wt.User
		Mwtmp["host"] = wt.Host
		Mwtmp["etermination"] = wt.Etermination
		Mwtmp["eexit"] = wt.Eexit
		Mwtmp["session"] = wt.Session
		Mwtmp["tvsec"] = wt.Tvsec
		Mwtmp["tvusec"] = wt.Tvusec
		Mwtmp["addrv6"] = fmt.Sprintf("%x %x %x %x, ", wt.Addrv6[0], wt.Addrv6[1], wt.Addrv6[2], wt.Addrv6[3])
		Mwtmp["unused"] = wt.Unused
		fmt.Println(Mwtmp)
		*/

		fmt.Printf("type : [%x], ", wt.Uttype)
		fmt.Printf("pid : [%d], ", wt.Pid)
		fmt.Printf("line : [%s], ", ConvByteToString(wt.Line[:]) )
		fmt.Printf("id : [%s], ", ConvByteToString(wt.Id[:]))
		fmt.Printf("user : [%s], ", ConvByteToString(wt.User[:]))
		fmt.Printf("host : [%s], ", ConvByteToString(wt.Host[:]))
		fmt.Printf("etermination : [%d], ", wt.Etermination)
		fmt.Printf("eexit : [%x], ", wt.Eexit)
		fmt.Printf("session : [%x], ", wt.Session)
		fmt.Printf("time : %s, ", time.Unix(int64(wt.Tvsec), int64(wt.Tvusec)))
		//fmt.Printf("tvsec : [%d], ", wt.Tvsec)
		//fmt.Printf("tvusec : [%x], ", wt.Tvusec)

		fmt.Printf("addrv6 : [%x %x %x %x], ", wt.Addrv6[0], wt.Addrv6[1], wt.Addrv6[2], wt.Addrv6[3])

		fmt.Printf("unused : [%s]\r\n", ConvByteToString(wt.Unused[:]))
	}

	fmt.Printf("Read count : %d ", i)
}