#!/bin/bash
cat _poly/type.go | genny gen "Typ=bool" > type_poly.go
cat _poly/type_slice.poly | genny gen "Typ=bool,int,string,interface{}" > type_slice_poly.go
cat _poly/type_slice_ret.poly | genny gen "Typ=bool,int,string,interface{} Rtyp=bool,int,string,interface{}" > type_slice_ret_poly.go
cat _poly/type_map.poly | genny gen "TKey=bool,int,string,interface{} TValue=bool,int,string,interface{}" > type_map_poly.go
cat _poly/type_map_ret.poly | genny gen "TKey=bool,int,string,interface{} TValue=bool,int,string,interface{} RtKey=bool,int,string,interface{} RtValue=bool,int,string,interface{}" > type_map_ret_poly.go
cat _poly/type_collection.poly | genny gen "Typ=bool,int,string,interface{}" > type_collection_poly.go

