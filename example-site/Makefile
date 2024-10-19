#!/usr/bin/make -f

dev:
	echo "Navigate a browser to http://localhost:7070/" && \
		hugoinho-dev -content "./content" -templates "./templates" -target "./rendered" -with-drafts -with-future

generate:
	hugoinho -content "./content" -templates "./templates" -target "./rendered"

clean:
	rm -rf "./rendered" && mkdir "./rendered"
