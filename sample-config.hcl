/*
Sample configuration, covers the idea, nothing is implemted yet
*/
sys-monitor {
    config {
        mqtt {
            host = "localhost"
            port = 1883
            username = "test-user"
            password = "random-password"
        }
    }
    
    nodes {
        cpu { // id of the node
            type = "Monitoring" 
            properties { // node properties
                usage { // name of property
                    type = "float"
                    init {
                        value = "0.0"
                        shell {
                            
                        }
                    }
                }

                avg {

                }

                enabled { // settable property
                    initial = true
                    type = bool
                    publish {
                        onconnect {

                        }
                    }
                    // handler to react to incoming message
                    handler {
                        javascript {
                            // $node and $mqtt are global variables in js vm
                            $node.enabled = $mqtt.payload
                        }

                        // specific handler to react to predefined values
                        switch {
                            ignoreCase = true
                            cases {
                                on {
                                    shell {

                                    }        
                                }
                            }
                        }
                    }        
                }        
            }
        }
    }

    publishers {
        periodicPropertyPublisher { // publish properties value of specified nodes
            nodes = ["cpu"]
            every = "5s"
        }

        onMessagePublisher { // TODO: this is like a handler
            node = "cpu"
            property = "enabled"
            topic = "devices/esp8266/temperature"
            payloadConditionExpression  = "<18.0"  
        }
        onStartPublisher {
            nodes = ["cpu"]    
        }
        publisher { // node publisher configuration, only one type of publisher can be defined, use array-type to have multiple
            // how frequently publish payload
            // the idea is to support multiple when, invoked in specified order in this config
            when {
                periodic {
                    every = "5s"
                }
                
                // publish    
                onmqttmessage {
                    topic = "devices/esp8266/temperature"
                    // message payload, no idea how to implement something like `less than` `<` or others now 
                    message = "<18.0" 
                }

                // whenever process started, it is usually will be invoked before onconnect (see next)
                onstart {

                }

                // when connected (or re-connected) to mqtt broker    
                onconnect {

                }
            }

            // publish what value
            what {
                // shell handler, output of shell will be send as json {stdout: "", stderr: "", exitcode: ""}
                shell {
                    grep 'cpu ' /proc/stat | awk '{usage=($2+$4)*100/($2+$4+$5)} END {print usage "%"}'
                }

                // raspberry pi pin value    
                pi-pin {
                    pin = "04"
                }

            }
        }

    }
}

// VPN example device  
vpn {
    config {

    }
    nodes {
        openconnect {
           name = "VPN"
           type = "VPN Service"

           properties {
               users {
                   when {
                        periodic {
                            every = "5s"
                        }
                   }
               }
           } 
        }
    }
}