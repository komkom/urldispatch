package urldispatch

import (
	"fmt"
	"net/url"
	"testing"
)

type link struct {
	dispIdx int
	urlIdx  int
}

func TestTest(t *testing.T) {

	dispatchers := []string{"test/:name1/:name2/xx/:k1/:key.../part3/:v1"}

	urls := []string{"https://test.com/test/value/value2/xx/vk1/1/2/3/4/5/part3/somevalue"}

	links := []link{link{dispIdx: 0, urlIdx: 0}}

	for _, l := range links {

		u, err := url.Parse(urls[l.dispIdx])
		if err != nil {
			panic(err)
		}

		du, err := url.Parse(dispatchers[l.urlIdx])
		if err != nil {
			panic(err)
		}

		/*
			at := []int{0}

			ptr := &at[0]
			fmt.Printf("%v\n", ptr)
		*/

		d := Dispatcher{}
		err = d.AddRoute(dispatchURL)
		if err != nil {
			panic(err)
		}

		oa, err := d.Dispatch(u)
		if err != nil {
			panic(err)
		}

		fmt.Printf("__args %v\n", oa.ar)

	}
}
