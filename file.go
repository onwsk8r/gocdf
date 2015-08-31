package netcdf

// #cgo LDFLAGS: -lnetcdf
// #include <stdlib.h>
// #include <netcdf.h>
import "C"
import (
    "unsafe"
    "fmt"
)

const (
	// Dummy string for converting between C and go string formats.
	// If it isn't long enough, we end up with memory errors.
	cstr = "`````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````"
)

// Helper functions to open files
func Read(name string) *File {
    var file File
    var cname = C.CString(name)
    file.Open(cname, C.NC_NOWRITE)
    C.free(unsafe.Pointer(cname))
    return &file
}

func Write(name string) *File {
    var file File
    var cname = C.CString(name)
    file.Open(cname, C.NC_WRITE)
    C.free(unsafe.Pointer(cname))
    return &file
}

// File struct
type File struct {
    ncid      C.int
    Dimension map[string]Dimension
    Variable  []Variable
    Attribute []Attribute
}

// Modern looking function
// Fixme throw a warning if any dims already exist
func (f *File) Parse() {
    var ndims, nvars, natts, unlimdimidp C.int
    cerr, err := C.nc_inq(f.ncid, &ndims, &nvars, &natts, &unlimdimidp)
    if err != nil || cerr != 0 {
        fmt.Println(err,cerr);
    }
    
    // Get the dimensions - make this concurrent
    fmt.Printf("Found %v dimensions\n", ndims)
    f.Dimension = make(map[string]Dimension, int(ndims))
    dim := C.int(0)
    for dim < ndims {
        var name = C.CString(cstr)
        var length C.size_t
        C.nc_inq_dim(f.ncid, dim, name, &length)
        
        var d = Dimension{Id: int(dim)}
        d.Length = int(length)
        d.Name  = C.GoString(name)
        f.Dimension[d.Name] = d
        C.free(unsafe.Pointer(name))
        dim++
    }
    fmt.Println("Successfully parsed dimensions")
    
    // Get the variables
    fmt.Printf("Found %v variables\n", nvars)
    f.Variable = make([]Variable, int(nvars))
    varid := C.int(0)
    for varid < nvars {
        var name = C.CString(cstr)
        var var_type C.nc_type
        var ndims C.int
        var dimids C.int
        var natts C.int
        C.nc_inq_var(f.ncid, C.int(varid), name, &var_type, &ndims, &dimids, &natts)
        
        var v Variable
        v.SetNcid(f.ncid)
        v.Id = int(varid)
        v.Name = C.GoString(name)
        v.Var_type = int(var_type)
        f.Variable[int(varid)] = v
        C.free(unsafe.Pointer(name))
        varid++
    }
    fmt.Println("Successfully parsed variables")
}

// Goroutines

// Functions you don't care about
func (f *File) Open(path *C.char, mode C.int) {
    var chunksizehint C.size_t
    var ncid C.int
    
    cerr, err := C.nc__open(path,mode,&chunksizehint,&ncid)
    if err != nil || cerr != 0 {
        fmt.Println(err,cerr);
    }
    
    f.ncid = ncid
}

func (n *File) Close() {
    ncid := C.int(n.ncid)
    C.nc_close(ncid)
}

// Primitives to get a single dimension FIXME - return a pointer to a dimension under the file
func (f File) getDimensionByName(name string) Dimension {
    
    var ncid  = C.int(f.ncid)
    var id      C.int
    var cname = C.CString(name)
    var length  C.size_t
    
    C.nc_inq_dimid(ncid, cname, &id)
    C.nc_inq_dimlen(ncid, id, &length)
    
    var d = Dimension{ Name: name}
    d.Id     = int(id)
    d.Length = int(length)
    
    C.free(unsafe.Pointer(cname))
    return d
}

func (f File) getDimensionById(id int) Dimension {
    var ncid   = C.int(f.ncid)
    var cid    = C.int(id)
    var name   = C.CString(cstr)
    var length   C.size_t
    
    C.nc_inq_dim(ncid, cid, name, &length)
    var d Dimension
    d.Id     = id
    d.Name   = C.GoString(name)
    d.Length = int(length)
    
    C.free(unsafe.Pointer(name))
    return d
}