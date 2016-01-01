package monkeyd

import (
    "fmt"
    "github.com/pelletier/go-toml"
)
// Monkeyd class
type Monkeyd struct {
	configFile string // toml config file path
    config *toml.TomlTree
}

// New Monkeyd
func New(configFile string) *Monkeyd {

    monkeyd := new(Monkeyd)

    monkeyd.configFile = configFile

    config, err := toml.LoadFile(configFile)
    if (err != nil) {
        fmt.Println("Error ", err.Error())
        panic(err)
    }

    monkeyd.config = config

    return monkeyd

}

func (this *Monkeyd) Run(sectionStr string) {
    fmt.Println("section string: %s", sectionStr)
    tree := this.config.Get(sectionStr).(*toml.TomlTree)
    fmt.Println(tree.ToString())
    for _, key := range tree.Keys() {
        value := tree.Get(key)
        fmt.Printf("Key: %s, Value: %s\n", key, value)
    }
}
