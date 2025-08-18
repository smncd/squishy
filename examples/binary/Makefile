VERSION := 0.4.0
ENDPOINT := "https://gitlab.com/smncd/squishy/-/releases/v$(VERSION)/downloads/squishy-linux-amd64-$(VERSION)"
EXE := ./bin/squishy

download:
	mkdir -p ./bin
	curl -L $(ENDPOINT) -o $(EXE)
	chmod +x $(EXE)

start:
	$(EXE) --config-file ./squishy.yaml
