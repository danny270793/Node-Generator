# Node Generator

Generate a node project scaffold

## Installation

### From Github releases page

Go to [Release page](https://github.com/danny270793/Node-Generator/releases) then download the binary which fits your environment

### From terminal

Get the last versión available on github

```bash
LAST_VERSION=$(curl https://api.github.com/repos/danny270793/Node-Generator/releases/latest | grep tag_name | cut -d '"' -f 4)
```

Download the last version directly to the binaries folder

For Linux (linux):

```bash
curl -L https://github.com/danny270793/Node-Generator/releases/download/${LAST_VERSION}/NodeGenerator_${LAST_VERSION}_linux_amd64.tar.gz -o ./nodegenerator.tar.gz
```

For MacOS (darwin):

```bash
curl -L https://github.com/danny270793/Node-Generator/releases/download/${LAST_VERSION}/NodeGenerator_${LAST_VERSION}_darwin_amd64.tar.gz -o ./nodegenerator.tar.gz
```

Untar the downloaded file

```bash
tar -xvf ./nodegenerator.tar.gz
```

Then copy the binary to the binaries folder

```bash
sudo cp ./NodeGenerator /usr/local/bin/nodegenerator
```

Make it executable the binary

```bash
sudo chmod +x /usr/local/bin/nodegenerator
```

```bash
nodegenerator --version
```

## Ussage

Run the binary and pass the path to the folders where you want to search for git projects

```bash
nodegenerator create /path/to/project
```

## Follow me

- [Youtube](https://www.youtube.com/channel/UC5MAQWU2s2VESTXaUo-ysgg)
- [Github](https://www.github.com/danny270793/)
- [LinkedIn](https://www.linkedin.com/in/danny270793)

## LICENSE

Licensed under the [MIT](license.md) License

## Version

NodeGenerator version 1.0.0

Last update 05/12/2023
