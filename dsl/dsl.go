package dsl

import (
	"io/ioutil"
	"log"

	"github.com/hashicorp/hcl"

	homie "github.com/masgari/homie-go/homie"
)

type propertyInitConfig struct {
	Value string
}

type propertyConfig struct {
	Type string
	Init propertyInitConfig
}

type nodeConfig struct {
	Type       string
	Properties map[string]propertyConfig
}

type dslConfig struct {
	Name       string
	Config     homie.Config
	Nodes      map[string]nodeConfig
	Publishers map[string]interface{}
}

// DeviceWrapper wrapper around homie.Device
type DeviceWrapper struct {
	Device          homie.Device
	DevicePublisher homie.PeriodicPublisher
	NodePublishers  []interface{}
}

// Run create homie device from file and run it with default publisher
func Run(file string) {
	dsl := LoadFile(file)[0]
	dsl.Device.Run(true)
}

// LoadFile create homie devices(s) from file
func LoadFile(file string) []DeviceWrapper {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Panic(err)
	}
	return Load(data)
}

// Load create homie instance(s) from DSL definition
func Load(bs []byte) []DeviceWrapper {
	var cfgMap map[string]*dslConfig
	if err := hcl.Unmarshal(bs, &cfgMap); err != nil {
		log.Panic(err)
	}
	dsls := make([]DeviceWrapper, 0, len(cfgMap))
	for name, cfg := range cfgMap {
		device := homie.NewDevice(name, &cfg.Config)
		configureDevice(device, cfg)
		dsl := DeviceWrapper{
			Device:          device,
			DevicePublisher: homie.NewDevicePublisher(device),
			NodePublishers:  configurePublishers(device, cfg),
		}

		dsls = append(dsls, dsl)
	}
	return dsls
}

func configureDevice(device homie.Device, cfg *dslConfig) {
	for name, c := range cfg.Nodes {
		node := device.NewNode(name, c.Type)
		configureNode(node, &c)
	}
}

func configureNode(node homie.Node, cfg *nodeConfig) {
	for name, c := range cfg.Properties {
		p := node.NewProperty(name, c.Type)
		p.SetValue(c.Init.Value)
	}
}
