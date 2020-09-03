module github.com/masgari/homie-dsl

go 1.12

require (
	github.com/hashicorp/hcl v1.0.0
	github.com/masgari/homie-go v0.0.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/stretchr/testify v1.3.0
)

replace github.com/masgari/homie-go => ../homie-go
