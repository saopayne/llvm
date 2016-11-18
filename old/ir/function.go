package ir

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/llir/llvm/ir/internal/enc"
	"github.com/llir/llvm/ir/types"
	"github.com/mewkiz/pkg/errutil"
)

// A Function declaration specifies the name and type of a function. A function
// definition contains a set of basic blocks, interconnected by control flow
// instructions (such as br), which forms the nodes in a Control Flow Graph of
// the function [1,2].
//
//    [1]: http://llvm.org/docs/LangRef.html#functions
//    [2]: http://llvm.org/docs/LangRef.html#terminators
type Function struct {
	// Function name.
	name string
	// Function signature.
	sig *types.Func
	// Basic blocks of the function, or nil if function declaration.
	blocks []*BasicBlock
	// localID represents a counter used for assigning unique local IDs to
	// unnamed local variables and basic blocks.
	localID int
}

// NewFunction returns a new function based on the given name and function
// signature.
//
// The caller is responsible for invoking AssignIDs once all basic blocks have
// been added to the function.
func NewFunction(name string, sig *types.Func) *Function {
	return &Function{name: name, sig: sig}
}

// Name returns the name of the function.
func (f *Function) Name() string {
	return f.name
}

// Sig returns the function signature.
func (f *Function) Sig() *types.Func {
	return f.sig
}

// Type returns the type of the value.
func (f *Function) Type() types.Type {
	return f.sig
}

// TODO: Try to figure out a better API for accessing the function body.

// Blocks returns the basic blocks of the function.
func (f *Function) Blocks() []*BasicBlock {
	return f.blocks
}

// SetBlocks sets the function body to the given basic blocks and the parent
// function of each basic block to the given function, and assigns unique local
// IDs to unnamed basic blocks and local variable definitions.
func (f *Function) SetBlocks(blocks []*BasicBlock) error {
	// Set the function body of each basic block to the given function.
	for _, block := range blocks {
		block.SetParent(f)
	}
	f.blocks = blocks

	// Assign unique local IDs to unnamed basic blocks and local variable
	// definitions.
	if err := f.AssignIDs(); err != nil {
		return errutil.Err(err)
	}

	return nil
}

// AppendBlock appends the given basic block to the function.
//
// The caller is responsible for invoking AssignIDs once all basic blocks have
// been added to the function.
func (f *Function) AppendBlock(block *BasicBlock) {
	block.SetParent(f)
	f.blocks = append(f.blocks, block)
}

// AssignIDs assigns unique local IDs to unnamed basic blocks and local variable
// definitions.
func (f *Function) AssignIDs() error {
	for _, block := range f.blocks {
		if err := block.assignIDs(); err != nil {
			return errutil.Err(err)
		}
	}
	return nil
}

// String returns the string representation of the function declaration.
func (f *Function) String() string {
	// Function signature; e.g.
	//    "void @foo()"
	//    "i32 @printf(i8*, ...)"
	paramsBuf := new(bytes.Buffer)
	for i, param := range f.Sig().Params() {
		if i > 0 {
			paramsBuf.WriteString(", ")
		}
		paramsBuf.WriteString(param.String())
	}
	if f.Sig().IsVariadic() {
		if len(f.Sig().Params()) > 0 {
			paramsBuf.WriteString(", ")
		}
		paramsBuf.WriteString("...")
	}
	sig := fmt.Sprintf("%s %s(%s)", f.Sig().Result(), f.ValueString(), paramsBuf)

	// Function declaration; e.g.
	//    declare i32 @printf(i8*, ...)
	if f.blocks == nil {
		return fmt.Sprintf("declare %s", sig)
	}

	// Function definition; e.g.
	//     define i32 @main() {
	//        ret i32 42
	//     }
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "define %s {\n", sig)
	for _, block := range f.blocks {
		fmt.Fprintln(buf, block)
	}
	buf.WriteString("}")
	return buf.String()
}

// ValueString returns a string representation of the value.
func (f *Function) ValueString() string {
	return enc.Global(f.Name())
}

// nextID returns the next unique local ID of the given function, and increments
// the internal localID counter.
func (f *Function) nextID() string {
	id := strconv.Itoa(f.localID)
	f.localID++
	return id
}