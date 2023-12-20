db-url := "sqlite:./results.db"

run:
  go run main.go -debug server

db-up:
  dbmate --url={{db-url}} up

db-down:
  dbmate --url={{db-url}} down

webpack:
  npx webpack

webpack-debug:
  npx webpack --mode=development

fmt:
  biome format . --write