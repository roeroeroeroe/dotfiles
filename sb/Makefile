BIN := sb

.PHONY: build
build:
ifneq (,$(findstring with_pulse,$(TAGS)))
	go run ./cmd/gen_volume_pulse_constants
endif
	go build -tags "$(TAGS)" -ldflags="-s -w" -v -o $(BIN)

.PHONY: clean
clean:
	rm -vf $(BIN) ./components/volume_pulse.h
