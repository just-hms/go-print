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

func mdToHtml(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Mmark
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.TOC

	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: codeHook,
	}

	renderer := html.NewRenderer(opts)

	renderedHTML := markdown.Render(doc, renderer)

	// replace template placeholder
	return []byte(
		strings.Replace(defaults.Template, "{{}}", string(renderedHTML), 1),
	)
}

func htmlToPdf(input []byte, url string) (io.ReadCloser, error) {
	if url == "" {
		url = "about:blank"
	}

	page := rod.New().MustConnect().MustPage(url)

	if err := page.SetDocumentContent(string(input)); err != nil {
		return nil, err
	}

	if err := page.WaitLoad(); err != nil {
		return nil, err
	}

	time.Sleep(2 * time.Second)

	return page.PDF(&proto.PagePrintToPDF{})

}
func concatFiles(dir string) ([]byte, error) {
	var result []byte

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {

			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			result = append(result, content...)
			result = append(result, '\n')
			result = append(result, []byte("---")...)
			result = append(result, '\n')
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func main() {
	// Define a string flag named "input" with a default value and usage message
	input := flag.String("input", "", "Usage: specify an input file or a folder")
	output := flag.String("output", "", "Usage: specify an output file, default is the input.pdf")

	flag.Parse()

	// handle flags
	if *input == "" {
		fmt.Println(errorStyle.Render("Error: must pass a .md file"))
		return
	}

	// check if file exists
	fileInfo, err := os.Stat(*input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if *output == "" {
		if fileInfo.IsDir() {
			*output = *input + ".pdf"
		} else {
			i := *input
			ext := filepath.Ext(*input)
			if ext != "" {
				*output = i[:len(i)-len(ext)]
			}
			*output += ".pdf"
		}
	}

	var (
		content []byte
		dir     string
	)

	if fileInfo.IsDir() {
		dir, err = filepath.Abs(*input)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		content, err = concatFiles(dir)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	} else {
		dir, err = filepath.Abs(filepath.Dir(*input))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		content, err = os.ReadFile(*input)
		if err != nil {
			panic(err)
		}
	}

	page := mdToHtml(content)

	// log
	log, _ := os.Create("data/page.html")
	log.WriteString(string(page))

	res, err := htmlToPdf(page, "file://"+dir+"/")
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
