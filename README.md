# Description

I2P is a tool for converting insomina data (Insomina v5) to postman (Postman v2.1) and vice versa

## Build


```bash
./build.sh
```

## Usage
Insomnia to postman
```
i2p convert --input-file insomnia.yaml --output-file postman.json
```
Postman to insomnia
```
i2p convert --input-file postman.json --output-file insomnia.yaml
```

## Supported commands

```bash
Usage:
  i2p convert [--input-file FILE] [--output-file FILE]

  i2p help
  i2p version

Options:
  --input-file FILE   Specify the input file (default: insomnia.yaml)
  --output-file FILE  Specify the output file (default: postman_collection.json)
```
