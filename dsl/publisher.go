package dsl

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"strings"
	"time"

	"github.com/masgari/homie-go/homie"
)

type periodicPublisherConfig struct {
	Every string // how frequently publish, in time.Duration format: 300ms, 2s, 1m, 5m, 1h
	Nodes []map[string]interface{}
}

type propertyPublisherFunc func(p homie.Property)

type dslNodePublisher struct {
	propertyFuncs map[string]propertyPublisherFunc
}

func (d *dslNodePublisher) publisherFunc() homie.NodePublisher {
	return func(n homie.Node) {
		for pName, f := range d.propertyFuncs {
			p := n.GetProperty(pName)
			f(p)
		}
	}
}

func configurePublishers(device homie.Device, cfg *dslConfig) []interface{} {
	publishers := make([]interface{}, 0, len(cfg.Publishers))

	for name, c := range cfg.Publishers {
		switch strings.ToLower(name) {
		case "periodicpublisher":
			periodicPublisher := configurePeriodicPublisher(device, c)
			publishers = append(publishers, periodicPublisher)
		default:
			panic(fmt.Errorf("unsupported publisher type: %s", name))
		}
	}
	return publishers
}

func configurePeriodicPublisher(device homie.Device, c interface{}) homie.PeriodicPublisher {
	var pCfg []periodicPublisherConfig
	if err := mapstructure.Decode(c, &pCfg); err != nil {
		panic(err)
	}
	cfg := pCfg[0]

	duration, err := time.ParseDuration(cfg.Every)
	if err != nil {
		panic(err)
	}

	if len(cfg.Nodes) < 1 {
		panic(fmt.Errorf("no node configured for periodic publisher: %v", cfg))
	}

	publisher := homie.NewPeriodicPublisher(duration)

	for nodeName, propsCfg := range cfg.Nodes[0] {
		node := device.GetNode(nodeName)
		if node == nil {
			panic(fmt.Errorf("unknown node: %s in periodic publisher", nodeName))
		}
		publisher.AddNodePublisher(node, makeNodePublisher(node, propsCfg))
	}
	return publisher
}

func decodeGenericMap(c interface{}) []map[string]interface{} {
	var cfg []map[string]interface{}
	if err := mapstructure.Decode(c, &cfg); err != nil {
		panic(err)
	}
	return cfg
}

func makeNodePublisher(node homie.Node, c interface{}) homie.NodePublisher {
	cfg := decodeGenericMap(c)
	if len(cfg) < 1 {
		panic(fmt.Errorf("no function configured for node: %s in periodic publisher: %v", node.Name(), cfg))
	}

	holder := dslNodePublisher{}
	holder.propertyFuncs = make(map[string]propertyPublisherFunc, len(cfg[0]))
	for propName, funcCfg := range cfg[0] {
		prop := node.GetProperty(propName)
		if prop == nil {
			panic(fmt.Errorf("unknown property: %s for node: %s in periodic publisher: %v", propName, node.Name(), funcCfg))
		}
		funcConfigMap := decodeGenericMap(funcCfg)[0]
		if len(funcConfigMap) != 1 {
			panic(fmt.Errorf("only one function configuration is allowed, property: %s of node: %s in periodic publisher: %v", propName, node.Name(), funcConfigMap))
		}
		holder.propertyFuncs[propName] = makePropertyFunc(prop, funcConfigMap)
	}

	return holder.publisherFunc()
}

func makePropertyFunc(property homie.Property, cfg map[string]interface{}) propertyPublisherFunc {
	for funcName, funcCfg := range cfg {
		switch strings.ToLower(funcName) {
		case "increment", "inc":
			return makeNumericPropertyModifier("increment", funcCfg).propFunc()
		case "decrement", "dec":
			return makeNumericPropertyModifier("decrement", funcCfg).propFunc()
		case "current-value", "current":
			return makeNumericPropertyModifier("current-value", funcCfg).propFunc()
		default:
			panic(fmt.Errorf("unknown function, property: %s of node: %s in periodic publisher: %s", property.Name(), property.Node().Name(), funcName))
		}
	}
	panic(fmt.Errorf("unhandled function, property: %s of node: %s in periodic publisher: %v", property.Name(), property.Node().Name(), cfg))
}
