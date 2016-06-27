package urldispatch

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

type targs struct {
	paramCount  int
	paramNames  []string
	paramValues []string
	arrayCount  int
}

type turl struct {
	url              string
	expectedParams   map[string]string
	unexpectedParams []string
	expectedArrays   map[string][]string
}

type tunit struct {
	dispatch string
	tus      []turl
}

func TestTest(t *testing.T) {

	units := []tunit{
		tunit{
			dispatch: "start/test/:name1/:name2/xx/:k1/:key.../part3/:v1?qk",
			tus: []turl{turl{
				url:              "https://test.com/start/test/value/xx/vk1/1/2/3/4/5/part3/somevalue?qk=1111",
				expectedParams:   map[string]string{"name1": "value", "k1": "vk1", "v1": "somevalue", "qk": "1111"},
				unexpectedParams: []string{"name2"},
				expectedArrays:   map[string][]string{"key": []string{"1", "2", "3", "4", "5"}}}}},
		tunit{
			dispatch: "end/test/:name1/:name2/:keys...",
			tus: []turl{
				turl{
					url:              "https://test.com/end/test/value/xx/1/2/3/4/5/6/7",
					expectedParams:   map[string]string{"name1": "value", "name2": "xx"},
					unexpectedParams: []string{"name3"},
					expectedArrays:   map[string][]string{"keys": []string{"1", "2", "3", "4", "5", "6", "7"}}}}},

		tunit{
			dispatch: "somepath/:name1/:name2/:keys.../",
			tus: []turl{
				turl{
					url:              "https://test.com/somepath/v1/v2/1/2/",
					expectedParams:   map[string]string{"name1": "v1", "name2": "v2"},
					unexpectedParams: []string{"name3"},
					expectedArrays:   map[string][]string{"keys": []string{"1", "2"}}}}},
		tunit{
			dispatch: "somepath2/folder1/folder2",
			tus: []turl{
				turl{
					url: "https://test.com/somepath2/folder1/folder2",
				}}},
	}

	d := Dispatcher{}

	for i, u := range units {
		du, err := url.Parse(u.dispatch)
		if err != nil {
			panic(err)
		}

		err = d.AddRoute(du, i)
		if err != nil {
			panic(err)
		}
	}

	for i, u := range units {
		fmt.Printf("tunit %v\n", i)

		for _, tu := range u.tus {

			u, err := url.Parse(tu.url)
			if err != nil {
				panic(err)
			}

			oa, err := d.Dispatch(u)
			if err != nil {
				panic(err)
			}

			if oa.amap.tag != i {
				t.Fatal("not expected tag.")
			}

			// test the expected params.
			for k, v := range tu.expectedParams {
				value, err := oa.ParamWithName(k)
				if err != nil {
					t.Fatal(err.Error())
				}

				if v != value {
					t.Fatalf("%v != %v\n", v, value)
				}
			}

			// test unexpected params.
			for _, n := range tu.unexpectedParams {
				_, err := oa.ParamWithName(n)
				if err == nil {
					t.Fatalf("param found for key which should not be there (%v)", n)
				}
			}

			// test expected arrays
			for k, a := range tu.expectedArrays {
				arr, err := oa.ArrayWithName(k)
				if err != nil {
					t.Fatal(err.Error())
				}

				if !reflect.DeepEqual(a, arr) {
					t.Fatal("wrong array")
				}
			}
		}
	}
}
