YT - YAML Tool
===================

A command line hack in golang to pull pieces of a YAML document apart
using JSONPath.

[vmware-labs/yaml-jsonpath](https://github.com/vmware-labs/yaml-jsonpath)
does the hard work.

Example
--------

Given the yaml document:
```yaml
---
manufacturer:
  european:
    bmw:
      - model: R1250
        cylinders: 2
        displacement: 1254 cc
        hp: 136 bhp
        torque: 105 lb-ft
        curb weight: 479 lbs
    ktm:
      - model: 390 Adventure
        cylinders: 1
        displacement: 373 cc
        hp: 43 bhp
        torque: 27 lb-ft
        curb weight: 379 lbs
      - model: 790 Adventure
        cylinders: 2
        displacement: 799 cc
        hp: 95 bhp
        torque: 65 lb-ft
        curb weight: 417 lbs
```

Pull out the models:
```bash
$ ./yt test/mc.yaml '..model'
R1250
390 Adventure
790 Adventure
```

Or the models for a manufacturer:
```bash
$ ./yt test/mc.yaml '.manufacturer.european.bmw[*].model'
R1250
```
