// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package tool

type Set_TBool struct {
	m *map[bool]bool
}

func NewSet_TBool() *Set_TBool {
	return &Set_TBool{
		m: &map[bool]bool{},
	}
}

func (s *Set_TBool) Add(o bool) {
	(*s.m)[o] = true
}
func (s *Set_TBool) Remove(o bool) {
	delete(*s.m, o)
}
func (s *Set_TBool) Has(o bool) bool {
	_, ok := (*s.m)[o]
	return ok
}
func (s *Set_TBool) Len() int {
	return len(*s.m)
}
func (s *Set_TBool) Clear() {
	s.m = &map[bool]bool{}
}
func (s *Set_TBool) IsEmpty() bool {
	return s.Len() == 0
}
func (s *Set_TBool) Slice() *[]bool {
	return MapKeys_KBoolVBool(s.m)
}

type Set_TInt struct {
	m *map[int]bool
}

func NewSet_TInt() *Set_TInt {
	return &Set_TInt{
		m: &map[int]bool{},
	}
}

func (s *Set_TInt) Add(o int) {
	(*s.m)[o] = true
}
func (s *Set_TInt) Remove(o int) {
	delete(*s.m, o)
}
func (s *Set_TInt) Has(o int) bool {
	_, ok := (*s.m)[o]
	return ok
}
func (s *Set_TInt) Len() int {
	return len(*s.m)
}
func (s *Set_TInt) Clear() {
	s.m = &map[int]bool{}
}
func (s *Set_TInt) IsEmpty() bool {
	return s.Len() == 0
}
func (s *Set_TInt) Slice() *[]int {
	return MapKeys_KIntVBool(s.m)
}

type Set_TString struct {
	m *map[string]bool
}

func NewSet_TString() *Set_TString {
	return &Set_TString{
		m: &map[string]bool{},
	}
}

func (s *Set_TString) Add(o string) {
	(*s.m)[o] = true
}
func (s *Set_TString) Remove(o string) {
	delete(*s.m, o)
}
func (s *Set_TString) Has(o string) bool {
	_, ok := (*s.m)[o]
	return ok
}
func (s *Set_TString) Len() int {
	return len(*s.m)
}
func (s *Set_TString) Clear() {
	s.m = &map[string]bool{}
}
func (s *Set_TString) IsEmpty() bool {
	return s.Len() == 0
}
func (s *Set_TString) Slice() *[]string {
	return MapKeys_KStringVBool(s.m)
}

type Set_TInterface struct {
	m *map[interface{}]bool
}

func NewSet_TInterface() *Set_TInterface {
	return &Set_TInterface{
		m: &map[interface{}]bool{},
	}
}

func (s *Set_TInterface) Add(o interface{}) {
	(*s.m)[o] = true
}
func (s *Set_TInterface) Remove(o interface{}) {
	delete(*s.m, o)
}
func (s *Set_TInterface) Has(o interface{}) bool {
	_, ok := (*s.m)[o]
	return ok
}
func (s *Set_TInterface) Len() int {
	return len(*s.m)
}
func (s *Set_TInterface) Clear() {
	s.m = &map[interface{}]bool{}
}
func (s *Set_TInterface) IsEmpty() bool {
	return s.Len() == 0
}
func (s *Set_TInterface) Slice() *[]interface{} {
	return MapKeys_KInterfaceVBool(s.m)
}
