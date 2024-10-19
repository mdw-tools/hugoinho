# hugoinho

Yet another static site generator--a [tiny](https://en.wiktionary.org/wiki/-inho) [hugo](https://gohugo.io).


## Disclaimer:

1. I wrote this to generate static html for my own website. I sometimes modify its behavior to suit my purposes. There is no intention to support general-purpose use. See the license for additional disclaimers.
2. Topics/tags are only rendered once 3 separate articles reference them.
3. Rather than use this repo outright, why not create your own fork, or create your own static site generator from scratch? It's really not that difficult, and it's a fun, relatively small-sized project.

## Installation

1. Clone this repo
2. Change directory to newly cloned repo
3. Run `make install`


## Example Site

1. Install per above instructions
2. Change directory to `./example-site`
3. Run `make dev`
   - Will start a web server--open browser to http://localhost:7070 to load the site.
   - Refresh the browser after saving changes to content.
   - Restart the web server process after saving changes to templates, then refresh the browser.
4. Or, run `make generate` to generate the static html and exit.
   - From there you can copy/upload the html wherever you want.
