package reflect2

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

// typelinks2 for 1.7 ~
//
//go:linkname typelinks2 reflect.typelinks
func typelinks2() (sections []unsafe.Pointer, offset [][]int32)

// initOnce guards initialization of types and packages
var initOnce sync.Once

var types = map[string]reflect.Type{}
var packages = map[string]map[string]reflect.Type{}

func Register[T any]() {
	t := reflect.TypeFor[T]()
	types[fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())] = t
}

func Register2(in any) {
	t := reflect.TypeOf(in)
	types[fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())] = t
}

// discoverTypes initializes types and packages
func discoverTypes() {
	loadGoTypes()
}

func loadGoTypes() {
	var obj interface{} = reflect.TypeOf(0)
	sections, offset := typelinks2()
	for i, offs := range offset {
		rodata := sections[i]
		for _, off := range offs {
			(*emptyInterface)(unsafe.Pointer(&obj)).word = resolveTypeOff(unsafe.Pointer(rodata), off)
			typ := obj.(reflect.Type)
			if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
				loadedType := typ.Elem()
				pkgTypes := packages[loadedType.PkgPath()]
				if pkgTypes == nil {
					pkgTypes = map[string]reflect.Type{}
					packages[loadedType.PkgPath()] = pkgTypes
				}
				types[loadedType.String()] = loadedType
				pkgTypes[loadedType.Name()] = loadedType
			}
		}
	}
}

type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

// TypeByName return the type by its name, just like Class.forName in java
func TypeByName(typeName string) (Type, error) {
	initOnce.Do(discoverTypes)
	t, ok := types[typeName]
	if !ok {
		return nil, errors.New("invalid type name")
	}
	return Type2(t), nil
}

// TypeByPackageName return the type by its package and name
func TypeByPackageName(pkgPath string, name string) (Type, error) {
	initOnce.Do(discoverTypes)
	pkgTypes, ok := packages[pkgPath]
	if !ok || pkgTypes == nil {
		return nil, errors.New("invalid package name")
	}
	pkg, ok := pkgTypes[name]
	if !ok {
		return nil, errors.New("invalid type name")
	}
	return Type2(pkg), nil
}
