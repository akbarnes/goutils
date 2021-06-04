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
`gover -co snapshot_time -out output_folder`
`gover -co snapshot_time -o output_folder`

