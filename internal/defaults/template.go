package defaults

var Template = `
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" lang="$lang$" xml:lang="$lang$"$if(dir)$ dir="$dir$"$endif$>
    <head>
        <meta charset="utf-8" />
        <meta name="generator" content="pandoc" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=yes" />

        <!-- tailwind css styles -->
        <script src="https://cdn.tailwindcss.com"></script>
        <script src="https://cdn.tailwindcss.com?plugins=forms,typography,aspect-ratio,line-clamp"></script>
        <style>
            pre {
                white-space: pre !important;
                overflow-x: hidden !important;
            }
        </style>
    </head>
    <body>
        <div class="prose prose-sm prose-img:rounded-lg prose-img:max-h-96 prose-img:max-w-sm !container w-full px-4 mx-auto">
            {{}}
        </div>

        <!-- mathjax -->
        <script src="https://cdn.mathjax.org/mathjax/latest/MathJax.js?config=TeX-AMS-MML_HTMLorMML" type="text/javascript"></script>
    </body>
</html>
`
