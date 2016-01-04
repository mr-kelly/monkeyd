package monkeyd

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestStringEqual(t *testing.T) {
	if !strings.EqualFold("ABC", "ABC") {
		t.Error("Test string equal fail")
	}
}

// 创建一个InPort虚拟服务器,每隔一秒,不停发随机数据
func CreateTestInPortServer(inPort int64) {

	server, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", inPort))
	if err != nil {
		log.Errorf("Error: %s", err.Error())
		return
	}
	go func() {

		defer server.Close()
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Errorf(err.Error())
				return
			}
			log.Infof("[SimInPortServer]Accept! begin send data!")
			go func() {
				for {
					log.Infof("[SimInPortServer]Send data!")
					conn.Write([]byte("Hello!"))
					time.Sleep(3 * time.Second)
				}
			}()
		}
	}()

}
func TestNewMonkeyd(t *testing.T) {

	m := NewWithContent(`
[test_serve]
type = "server"
forwardPort = 33890
clientPort = 33891

[test_forward]
type = "forwarder"
inPort = 3389
serverAddress = "127.0.0.1:33891"
`)

	// 虚拟本地服务器
	CreateTestInPortServer(3389)

	signal := make(chan bool)
	go func() {
		m.Run("test_serve") // 坚挺client
		signal <- true
	}()

	go func() {
		m.Run("test_forward") // 坚挺forwarder
	}()

	fmt.Println("Wait 1s")
	time.Sleep(1 * time.Second)

	// client 模拟用户机
	client, err := net.Dial("tcp", "127.0.0.1:33891")
	if err != nil {
		fmt.Println("error on dial test" + err.Error())
		return
	}
	defer client.Close()
	testStrArr := []byte("test string")
	testStr := string(testStrArr)
	//client.Write(testStrArr)

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
