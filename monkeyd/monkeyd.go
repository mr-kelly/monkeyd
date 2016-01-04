package monkeyd

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

var log = logging.MustGetLogger("monkeyd")

// Monkeyd class
type Monkeyd struct {
	config *toml.TomlTree
}

// New a struct with config file content
func NewWithContent(configContent string) *Monkeyd {

	monkeyd := new(Monkeyd)

	config, err := toml.Load(configContent)
	if err != nil {
		log.Errorf("Error ", err.Error())
		panic(err)
	}

	monkeyd.config = config

	return monkeyd

}

// New Monkeyd, with file path
func New(configFilePath string) (*Monkeyd, error) {
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	return NewWithContent(string(content)), err
}

// Run a section, base on the `type` key type
func (this *Monkeyd) Run(sectionStr string) {
	tree := this.config.Get(sectionStr).(*toml.TomlTree)
	fmt.Println(tree.ToString())
	typeStr := tree.Get("type")
	if typeStr == "server" {
		forwardPort := tree.Get("forwardPort").(int64)
		clientPort := tree.Get("clientPort").(int64)
		this.RunServer(forwardPort, clientPort)
	} else if typeStr == "forwarder" {
		inPort := tree.Get("inPort").(int64)
		serverAddress := tree.Get("serverAddress").(string)
		this.RunForwarder(inPort, serverAddress)
	} else {
		panic(fmt.Sprintf("Unknowd type str: %s", typeStr))
	}
}
