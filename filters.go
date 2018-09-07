package main

import (
	"bytes"
	//	"github.com/knieriem/markdown"
)

type IFilter interface {
	filter([]byte) bytes.Buffer
}

func doFilter(extName string) IFilter {
	//if extName == "md" || extName == "markdown" {
	//	return &mdFilter{}
	//}
	return nil
}

type mdFilter struct {
}

/*
func (*mdFilter) filter(input []byte) bytes.Buffer {
		p := markdown.NewParser(&markdown.Extensions{Smart: true})
		var buf bytes.Buffer
		buf.WriteString(`<html><header><style type="text/css">`)
		buf.WriteString(CSS)
		buf.WriteString("</style></header><body>")
		p.Markdown(bytes.NewReader(input), markdown.ToHTML(&buf))
		buf.WriteString("</body></html>")
		return buf

	return nil
}
*/
