package main

import (
	"flag"
	"fmt"
	"io"
	"kek/defaults"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"

	"github.com/gomarkdown/markdown"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

func codeHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {

	if code, ok := node.(*ast.CodeBlock); ok {
		quick.Highlight(w, string(code.Literal), string(code.Info), "html", "monokailight")
		return ast.GoToNext, true
	}

	return ast.GoToNext, false
}

func mdToHtml(md []byte, absolutepath string) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Mmark
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.TOC

	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: codeHook,
		AbsolutePrefix: absolutepath,
	}

	renderer := html.NewRenderer(opts)

	res := markdown.Render(doc, renderer)

	return []byte(strings.Replace(defaults.Template, "{{}}", string(res), 1))
}

func htmlToPdf(input []byte) (io.ReadCloser, error) {
	page := rod.New().MustConnect().MustPage("about:blank")

	if err := page.SetDocumentContent(string(input)); err != nil {
		return nil, err
	}

	if err := page.WaitLoad(); err != nil {
		return nil, err
	}

	time.Sleep(2 * time.Second)

	return page.PDF(&proto.PagePrintToPDF{})

}

func main() {
	// Define a string flag named "input" with a default value and usage message
	input := flag.String("input", "", "Usage: specify an input file")
	output := flag.String("output", "", "Usage: specify an output file, default is the input.pdf")

	flag.Parse()

	if *input == "" {
		fmt.Println(errorStyle.Render("Error: must pass a .md file"))
		return
	}

	if *output == "" {
		i := *input
		ext := filepath.Ext(*input)
		if ext != "" {
			*output = i[:len(i)-len(ext)]
		}
		*output += ".pdf"
	}

	dir, err := filepath.Abs(filepath.Dir(*input))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	content, err := os.ReadFile(*input)
	if err != nil {
		panic(err)
	}

	page := mdToHtml(content, "file://"+dir+"/")

	// log
	log, _ := os.Create("data/page.html")
	log.WriteString(string(page))

	res, err := htmlToPdf(page)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(*output)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, res)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
