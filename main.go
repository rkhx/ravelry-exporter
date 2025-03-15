package main

import (
	"fmt"
	"os"
)

func main() {
	username, err := getCurrentUsername()
	if err != nil {
		fmt.Println("Error fetching username:", err)
		os.Exit(1)
	}

	fmt.Println("Authenticated user:", username)

	bundles, err := getUserBundles(username)
	if err != nil {
		fmt.Println("Error fetching bundles:", err)
		os.Exit(1)
	}

	id := bundles[0].ID
	bla, err := getBundleContent(username, id)
	if err != nil {
		fmt.Println("Error fetching bundles:", err)
		os.Exit(1)
	}

	for _, item := range bla.BundledItems {
		aa, err := getBundleItem(item.Id)
		if err != nil {
			fmt.Println("Error fetching bundles:", err)
			os.Exit(1)
		}

		ss, err := GetPattern(aa.Item[0].Id)
		if err != nil {
			fmt.Println("Error fetching bundles:", err)
			os.Exit(1)
		}
		fmt.Println(ss)
	}
}
