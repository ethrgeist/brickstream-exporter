# XML Version Examples

There seem to be 5 different versions of XML output that a Brickstream device produces.

The only noticeable difference seem to be between the first two versions.

Version 1 lacks the following fields:

- `Version` - XML version
- `TransmitTime` - A unix timestamp of when the XML was transmitted
- `TimezoneName` - String representing the timezone of the device
- `SwRelease` - Software release version of the device

Versions 2 to 5 are identical field-wise.

It's not entirely clear what the difference is between versions 2 to 5, but it seems that they are just different
releases of the same software.

## Notes

- `Timezone` contains the UTC offset in hours, e.g. `-5` for GMT-5 or `+1` for GMT+1.
- `DST` contains the daylight saving time offset in hours, e.g. `-1` for GMT-5 with DST or `0` for GMT+1 without DST.