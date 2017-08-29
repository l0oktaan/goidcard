package MyCard

import (
    "fmt"
    "github.com/ebfe/scard"
	"bytes"    
    "io/ioutil"    
    "golang.org/x/text/encoding/traditionalchinese"
    "golang.org/x/text/transform"	
	
	//"github.com/varokas/tis620"
	"unicode/utf8"
	
)
const OFFSET = 0xd60
const WIDTH = 3

type Person struct{
	ID string
	Prefix string
	Fname string
	Lname string
	Addr string
	Birth string
	Sex string
}

func ToUTF8(tis620bytes []byte) []byte {
	l := findOutputLength(tis620bytes)
	output := make([]byte, l)

	index := 0
	buffer := make([]byte, WIDTH)
	for _, c := range tis620bytes {
		if !isThaiChar(c) {
			output[index] = c

			index++
		} else {
			utf8.EncodeRune(buffer, int32(c)+OFFSET)
			output[index] = buffer[0]
			output[index+1] = buffer[1]
			output[index+2] = buffer[2]

			index += 3
		}
	}
	return output
}

func findOutputLength(tis620bytes []byte) int {
	outputLen := 0
	for i, _ := range tis620bytes {
		if isThaiChar(tis620bytes[i]) {
			outputLen += WIDTH //always 3 bytes for thai char
		} else {
			outputLen += 1
		}
	}
	return outputLen
}

func isThaiChar(c byte) bool {
	return (c >= 0xA1 && c <= 0xDA) || (c >= 0xDF && c <= 0xFB)
}
func Decode(s []byte) ([]byte, error) {
    I := bytes.NewReader(s)
    O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
    d, e := ioutil.ReadAll(O)
    if e != nil {
        return nil, e
    }
    return d, nil
}
func CToGoString(c []byte) string {
    n := -1
    for i, b := range c {
        if b == 0 {
            break
        }
        n = i
    }
    return string(c[:n+1])
}

func ReadCard() string {
    // Establish a PC/SC context

    context, err := scard.EstablishContext()
    if err != nil {
        fmt.Println("Error EstablishContext:", err)
        return ""
    }

    // Release the PC/SC context (when needed)
    defer context.Release()

    // List available readers
    readers, err := context.ListReaders()
    if err != nil {
        fmt.Println("Error ListReaders:", err)
        return ""
    }

    // Use the first reader
    reader := readers[0]
    fmt.Println("Using reader:", reader)

    // Connect to the card
    card, err := context.Connect(reader, scard.ShareShared, scard.ProtocolAny)
    if err != nil {
        fmt.Println("Error Connect:", err)
        return ""
    }

    // Disconnect (when needed)
    defer card.Disconnect(scard.LeaveCard)

    // Send select APDU
    var cmd_select = []byte{0x00, 0xA4, 0x04, 0x00, 0x08, 0xA0, 0x00, 0x00, 0x00, 0x54, 0x48, 0x00, 0x01}
    rsp, err := card.Transmit(cmd_select)
    if err != nil {
        fmt.Println("Error Transmit:", err)
        return ""
    }
    fmt.Println(rsp)

    // Send command APDU
    var cmd_command1 = []byte{0x80, 0xb0, 0x00, 0x11, 0x02, 0x00, 0xd1}
	var cmd_command2 = []byte{0x00, 0xc0, 0x00, 0x00, 0xd1}
	
    rsp, err = card.Transmit(cmd_command1)
    if err != nil {
        fmt.Println("Error Transmit:", err)
        return ""
    }
    //fmt.Println(rsp)
	rsp, err = card.Transmit(cmd_command2)
    if err != nil {
        fmt.Println("Error Transmit:", err)
        return ""
    }

    //fmt.Println(rsp)
    //tt := []byte(rsp)
    //fmt.Println(tt)
	//t,err := Decode(rsp)
	//fmt.Println(string(b))
	//fmt.Fprintf(w, "<h1>%s</h1>",tt)
    //fmt.Println(string(b))
    //fmt.Println(err)
	//str := CToGoString(rsp[:])
	
	
	x := ToUTF8(rsp)
	return string(x)
    //fmt.Println(string(x))
    //fmt.Fprintf(w, "<h1>%s</h1>",x)    
    
}
