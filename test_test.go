package urldispatch

import (
	"fmt"
	"net/url"
	"testing"
)

type targs struct {
	paramCount  int
	paramNames  []string
	paramValues []string
	arrayCount  int
}

type turl struct {
	url string
}

type tunit struct {
	dispatch string
	tus      []turl
}

func TestTest(t *testing.T) {

	units := []tunit{tunit{dispatch: "start/test/:name1/:name2/xx/:k1/:key.../part3/:v1?qk", tus: []turl{turl{url: "https://test.com/start/test/value/xx/vk1/1/2/3/4/5/part3/somevalue?qk=1111"}}}}

	d := Dispatcher{}

	for _, u := range units {

		du, err := url.Parse(u.dispatch)
		if err != nil {
			panic(err)
		}

		err = d.AddRoute(du)
		if err != nil {
			panic(err)
		}

		for _, tu := range u.tus {

			u, err := url.Parse(tu.url)
			if err != nil {
				panic(err)
			}

			oa, err := d.Dispatch(u)
			if err != nil {
				panic(err)
			}

			/*
				fmt.Printf("__args %v\n", oa.ar)
				fmt.Printf("__pc %v\n", oa.ParamCount())

				for i := 0; i < oa.ParamCount(); i++ {
					val, err := oa.Value(i)
					if err != nil {
						fmt.Println("error: " + err.Error())
					}
					fmt.Printf("__value %v\n", val)
				}
			*/

			pNames := []string{"name1", "name2", "k1", "v1", "qk", "ii"}

			for _, pn := range pNames {
				value, err := oa.ParamWithName(pn)
				if err != nil {
					fmt.Println("error: " + err.Error())
				}

				fmt.Println("value " + value)
			}

			aNames := []string{"key"}

			for _, pn := range aNames {
				arr, err := oa.ArrayWithName(pn)
				if err != nil {
					fmt.Println("error: " + err.Error())
				}

				fmt.Printf("array: %v\n", arr)
			}

		}
	}
}
