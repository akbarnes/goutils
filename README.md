# GoUtils
This is a collection of command-line utilities in go

## GoCalc
A 4-level stack RPN calculator

## GoCoffee
Calculate how many tablespoons of coffee to put per cup of water

## GoTimer 
Command line timer

## GoSun
Sunrise/sunset calculator. Currently not giving accurate results

## GoWeather
Check the weather

## GoTimestamp
Print a time or date stamp in a format that can be used in file or directory names

## GoConvert
Convert units

## GoRes
Resistor color code calculator

## GoHash
Calculate the SHA256 hash of a file

## GoVer
This is a much simplified version of [dupver](https://github.com/akbarnes/dupver/) that uses a local repository and file-level deduplication. There are only three commands `commit`, `log`, and `checkout`. 

This has been moved to its own repo: [https://github.com/akbarnes/gover](https://github.com/akbarnes/gover)


To build:
```
go mod init goutils
go get github.com/bmatcuk/doublestar/v4
go install gover.go
```

### Commit
The `-msg` or `-m` message flag is optional, as is the `-commit` or `-ci` flag as commiting is the default action
`gover -commit -msg 'a message' file1 file2 file3`
`gover -ci -msg 'a message' file1 file2 file3`
`gover -m  'a message 'file1 file2 file3`
`gover file1 file2 file3`

### Log
This takes the optional `-json` or `-j` argument to output json for use with object shells. To list all the snapshots:
`gover -log`
`gover -json -log`
`gover -j -log`
`gover -l`

To list the files in a particular snapshot:
`gover -log snapshot_time`
`gover -log -json snapshot_time`

## Checkout
This takes an optional argument to specify an output folder. To checkout a snapshot:
`gover -checkout snapshot_time`
`gover -co snapshot_time`
`gover -out output_folder -co snapshot_time`
`gover -o output_folder -co snapshot_time`

## GoRand
Print a random text string 

## GoRandText
Create a random text file

# GoDiff
Simple program to diff two files. I wrote this to have consistent diff
behavior across operating systems. Usage:

`godiff file1 file2`

It will return `equal` if the files are equal and `different` otherwise.
If the `-errors` flag is enabled it will return `error` if one file does not
exist or is unreadable and `different` otherwise.

If the `json` flag is enabled it will return `false` if the files are equal
and `true` otherwise. If the `-errors` flag is also enabled it will return 
`null` if one file does not exist or is unreadable and `true` otherwise.

Finally, if the `-human` flag is enabled this will toggle verbose human-readable
output which includes automatically reporting if there was an error reading files.
