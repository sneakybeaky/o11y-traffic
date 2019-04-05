# o11y target generater for Vegeta

Generates targets in the [Vegeta generated targets](https://github.com/tsenart/vegeta#usage-generated-targets) json format.

For each image file found a POST request with a mulitpart body consisting of

- `file` : the file contents
- `name` : a randomized string


## Usage manual

```console
Usage:
  o11y-traffic [flags]

Flags:
  -d, --directory string   Directory to walk
  -f, --forever            Run forever (default true)
  -h, --help               help for imggen
  -u, --url string         URL to POST form data to (default "http://localhost:8080/api/images")
```

### Mandatory flags

##### `--directory` or `-d` 

The directory to walk looking for image files identified by file extension : 

- .jpeg / .jpg
- .tiff / .tif
- .gif
- .png
- .svg

The sequence of the found files will be randomized before being encoded as targets.

### Optional flags

##### `--forever` or `-f` 

When false the files found will be sent just once. When these files have all been sent the program will exit.

When set to true the files will be continually sent. Each cycle the sequence of the images will be randomized.

##### `--url` or `-u` 

The URL where the POST will be made to

##### `--help` or `-h` 

Show usage help

## Usage: feed targets into vegeta

You will need to [install vegeta](https://github.com/tsenart/vegeta/blob/master/README.md#install).

To find all images under `~/tmp/jpg` and send to `http://localhost:8080/api/images` at a rate of 1 per second for 30 seconds :

```console
o11y-traffic -d ~/tmp/jpg | \
  vegeta attack -rate=1/s -lazy -format=json -duration=30s | \
  tee results.bin | \
  vegeta report
```

