package api

import (
	"log"
	"net/http"

	"github.com/ponzu-cms/ponzu/system/item"

	"github.com/tidwall/gjson"
	"golang.org/x/net/http2"
)

func push(res http.ResponseWriter, req *http.Request, pt func() interface{}, data []byte) {
	// Push(target string, opts *PushOptions) error
	if pusher, ok := res.(http.Pusher); ok {
		if p, ok := pt().(item.Pushable); ok {
			// get fields to pull values from data
			fields := p.Push()

			// parse values from data to push
			values := gjson.GetManyBytes(data, fields...)

			// push all values from Pushable items' fields
			for i := range values {
				val := values[i]
				val.ForEach(func(k, v gjson.Result) bool {
					if v.String() == "null" {
						return true
					}

					err := pusher.Push(v.String(), nil)
					if err != nil && err != http2.ErrRecursivePush {
						log.Println("Error during Push of value:", v.String())
					}

					return true
				})
			}
		}
	}

}
