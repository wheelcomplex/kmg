kmgYaml
=================
origin code come from "launchpad.net/goyaml"

fix follow problem:
* struct key name will not change when Unmarshal and Marshal by default(goyaml will change them to lowercase).
* chinese string will not Marshal to "\uxxxx"
* value 1 Unmarshal to float64 problem