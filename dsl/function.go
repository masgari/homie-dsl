package dsl

import (
	"fmt"
	homie "github.com/masgari/homie-go/homie"
	"github.com/mitchellh/mapstructure"
	"strconv"
)

type numericPropertyModifier struct {
	step     float64
	operator string
}

func makeNumericPropertyModifier(operator string, funcCfg interface{}) *numericPropertyModifier {
	var propModifier []numericPropertyModifier
	if err := mapstructure.Decode(funcCfg, &propModifier); err != nil {
		panic(err)
	}
	propModifier[0].operator = operator
	fmt.Printf(">>>propModifier: %v\nfuncCfg:\n%v", propModifier[0], funcCfg)
	return &propModifier[0]
}

func (m *numericPropertyModifier) propFunc() propertyPublisherFunc {
	return func(p homie.Property) {
		nValue, err := strconv.ParseFloat(p.Value(), 64)
		if err != nil {
			panic(err)
		}
		fmt.Printf(">>>p:%s, value: %f, step: %f\n", p.Name(), nValue, m.step)
		switch m.operator {
		case "increment":
			nValue += m.step
		case "decrement":
			nValue += m.step
		case "current-value":
		}

		switch p.Type() {
		case "int", "integer", "number":
			p.SetValue(strconv.FormatInt(int64(nValue), 10))
		case "float", "double":
			p.SetValue(strconv.FormatFloat(nValue, 'f', -1, 64))
		}
	}
}
