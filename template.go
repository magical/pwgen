package main

import "html/template"

type tmplContext struct {
	Words  []string
	Digits []int
}

var tmpl = template.Must(template.New("password").Parse(`<!doctype html>
<meta charset="utf-8">
<title>Password Generator</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
  header { text-align: center; }
  body { max-width: 500px; margin: auto; }
  .passwords, .digits { font: 20px sans-serif; text-align: center; }
  .passwords, .digits { max-width: 80%; margin: auto; position: relative; }

  .passwords span { display: inline-block; padding: .5em 0; }
  .passwords span.short { width: 26.4%; }
  .passwords span.large { width: 40%; }

  .digits { margin-top: .5em; }
  .digits span { display: inline-block; padding: .5em .5em; }

  footer { font-size: small; color: #808080; max-width: 80%; margin: 4em auto 0; }

  @media screen and (min-width: 500px) {
    body { max-width: 800; }
    .passwords span.short { width: 20%; }
    .passwords span.large { width: 26.4%; }
  }
</style>

<header>
  <h1>Password Generator</h1>
  <p><i>Need a password? Here, have a dozen.</i>
</header>

<main>
  <div class="passwords">
    {{- range $i, $_ := .Words }}
      {{- if lt $i 12 }}
        <span class="short">{{ . }}</span>
      {{- else }}
        <span class="large">{{ . }}</span>
      {{- end }}
    {{- end }}
  </div>
  <div class="digits">
    {{- range .Digits }}
      <span class="digit">{{ . }}</span>
    {{- end }}
  </div>
</main>

<footer>
  <p>
    Maintained by Andrew Ekstedt.
    Written in Go.
    Randomness harvested from <code>/dev/urandom</code>, via the <a href="https://golang.org/pkg/crypto/rand/">crypto/rand</a> package.
    Word lists by the <a href="https://www.eff.org/deeplinks/2016/07/new-wordlists-random-passphrases">EFF</a>.
  </p>
</footer>
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
