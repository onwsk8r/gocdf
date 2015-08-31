package netcdf

// #cgo LDFLAGS: -lnetcdf
// #include <stdlib.h>
// #include <netcdf.h>
import "C"
import (
	"strings"
	"unsafe"
)

type Variable struct {
	ncid      C.int
	Id        int
	Name      string
	Var_type  int
	Dimension map[string]Dimension //make this a slice or array of pointers
	Attribute map[string]Attribute
}

func (v *Variable) SetNcid(ncid C.int) {
	v.ncid = ncid
}

func (v *Variable) Get(idx int) interface{} {
	index := C.size_t(idx)

	switch {
	case v.Var_type == 2: //string
		str := C.CString(cstr)
		C.nc_get_var1_string(v.ncid, C.int(v.Id), &index, &str)
        val := strings.Trim(C.GoString(str), "`")
	    C.free(unsafe.Pointer(str))
        return val
	case v.Var_type == 3: //int16
        var nval C.int
		C.nc_get_var1_int(v.ncid, C.int(v.Id), &index, &nval)
        return int(nval)
	case v.Var_type == 4: //int32
        var nval C.int
		C.nc_get_var1_int(v.ncid, C.int(v.Id), &index, &nval)
        return int(nval)
	case v.Var_type == 5: //float32
        var nval C.float
		C.nc_get_var1_float(v.ncid, C.int(v.Id), &index, &nval)
        return float32(nval)
	case v.Var_type == 6: //float64
        var nval C.float
		C.nc_get_var1_float(v.ncid, C.int(v.Id), &index, &nval)
        return float64(nval)
	}
    return -9999
}