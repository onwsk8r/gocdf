GoCDF
==

I wrote this lightweight library to wrap [libnetcdf](http://www.unidata.ucar.edu/software/netcdf/) because I needed to read NetCDF files in Go.
Currently it's just a friendly interface to the library that affords me the ability to iterate through the data and do stuff.

#### BUT IN THE FUTURE
I absolutely intend to turn this into a Pure Go NetCDF library.

### Usage
```go
file := Read('/path/to/netCDF');
file.parse(); //now we have all of the dims, vars, and attrs
defer nc.Close()

netCdfDimension := file.getDimensionByName('theDim');
otherDimension := file.getDimensionById(15);
//...then the Dimension struct has no methods so I guess we're done

length := nc.Dimension["recNum"].Length //Supposing you have a dimension called recNum that contains the number of records
width := len(nc.Variable)
for i := 0; i < length; i++ {
    row := make(map[string]interface{}, width)
    for _, value := range nc.Variable {
        row[value.Name] = value.Get(i)
    }
    doStuff(row)
}
```

### I KNOW IT'S ROUGH
It's supposed to accomplish one goal - let me run that loop.
Which it does. It runs through about 250 megs of HDF5 on one Xeon core in about six seconds.
I'll be putting in additional development time (or you can) later.
