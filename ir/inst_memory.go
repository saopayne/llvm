// === [ Memory instructions ] =================================================
//
// References:
//    http://llvm.org/docs/LangRef.html#memory-access-and-addressing-operations

package ir

import (
	"bytes"
	"fmt"

	"github.com/llir/llvm/internal/enc"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// --- [ alloca ] --------------------------------------------------------------

// InstAlloca represents an alloca instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#alloca-instruction
type InstAlloca struct {
	// Parent basic block.
	parent *BasicBlock
	// Name of the local variable associated with the instruction.
	name string
	// Type of the instruction.
	typ *types.PointerType
	// Element type.
	elem types.Type
	// Number of elements; or nil if one element.
	nelems value.Value
}

// NewAlloca returns a new alloca instruction based on the given element type.
func NewAlloca(elem types.Type) *InstAlloca {
	typ := types.NewPointer(elem)
	return &InstAlloca{typ: typ, elem: elem}
}

// Type returns the type of the instruction.
func (inst *InstAlloca) Type() types.Type {
	return inst.typ
}

// Ident returns the identifier associated with the instruction.
func (inst *InstAlloca) Ident() string {
	return enc.Local(inst.name)
}

// Name returns the name of the local variable associated with the instruction.
func (inst *InstAlloca) Name() string {
	return inst.name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstAlloca) SetName(name string) {
	inst.name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstAlloca) String() string {
	if nelems, ok := inst.NElems(); ok {
		return fmt.Sprintf("%s = alloca %s, %s %s",
			inst.Ident(),
			inst.ElemType(),
			nelems.Type(),
			nelems.Ident())
	}
	return fmt.Sprintf("%s = alloca %s",
		inst.Ident(),
		inst.ElemType())
}

// Parent returns the parent basic block of the instruction.
func (inst *InstAlloca) Parent() *BasicBlock {
	return inst.parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstAlloca) SetParent(parent *BasicBlock) {
	inst.parent = parent
}

// ElemType returns the element type of the alloca instruction.
func (inst *InstAlloca) ElemType() types.Type {
	return inst.elem
}

// NElems returns the number of elements of the alloca instruction and a boolean
// indicating if the number of elements were present.
func (inst *InstAlloca) NElems() (value.Value, bool) {
	if inst.nelems != nil {
		return inst.nelems, true
	}
	return nil, false
}

// SetNElems sets the number of elements of the alloca instruction.
func (inst *InstAlloca) SetNElems(nelems value.Value) {
	inst.nelems = nelems
}

// --- [ load ] ----------------------------------------------------------------

// InstLoad represents a load instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#load-instruction
type InstLoad struct {
	// Parent basic block.
	parent *BasicBlock
	// Name of the local variable associated with the instruction.
	name string
	// Type of the instruction.
	typ types.Type
	// Source address.
	src value.Value
}

// NewLoad returns a new load instruction based on the given source address.
func NewLoad(src value.Value) *InstLoad {
	t, ok := src.Type().(*types.PointerType)
	if !ok {
		panic(fmt.Sprintf("invalid source address type; expected *types.PointerType, got %T", src.Type()))
	}
	return &InstLoad{typ: t.Elem(), src: src}
}

// Type returns the type of the instruction.
func (inst *InstLoad) Type() types.Type {
	return inst.typ
}

// Ident returns the identifier associated with the instruction.
func (inst *InstLoad) Ident() string {
	return enc.Local(inst.name)
}

// Name returns the name of the local variable associated with the instruction.
func (inst *InstLoad) Name() string {
	return inst.name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstLoad) SetName(name string) {
	inst.name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstLoad) String() string {
	src := inst.Src()
	return fmt.Sprintf("%s = load %s, %s %s",
		inst.Ident(),
		inst.Type(),
		src.Type(),
		src.Ident())
}

// Parent returns the parent basic block of the instruction.
func (inst *InstLoad) Parent() *BasicBlock {
	return inst.parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstLoad) SetParent(parent *BasicBlock) {
	inst.parent = parent
}

// Src returns the source address of the load instruction.
func (inst *InstLoad) Src() value.Value {
	return inst.src
}

// SetSrc sets the source address of the load instruction.
func (inst *InstLoad) SetSrc(src value.Value) {
	inst.src = src
}

// --- [ store ] ---------------------------------------------------------------

// InstStore represents a store instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#store-instruction
type InstStore struct {
	// Parent basic block.
	parent *BasicBlock
	// Source value.
	src value.Value
	// Destination address.
	dst value.Value
}

// NewStore returns a new store instruction based on the given source value and
// destination address.
func NewStore(src, dst value.Value) *InstStore {
	return &InstStore{src: src, dst: dst}
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstStore) String() string {
	src, dst := inst.Src(), inst.Dst()
	return fmt.Sprintf("store %s %s, %s %s",
		src.Type(),
		src.Ident(),
		dst.Type(),
		dst.Ident())
}

// Parent returns the parent basic block of the instruction.
func (inst *InstStore) Parent() *BasicBlock {
	return inst.parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstStore) SetParent(parent *BasicBlock) {
	inst.parent = parent
}

// Src returns the source value of the store instruction.
func (inst *InstStore) Src() value.Value {
	return inst.src
}

// SetSrc sets the source value of the store instruction.
func (inst *InstStore) SetSrc(src value.Value) {
	inst.src = src
}

// Dst returns the destination address of the store instruction.
func (inst *InstStore) Dst() value.Value {
	return inst.dst
}

// SetDst sets the destination address of the store instruction.
func (inst *InstStore) SetDst(dst value.Value) {
	inst.dst = dst
}

// --- [ fence ] ---------------------------------------------------------------

// --- [ cmpxchg ] -------------------------------------------------------------

// --- [ atomicrmw ] -----------------------------------------------------------

// --- [ getelementptr ] -------------------------------------------------------

// InstGetElementPtr represents a getelementptr instruction.
//
// References:
//    http://llvm.org/docs/LangRef.html#getelementptr-instruction
type InstGetElementPtr struct {
	// Parent basic block.
	parent *BasicBlock
	// Name of the local variable associated with the instruction.
	name string
	// Type of the instruction.
	typ types.Type
	// Source address element type.
	elem types.Type
	// Source address.
	src value.Value
	// Element indices.
	indices []value.Value
}

// NewGetElementPtr returns a new getelementptr instruction based on the given
// source address and element indices.
func NewGetElementPtr(src value.Value, indices ...value.Value) *InstGetElementPtr {
	srcType, ok := src.Type().(*types.PointerType)
	if !ok {
		panic(fmt.Sprintf("invalid source address type; expected *types.PointerType, got %T", src.Type()))
	}
	elem := srcType.Elem()
	e := elem
	for i, index := range indices {
		if i == 0 {
			// Ignore checking the 0th index as it simply follows the pointer of
			// src.
			//
			// ref: http://llvm.org/docs/GetElementPtr.html#why-is-the-extra-0-index-required
			continue
		}
		if t, ok := e.(*types.NamedType); ok {
			e, ok = t.Def()
			if !ok {
				panic(fmt.Sprintf("invalid named type %q; expected underlying type definition, got nil", t))
			}
		}
		switch t := e.(type) {
		case *types.PointerType:
			// ref: http://llvm.org/docs/GetElementPtr.html#what-is-dereferenced-by-gep
			panic("unable to index into element of pointer type; for more information, see http://llvm.org/docs/GetElementPtr.html#what-is-dereferenced-by-gep")
		case *types.ArrayType:
			e = t.Elem()
		case *types.StructType:
			idx, ok := index.(*constant.Int)
			if !ok {
				panic(fmt.Sprintf("invalid index type for structure element; expected *constant.Int, got %T", index))
			}
			e = t.Fields()[idx.Int64()]
		default:
			panic(fmt.Sprintf("support for indexing element type %T not yet implemented", e))
		}
	}
	typ := types.NewPointer(e)
	return &InstGetElementPtr{typ: typ, elem: elem, src: src, indices: indices}
}

// Type returns the type of the instruction.
func (inst *InstGetElementPtr) Type() types.Type {
	return inst.typ
}

// Ident returns the identifier associated with the instruction.
func (inst *InstGetElementPtr) Ident() string {
	return enc.Local(inst.name)
}

// Name returns the name of the local variable associated with the instruction.
func (inst *InstGetElementPtr) Name() string {
	return inst.name
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *InstGetElementPtr) SetName(name string) {
	inst.name = name
}

// String returns the LLVM syntax representation of the instruction.
func (inst *InstGetElementPtr) String() string {
	buf := &bytes.Buffer{}
	src := inst.Src()
	fmt.Fprintf(buf, "%s = getelementptr %s, %s %s",
		inst.Ident(),
		inst.elem,
		src.Type(),
		src.Ident())
	for _, index := range inst.Indices() {
		fmt.Fprintf(buf, ", %s %s",
			index.Type(),
			index.Ident())
	}
	return buf.String()
}

// Parent returns the parent basic block of the instruction.
func (inst *InstGetElementPtr) Parent() *BasicBlock {
	return inst.parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *InstGetElementPtr) SetParent(parent *BasicBlock) {
	inst.parent = parent
}

// Src returns the source address of the getelementptr instruction.
func (inst *InstGetElementPtr) Src() value.Value {
	return inst.src
}

// SetSrc sets the source address of the getelementptr instruction.
func (inst *InstGetElementPtr) SetSrc(src value.Value) {
	inst.src = src
}

// Indices returns the element indices of the getelementptr instruction.
func (inst *InstGetElementPtr) Indices() []value.Value {
	return inst.indices
}

// SetIndices sets the element indices of the getelementptr instruction.
func (inst *InstGetElementPtr) SetIndices(indices []value.Value) {
	inst.indices = indices
}
