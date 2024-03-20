run:
	@templ generate
	@go run main.go

templ:
	@templ generate -watch -proxy=http://localhost:3000

sass:
	@sass --no-source-map --watch  public/stylesheets/tallytome.scss public/stylesheets/tallytome.css