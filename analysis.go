package wildyam

import (
	"fmt"
	"go/parser"
	"go/token"
)

func ConvertFile(fp string) {
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, fp, nil, parser.Mode(0))
	if err != nil {
		return
	}
	fmt.Println(f)
}
