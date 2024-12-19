package utils

import (
	"context"
	"net/http"
	"fmt"
	"strings"
	"github.com/xhd2015/xgo/runtime/core"
  "github.com/xhd2015/xgo/runtime/trap"
)

func AutoLog(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trap.AddInterceptor(&trap.Interceptor{
			Pre: func(ctx context.Context, f *core.FuncInfo, args core.Object, results core.Object) (interface{}, error) {
				index := strings.Index(f.File, "sphinx-tribes")
				trimmed := f.File
				if index != -1 {
					trimmed = f.File[index:]
				}
				fmt.Printf("%s:%d %s\n", trimmed, f.Line, f.Name)

				return nil, nil
			},
		})
		fn(w, r)
	}
}
