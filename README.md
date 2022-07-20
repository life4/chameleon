# Chameleon

Chameleon is web application (blog engine) that reflects content from markdown files from a git repository. Powers [articles.orsinium.dev](https://articles.orsinium.dev/).

Features:

+ Markdown (full support for [CommonMark](https://commonmark.org/) and [GitHub Flavored Markdown](https://github.github.com/gfm/))
+ [Jupyter notebooks](https://jupyter.org/)
+ Minimalistic UI
+ Easy to use, no CI or a special repo structure required
+ Zero configuration
+ Single binary
+ Automatically pull the repo by schedule
+ Built-in prose linter ([Vale](https://github.com/errata-ai/vale))
+ Syntax highlighting ([Prism](https://prismjs.com/))
+ Formulas ([MathJax](https://www.mathjax.org/))
+ Emoji ([enescakir/emoji](https://github.com/enescakir/emoji))
+ Views count
+ Great performance and server-side caching
+ Optional password protection
+ Search
+ Minification ([minify](https://github.com/tdewolff/minify#examples))

## Usage

Build:

```bash
git clone https://github.com/life4/chameleon.git
cd chameleon
go build -o chameleon.bin .
```

Run:

```bash
./chameleon.bin --path ./path/to/repo/
```
