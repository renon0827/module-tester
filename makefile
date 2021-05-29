build:
	go vet
	go build

Test:
	make build
	module-tester.exe Test --scenario=scenario.json --out=result.json