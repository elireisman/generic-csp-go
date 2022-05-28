.DEFAULT: run

.PHONY: run
run:
	@for f in `find . -name 'main.go'`; do echo; echo "[PROBLEM] $$f"; go run $$f; echo; done

