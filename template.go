package main

import "html/template"

type tmplContext struct {
	Words []string
}

var tmpl = template.Must(template.New("password").Parse(tmplSrc))
var tmplSrc = `<!doctype html>
<meta charset="utf-8">
<title>Password Generator</title>
<style>
  header { text-align: center; }
  body { max-width: 500px; margin: auto; }
  .passwords { font: 20px sans-serif; }
  .passwords span { display: block-inline; padding: .5em; }
</style>

<header>
  <h1>Password Generator</h1>
  <p><i>Need a password? Here, have a dozen.</i>
</header>

<div class="passwords">
  {{ range .Words }}
     <span>{{ . }}</span>
  {{ end }}
</div>
`
