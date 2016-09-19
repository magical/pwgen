package main

import "html/template"

type tmplContext struct {
	Words []string
}

var tmpl = template.Must(template.New("password").Parse(`<!doctype html>
<meta charset="utf-8">
<title>Password Generator</title>
<meta name="viewport" value="width=device-width, initial-scale=1">
<style>
  header { text-align: center; }
  body { max-width: 500px; margin: auto; }
  .passwords { font: 20px sans-serif; text-align: center; }
  .passwords { max-width: 80%; margin: auto; }

  .passwords span { display: inline-block; min-width: fit-content; width: 26.4%; padding: .5em 0; }

  @media screen and (min-width: 500px) {
    body { max-width: 800; }
    .passwords span { width: 20%; }
  }
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
`))

var _ = template.Must(tmpl.New("list").Parse(`<!doctype html>
<meta charset="utf-8">
<title>Password Generator - word list</title>
<style>

</style>

<header>
  <h1>Word List</h1>
</header>

<ol>
  {{ range .Words }}
    <li>{{ . }}</li>
  {{ end }}
</ol>
`))
