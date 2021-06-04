# Images Loader

## Details

Utility used to download some images from the https://www.pexels.com website automatically.

## Usage and options


```
go run main.go [options]
```

where options can be:

1. Request options

* -n int

        Number of results per page (default 20)

* -o string

        Orientation (landscape, portrait or square) (default "landscape")

* -p int

        Page to display (default 1)

* -q string

        Query (people, nature, ...) (default "people")

* -z string

        Minimum photo size (large (24MP), medium(12MP), small(4MP)) (default "small")

* -s string

        An assortment of different image sizes that can be used to display the photos. Authorized values are:

        Original: The image without any size changes. It will be the same as the width and height attributes.

        Large2x: The image resized W 940px X H 650px DPR 2

        Large: The image resized to W 940px X H 650px DPR 1

        Medium: The image scaled proportionally so that it's new height is 350px

        Small: The image scaled proportionally so that it's new height is 130px

        Portrait: The image cropped to W 800px X H 1200px

        Landscape: The image cropped to W 1200px X H 627px

        Tiny: The image cropped to W 280px X H 200px

* -l string

        log level (debug, info, warn, error) (default "info")

* -c config.yml

        use a YAML configuration file with all options already defined. If the config file is defined, then the 
        command line options are ignored.

All command line options are optional and have default values. If you use the config file, all options MUST be defined.

## Example

go run main.go parseJson.go request.go checks.go fileutil.go config.go -c config.yml

go run main.go parseJson.go request.go checks.go fileutil.go config.go

go run main.go parseJson.go request.go checks.go fileutil.go config.go -q people -p 2 -n 80 -s Small -z small -f /tmp -l warn
