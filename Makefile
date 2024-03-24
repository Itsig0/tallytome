run:
	@templ generate
	@go run main.go

templ:
	@templ generate -watch -proxy=http://localhost:3000

sass:
	@sass --no-source-map --watch  view/stylesheets/tallytome.scss public/tallytome.css