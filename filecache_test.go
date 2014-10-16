package filecache_test

import (
	"fmt"
	"github.com/DannyBen/filecache"
)

func Example() {
	// Get a cache handler with default (temp) directory as
	// the cache directory, and a one hour cache life
	handler := filecache.Handler{Life: 1}

	// Data to store in cache
	data := []byte("Joey doesn't share food")

	// Store the data in the cache
	handler.Set("testkey", data)

	// Retrieve the data from the cache
	r := handler.Get("testkey")

	// Show the result
	fmt.Println(string(r))

	// Output:
	// Joey doesn't share food
}
