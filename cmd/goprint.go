package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/just-hms/goprint/internal/defaults"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"golang.org/x/sync/errgroup"

	"github.com/alecthomas/chroma/quick"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"

	"github.com/gomarkdown/markdown"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

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
	// MAYBE: add this back html.TOC
	htmlFlags := html.CommonFlags | html.HrefTargetBlank

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

func mdToPdf(dir string, content []byte) (io.ReadCloser, error) {
	page := mdToHtml(content)
	return htmlToPdf(page, "file://"+dir+"/")
}

func main() {
	// Define a string flag named "input" with a default value and usage message
	input := flag.String("input", "", "Usage: specify an input file or a folder")
	output := flag.String("output", "", "Usage: specify an output file, default is the input.pdf")

	flag.Parse()

	// handle flags
	if *input == "" {
		fmt.Println("Error: must pass a .md file or a folder")
		return
	}

	// check if file exists
	fileInfo, err := os.Stat(*input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// TODO: make this better
	if *output == "" {
		*output = "data/kek.pdf"
	}

	var (
		results []string
		dir     string
	)

	// TODO: cretae a pdf merda here
	if fileInfo.IsDir() {
		dir, err = filepath.Abs(*input)
		if err != nil {
			panic(err)
		}

		files := []string{}

		err := filepath.Walk(dir, func(path string, info os.FileInfo, _ error) error {
			if info.IsDir() {
				return nil
			}
			files = append(files, path)
			return nil
		})
		if err != nil {
			panic(err)
		}

		wg := errgroup.Group{}
		mu := sync.Mutex{}

		for i, path := range files {
			path := path
			i := i
			wg.Go(func() error {
				content, err := os.ReadFile(path)
				if err != nil {
					panic(err)
				}
				res, err := mdToPdf(dir, content)
				if err != nil {
					return err
				}

				// put the res into a temp file and store its name
				tempFile, err := os.CreateTemp("", fmt.Sprintf("%d_temp_pdf_", i))
				if err != nil {
					return nil
				}
				_, err = io.Copy(tempFile, res)
				if err != nil {
					return err
				}
				mu.Lock()
				results = append(results, tempFile.Name())
				mu.Unlock()
				return nil
			})
		}
		if err := wg.Wait(); err != nil {
			panic(err)

		}

	} else {
		dir, err := filepath.Abs(filepath.Dir(*input))
		if err != nil {
			panic(err)
		}

		content, err := os.ReadFile(*input)
		if err != nil {
			panic(err)
		}

		res, err := mdToPdf(dir, content)
		if err != nil {
			panic(err)
		}

		// put the res into a temp file and store its name

		tempFile, err := os.CreateTemp("", "temp_pdf_")
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(tempFile, res)
		if err != nil {
			panic(err)
		}
		results = append(results, tempFile.Name())
	}

	config := model.NewDefaultConfiguration()

	// Merge the temporary PDF files.
	if err := api.MergeCreateFile(results, *output, config); err != nil {
		fmt.Printf("Error merging PDFs: %v\n", err)
		os.Exit(1)
	}

	for _, tempFile := range results {
		err := os.Remove(tempFile)
		if err != nil {
			fmt.Printf("Error deleting temporary file %s: %v\n", tempFile, err)
		}
	}
}
