db-url := "sqlite:./results.db"

run:
  just db-up
  go run main.go -debug server

clean:
  just db-down

pack:
  npx webpack --mode=development

db-up:
  dbmate --url={{db-url}} up

db-down:
  dbmate --url={{db-url}} down

build:
  npx webpack


fmt:
  biome format . --write

test:
  npx vitest run

  go test ./...