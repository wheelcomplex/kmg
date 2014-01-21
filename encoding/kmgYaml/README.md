kmgYaml
=================
origin code come from "launchpad.net/goyaml"

add following bugfix:
* improve: use github for source control.
* improve: struct key name will not change when Unmarshal and Marshal by default(goyaml will change them to lowercase).
* improve: chinese string will not Marshal to "\uxxxx"
* bug:"1" can not unmarshal to float64 problem
* improve: Can unmarshal array problem now.
* improve: Can unmarshal/marshal time.Time correct.