package main

import comp "me.compiler/compiler"

func main() {
	c := comp.NewCompiler()
	c = c.SetData(make(map[string]any)).SetTxt("lorem ipsum")
	r, _, _ := c.Compile().GetResponse()
	print(r)
}
