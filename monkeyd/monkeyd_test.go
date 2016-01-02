package monkeyd

import (
	"fmt"
	"net"
	"testing"
	"time"
	"bytes"
	"strings"
)

func TestStringEqual(t *testing.T) {
	if !strings.EqualFold("ABC", "ABC") {
		t.Error("Test string equal fail")
	}
}
func TestNewMonkeyd(t *testing.T) {

	m := NewWithContent(`[test_serve]
type = "server"
forwardPort = 3389
servePort = 33890
`)
	signal := make(chan bool)
	go func() {
		m.Run("test_serve")
		signal <- true
	}()

	fmt.Println("Wait 1s")
	time.Sleep(1 * time.Second)

	// client -> 用户机
	client, err := net.Dial("tcp", "127.0.0.1:33890")
	if err != nil {
		fmt.Println("error on dial test" + err.Error())
		return
	}
	defer client.Close()
	testStrArr := []byte("test string")
	testStr := string(testStrArr)
	//client.Write(testStrArr)

	// forward -> 转发机
	// little machine in intranet
    forward, err := net.Dial("tcp", "127.0.0.1:3389")
    if err != nil {
         fmt.Println(err.Error())
         return
    }
    defer forward.Close()
    // 转发机发送, client端收到
    forward.Write(testStrArr)

    // Client 接受到 foward的消息,经过 server的转发
	buf := make([]byte, 1024)
	_, err = client.Read(buf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

    readStr := string(bytes.TrimRight(buf[:], "\x00"))
	if len(readStr) != len(testStr) {
		t.Error(fmt.Sprintf("Test string length error, expect: %d, but: %d", len(readStr), len(testStr)))
	}
	if !strings.EqualFold(readStr, testStr) {
		t.Error(fmt.Sprintf("Test error expect: '%s', but: '%s' ", readStr, testStr))
	}




}
