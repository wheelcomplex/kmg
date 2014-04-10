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
* bug: use "总数 :" as value of a map will generate a not valid yaml ,fixed in yaml_emitter_analyze_scalar():w = width(value[i]) //bug origin width(value[0])