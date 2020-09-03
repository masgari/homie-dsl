# Homie DSL

This is a Domain Specific Language (DSL) for [homie-go](github.com/masgari/homie-go/homie).

It is based on [HashiCorp Configuration Language](github.com/hashicorp/hcl)

Examples:

```hcl
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
```