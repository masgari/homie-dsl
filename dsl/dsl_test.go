package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"
	homie "github.com/masgari/homie-go/homie"
)

func TestCreateDeviceFromConfig(t *testing.T) {
	//d := homie.NewDevice()
	definition := `
	test {
		config {
			mqtt {
				host = "localhost"
				port = 1883
				username = "test-user"
				password = "random-password"
			}
			baseTopic = "test_devices/"
			statsReportInterval = 37
		}

		nodes {
			n1 {
				type = "Test Node"

				properties {
					p1 {
						type = "integer"
						init {
							value = "7"
						}
					}
					p2 {
						type = "float"
						init {
							value = "2.3"
						}
					}
				}
			}
			n2 {
				type = "TestNode2"
				properties {
					p3 {
						type = "string"
						init {
							value = "json"
						}
					}
				}
			}
		}
		publishers {
			periodicPublisher {
				every = "2s"
				nodes {
					n1 {
						p1 {
							"increment" {
								step = 2
							}
						}
						p2 {
							dec {
								step = -0.1
							} 
						}
					}
					n2 {
						p3 {
							current-value {
								step = 2
							}
						}
					}
				}
			}
		}
	}
	`

	dsls := Load([]byte(definition))
	assert.NotEmpty(t, dsls)
	assert.Equal(t, 1, len(dsls))
	device := dsls[0].Device
	assert.NotEmpty(t, device)
	assert.Equal(t, "test", device.Name())
	assert.Equal(t, "localhost", device.Config().Mqtt.Host)
	assert.Equal(t, 1883, device.Config().Mqtt.Port)
	assert.Equal(t, "test-user", device.Config().Mqtt.Username)
	assert.Equal(t, "random-password", device.Config().Mqtt.Password)
	assert.Equal(t, "test_devices/", device.Config().BaseTopic)
	assert.Equal(t, 37, device.Config().StatsReportInterval)

	n1 := device.GetNode("n1")
	assert.NotEmpty(t, n1)
	assert.Equal(t, "Test Node", n1.Type())

	assert.NotEmpty(t, n1.NodePublisher())

	assert.NotEmpty(t, n1.GetProperty("p1"))
	assert.NotEmpty(t, n1.GetProperty("p2"))
	assert.Equal(t, "float", n1.GetProperty("p2").Type())

	assert.Equal(t, "7", n1.GetProperty("p1").Value())
	assert.Equal(t, "2.3", n1.GetProperty("p2").Value())

	assert.NotEmpty(t, dsls[0].NodePublishers)
	assert.Implements(t, (*homie.PeriodicPublisher)(nil), dsls[0].NodePublishers[0])

	pp := dsls[0].NodePublishers[0].(homie.PeriodicPublisher)
	np := pp.GetNodePublisher(n1)
	assert.NotEmpty(t, np)
}
