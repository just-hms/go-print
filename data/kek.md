# Large Markdown File

## Introduction

This is a large markdown file containing various sections and elements to showcase the formatting capabilities of Markdown.

## Headers

### Header 3

#### Header 4

##### Header 5

###### Header 6

## Text Formatting

Markdown allows you to format text in various ways:

- **Bold text** can be created using double asterisks or double underscores: **Bold Text**.

- *Italic text* can be created using single asterisks or single underscores: *Italic Text*.

- ~~Strikethrough text~~ can be created using double tilde: ~~Strikethrough Text~~.

## Lists

Markdown supports both ordered and unordered lists.

### Unordered List

- Item 1
- Item 2
  - Subitem A
  - Subitem B
- Item 3

### Ordered List

1. First Item
2. Second Item
   1. Subitem 1
   2. Subitem 2
3. Third Item

## Links and Images

You can create links like this: [Visit OpenAI](https://www.openai.com).

You can embed images like this:
![Markdown Logo](kek.jpeg)

## Code

You can include inline code by wrapping it in backticks, like `print("Hello, World!")`.

For code blocks, use triple backticks:

```py
def greet(name):
    print(f"Hello, {name}!")
```

```html
<p><span class="math display">\[
  \xi
\]</span></p>
        </div>

        <!-- mathjax -->
        <script src="https://cdn.mathjax.org/mathjax/latest/MathJax.js?config=TeX-AMS-MML_HTMLorMML" type="text/javascript"></script>
    </body>
```


```go
func codeHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {

	if code, ok := node.(*ast.CodeBlock); ok {
		quick.Highlight(w, string(code.Literal), string(code.Info), "html", "monokailight")
		return ast.GoToNext, true
	}

	return ast.GoToNext, false
}
```

## Formulas


$$
\xi
$$

$$
S (\omega)=\frac{\alpha g^2}{\omega^5} \,
e ^{[-0.74\bigl\{\frac{\omega U_\omega 19.5}{g}\bigr\}^{-4}]}
$$