.PHONY: run

ALL: run

run:
	go mod tidy
	go run .

clean:
	rm data.db