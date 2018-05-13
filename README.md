# Time-based Unique Identifier (TUID)

[![GoDoc](https://godoc.org/github.com/semrekkers/tuid?status.svg)](https://godoc.org/github.com/semrekkers/tuid)

A TUID is a 64-bit unsigned integer with a time part and a random part. The time part is 42 bits long and contains the amount of milliseconds after a custom epoch (`2018-01-01 00:00:00.00 UTC`). This gives us a timespan of around 140 years (until ~2158). The other part contains 22 random bits. With this design, you can have 4,194,303 unique identifiers per millisecond, in theory. At this moment, there isn't any research done on the probability of a collision but it is considered fairly low.

This design and implementation is inspired by [this blog post](https://instagram-engineering.com/sharding-ids-at-instagram-1cf5a71e5a5c) about Sharding & IDs at Instagram.

## License

MIT License.
