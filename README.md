This is a refactor of an old Go Lambda that creates an API that will return a visitor's geolocation based on IP address.
It now uses the Aegis deploy tool and framework, which uses native Go Lambdas now instead of needing a shim.

### Prerequisites

 - An AWS account (and configured credentials for CLI)
 - [Maxmind's GeoLite2 City database.](http://dev.maxmind.com/geoip/geoip2/geolite2/)
 - [Aegis](https://github.com/tmaiaroto/aegis) deploy tool and framework
 - [go-bindata](https://github.com/jteeuwen/go-bindata) 

The filename will likely be ```GeoLite2-City.mmdb```.


## Instructions

Assuming you have AWS credentials configured and you've got the `aegis` and `go-bindata` binaries
in your path ready to use. You can run the following to retrieve the geoip data set, build, and deploy:

```
curl -o GeoLite2-City.tar.gz http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.tar.gz

tar -xvf GeoLite2-City.tar.gz

mv GeoLite2-City_20180327/GeoLite2-City.mmdb GeoLite2-City.mmdb

go-bindata GeoLite2-City.mmdb

aegis deploy
```

It should build and deploy the Lambda. It may take a little bit due to the size of the database file.
It took my 2min 12s to build and deploy (so says my zsh).

------

You can see the old way we had to use Go in AWS Lambda by looking at the `legacy-shim` branch 
[and here's the original article.](https://medium.com/@shift8creative/go-amazon-lambda-7e95a147cec8#.ab93bgu8s)