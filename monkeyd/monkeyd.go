package monkeyd

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"log"
	"net"
)

// Monkeyd class
type Monkeyd struct {
	config *toml.TomlTree
}

func NewWithContent(configContent string) *Monkeyd {

	monkeyd := new(Monkeyd)

	config, err := toml.Load(configContent)
	if err != nil {
		fmt.Println("Error ", err.Error())
		panic(err)
	}

	monkeyd.config = config

	return monkeyd

}

// New Monkeyd
func New(configFile string) (*Monkeyd, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	return NewWithContent(string(content)), err
}

func (this *Monkeyd) Run(sectionStr string) {
	tree := this.config.Get(sectionStr).(*toml.TomlTree)
	fmt.Println(tree.ToString())
	typeStr := tree.Get("type")
	if typeStr == "server" {
		this.RunServer(tree.Get("forwardPort").(int64), tree.Get("servePort").(int64))
	} else {
		panic(fmt.Sprintf("Unknowd type str: %s", typeStr))
	}
}

func (this *Monkeyd) ServeConnHandler(conn net.Conn) {

	defer conn.Close()

	buf := make([]byte, 10240)
	for {
		// Recv
		_, err := conn.Read(buf)
		if err != nil {
			log.Printf("Read fail")
			return
		}

		// Send reply
		_, err = conn.Write(buf)
		if err != nil {
			log.Printf("Write fail")
			return
		}

	}
}

/*
运行服务器模式，传入服务端口，和转发到的端口
*/
func (this *Monkeyd) RunServer(forwardPort int64, servePort int64) {

	// serve port Conn Channel
	serveConnChan := make(chan net.Conn)
	go func() {
		for connChan := range serveConnChan {
            go this.ServeConnHandler(connChan)
		}
	}()

	// 同步开始监听端口
    forwardAddr := fmt.Sprintf("0.0.0.0:%d", forwardPort)
	serveAddr := fmt.Sprintf("0.0.0.0:%d", servePort)

    fmt.Printf("[RunServer]Server listen, on ServePort: %d, ForwardPort: %d \n", serveAddr, forwardPort)

    // serve listener
	serveListener, err := net.Listen("tcp", serveAddr)
	if err != nil {
		log.Printf("Error listen: %s", err)
		return
	}
	// 异步不停获取新连接 / serve
	go func() {

		defer serveListener.Close()

		for {
            fmt.Printf("!!!")
            serveConn, err := serveListener.Accept()
            fmt.Printf(".........")
            if err != nil {
                log.Printf("Error accept:%s", err.Error())
                return
            }
			serveConnChan <- serveConn
		}
	}()


    // forward listener
	forwardConnChan := make(chan net.Conn)
    go func() {
		for connChan := range forwardConnChan {
            go this.ForwardConnHandler(connChan)
		}
    }
    forwardListener, err := net.Listen("tcp", forwardAddr)

    if err != nil {
        log.Printf("Error accept forward: %s", err.Error())
        return
    }
    go func() {
       defer forwardListener.Close()

       for {
            forwardConn, err := forwardListener.Accept()
            if err != nil {
                 log.Printf("Error %s", err.Error())
                 return
            }
            forwardConnChan <- forwardConn
       }
    }()

}

func (this *Monkeyd) ForwardConnHandler(connChan chan net.Conn) {
}
