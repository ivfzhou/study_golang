// 原型模式是一种创建型设计模式，使你能够复制已有对象，而又无需使代码依赖它们所属的类。
// 如果你需要复制一些对象，同时又希望代码独立于这些对象所属的具体类，可以使用原型模式。
// 如果子类的区别仅在于其对象的初始化方式，那么你可以使用该模式来减少子类的数量。别人创建这些子类的目的可能是为了创建特定类型的对象。
package design_pattern_test

import (
	"fmt"
	"testing"
)

func TestPrototype(t *testing.T) {
	var file1 Inode = &File_{name: "File1"}
	var file2 Inode = &File_{name: "File2"}
	var file3 Inode = &File_{name: "File3"}

	var folder1 Inode = &Folder_{
		children: []Inode{file1},
		name:     "Folder1",
	}

	var folder2 Inode = &Folder_{
		children: []Inode{folder1, file2, file3},
		name:     "Folder2",
	}
	fmt.Println("Printing hierarchy for Folder2")
	folder2.Print("  ")

	cloneFolder := folder2.Clone()
	fmt.Println("Printing hierarchy for Clone Folder")
	cloneFolder.Print("  ")
}

// ===

type Inode interface {
	Print(string)
	Clone() Inode
}

// =

type File_ struct {
	name string
}

func (f *File_) Print(indentation string) {
	fmt.Println(indentation + f.name)
}

func (f *File_) Clone() Inode {
	return &File_{name: f.name + "_clone"}
}

// =

type Folder_ struct {
	children []Inode
	name     string
}

func (f *Folder_) Print(indentation string) {
	fmt.Println(indentation + f.name)
	for _, i := range f.children {
		i.Print(indentation + indentation)
	}
}

func (f *Folder_) Clone() Inode {
	cloneFolder := &Folder_{name: f.name + "_clone"}
	var tempChildren []Inode
	for _, i := range f.children {
		copy := i.Clone()
		tempChildren = append(tempChildren, copy)
	}
	cloneFolder.children = tempChildren
	return cloneFolder
}
