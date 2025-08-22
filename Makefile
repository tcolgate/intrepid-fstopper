.PHONY: flash build

flash:
	tinygo flash \
		-baudrate 9600 \
		-size full \
		-monitor \
		-scheduler none \
		-gc leaking \
		-target ./timer-1.2-target.json \
		-no-debug \
		-work \
		-stack-size 256B \
		-print-allocs .

build:
	tinygo build \
		-target ./timer-1.2-target.json \
		-stack-size 256B \
		-scheduler none \
		-gc leaking \
		-size full \
		-print-allocs . \
		-o intrep-ftimer.hex
	#	-o intrep-ftimer.elf
	#	-no-debug \
	du -h intrep-ftimer.hex

build-elf:
	tinygo build \
		-target ./timer-1.2-target.json \
		-stack-size 256B \
		-scheduler none \
		-gc leaking \
		-size full \
		-print-allocs . \
		-o intrep-ftimer.elf
	du -h intrep-ftimer.hex

monitor:
	tinygo monitor -baudrate 9600

test:
	go test ./num
