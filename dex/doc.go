/*
	Package dex provides key management for retrieval by indexes,
	extending Level through use of interfaces.
		const catName = "Michael"

		if err := dx.Store(Cat{catName}, 0); err != nil{
			panic(err)
		}

		var Michael Cat
		if err := dx.Retrieve(&Michael, 0); err != nil{
			panic(err)
		}
*/
package dex
