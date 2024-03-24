run:
	@templ generate
	@go run main.go

templ:
	@templ generate -watch -proxy=http://localhost:3000

sass:
	@sass --no-source-map --watch  view/stylesheets/tallytome.scss public/tallytome.css

build:
	@rm -rf tmp/release/tallytome
	@mkdir -p tmp/release/tallytome
	@templ generate
	@go build main.go
	@mv main tallytome
	@mv tallytome tmp/release/tallytome
	@cp -r public tmp/release/tallytome