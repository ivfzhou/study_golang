// 组合模式是一种结构型设计模式，你可以使用它将对象组合成树状结构，并且能像使用独立对象一样使用它们。
// 如果你需要实现树状对象结构，可以使用组合模式。
// 如果你希望客户端代码以相同方式处理简单和复杂元素，可以使用该模式。
package design_pattern_test

import (
	"fmt"
	"testing"
)

func TestComposite(t *testing.T) {
	var file1 Component = &File{name: "File1"}
	var file2 Component = &File{name: "File2"}
	var file3 Component = &File{name: "File3"}

	folder1 := &Folder{
		name: "Folder1",
	}

	folder1.Add(file1)

	folder2 := &Folder{
		name: "Folder2",
	}
	folder2.Add(file2)
	folder2.Add(file3)
	folder2.Add(folder1)

	folder2.Search("rose")
}

// ===

type Component interface {
	Search(string)
}

// =

type File struct {
	name string
}

func (f *File) Search(keyword string) {
	fmt.Printf("Searching for keyword %s in file %s\n", keyword, f.name)
}

func (f *File) GetName() string {
	return f.name
}

// =

type Folder struct {
	components []Component
	name       string
}

func (f *Folder) Search(keyword string) {
	fmt.Printf("Serching recursively for keyword %s in folder %s\n", keyword, f.name)
	for _, composite := range f.components {
		composite.Search(keyword)
	}
}

func (f *Folder) Add(c Component) {
	f.components = append(f.components, c)
}
