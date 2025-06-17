# brickstream-exporter

This is a Prometheus exporter for Brickstream devices, it can be used as a HTTP target for Brickstream to send XML data to.
Since this exporter does not directly query the Brickstream, it is not realtime, the report interval is defined via
configuration of the Brickstream device itself.

It also can store data points in a SQLite database, which can be used to query historical data or export into other formats
like CSV or JSON.

## Multizone Support

The exporter was only tested with a single zone per device due missing licenses for multiple zones. While the exporter
takes multiple outputs in the XML report into account, it is unclear, if this actually works. If you find issues with this,
please open an [issue on GitHub](https://github.com/ethrgeist/brickstream-exporter/issues), ideally with a sample XML
report that fails.

## XML Samples

For examples of the XML output, see the [XML Version Examples](examples/xml-versions/README.md).