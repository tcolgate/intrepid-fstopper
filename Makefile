.PHONY: flash build

flash:
	tinygo flash \
		-x \
		-baudrate 9600 \
		-monitor \
		-target ./timer-1.2-target.json \
		-no-debug \
		-print-allocs main\

build:
	tinygo build \
		-x \
		-baudrate 9600 \
		-monitor \
		-target ./timer-1.2-target.json \
		-no-debug \
		-o intrep-ftimer.hex
	du -h intrep-ftimer.hex

monitor:
	tinygo monitor -baudrate 9600

test:
	go test ./num
