.PHONY: flash build

flash:
	tinygo flash \
		-baudrate 9600 \
		-size full \
		-monitor \
		-scheduler tasks \
		-gc leaking \
		-target ./timer-1.2-target.json \
		-no-debug \
		-work \
		-stack-size 256B \
		-print-allocs .
	

build:
	tinygo build \
		-baudrate 9600 \
		-scheduler tasks \
		-gc leaking \
		-size full \
		-monitor \
		-target ./timer-1.2-target.json \
		-no-debug \
		-print-allocs . \
		-o intrep-ftimer.hex
	du -h intrep-ftimer.hex

monitor:
	tinygo monitor -baudrate 9600

test:
	go test ./num
