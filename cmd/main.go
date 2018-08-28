package main

import (
	"e212/store"
	"os"
)

func usage() {
	println("Usage:", os.Args[0], " add mcc mnc country operator")
	println("      ", os.Args[0], " remove mcc mnc")
	os.Exit(2)

}

func main() {
	if len(os.Args) == 1 {
		usage()
	}
	err := store.Init("mccmnc.db")
	if err != nil {
		println("Failed to open store:", err.Error())
		os.Exit(1)
	}

	if os.Args[1] == "add" && len(os.Args) == 6 {
		e := store.NewE212Entry(os.Args[2], os.Args[3], os.Args[4], os.Args[5])

		err = store.E212Add(e)
		if err != nil {
			println("Failed to add:", err.Error())
			os.Exit(1)
		}
	} else if os.Args[1] == "remove" && len(os.Args) == 4 {
		e := store.MccMnc{Mcc: os.Args[2], Mnc: os.Args[3]}

		err = store.E212Remove(&e)
		if err != nil {
			println("Failed to remove:", err.Error())
			os.Exit(1)
		}
	} else {
		usage()
	}
}
