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
	url              string
	expectedParams   map[string]string
	unexpectedParams []string
	expectedArrays   map[string][]string
}

type tunit struct {
	dispatch        string
	shouldFailOnAdd bool
	tus             []turl
}

func TestTest(t *testing.T) {

	units := []tunit{
		tunit{
			dispatch: "somepath2/folder2/:p1",
			tus: []turl{
				turl{
					url:            "https://test.com/somepath2/folder2/v1",
					expectedParams: map[string]string{"p1": "v1"},
				}}},
		tunit{
			dispatch: "start/test/:name1/:name2/xx/:k1/:key.../part3/:v1?qk&qk2",
			tus: []turl{
				turl{
					url:              "https://test.com/start/test/value/xx/vk1/1/2/3/4/5/part3/somevalue?qk=1111",
					expectedParams:   map[string]string{"name1": "value", "k1": "vk1", "v1": "somevalue", "qk": "1111"},
					unexpectedParams: []string{"name2", "qk2"},
					expectedArrays:   map[string][]string{"key": []string{"1", "2", "3", "4", "5"}}},
				turl{
					url:              "https://test.com/start/test/value/xx/vk1/part3/somevalue?qk=1111&qk2=2222",
					expectedParams:   map[string]string{"name1": "value", "k1": "vk1", "v1": "somevalue", "qk": "1111", "qk2": "2222"},
					unexpectedParams: []string{"name2"},
					expectedArrays:   map[string][]string{"key": []string{}}},
				turl{
					url:              "https://test.com/start/test/value/xx/vk1/tt/part3?qk=1111",
					expectedParams:   map[string]string{"name1": "value", "k1": "vk1", "qk": "1111"},
					unexpectedParams: []string{"name2", "v1"},
					expectedArrays:   map[string][]string{"key": []string{"tt"}}},
			},
		},
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
		tunit{
			dispatch:        "somepath2/folder1",
			shouldFailOnAdd: true},
		tunit{
			dispatch:        "somepath2/folder1/folder2/folder3",
			shouldFailOnAdd: true},
		tunit{
			dispatch:        "somepath2/folder1/folder2/:arr...",
			shouldFailOnAdd: true},
		tunit{
			dispatch: "dispatch/:ar1.../folder1/folder2",
			tus: []turl{
				turl{
					url:            "https://test.com/dispatch/1/2/3/4/folder1/folder2",
					expectedArrays: map[string][]string{"ar1": []string{"1", "2", "3", "4"}},
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
			if !u.shouldFailOnAdd {
				panic(err)
			}
			continue
		}
	}

	for i, u := range units {
		fmt.Printf("tunit %v\n", i)

		for _, tu := range u.tus {

			fmt.Printf("dispatch %v\n", tu.url)

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

				if !compareStringSlices(a, arr) {
					t.Fatal(fmt.Sprintf("wrong array %v,%v", arr, a))
				}
			}
		}
	}
}

func compareStringSlices(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}

	for i, l := range left {
		if l != right[i] {
			return false
		}
	}

	return true
}
