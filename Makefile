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
	 	-no-debug \
		-o intrepid-fstopper-$(shell git describe).hex
	du -h intrepid-fstopper-$(shell git describe).hex

build-elf:
	tinygo build \
		-target ./timer-1.2-target.json \
		-stack-size 256B \
		-scheduler none \
		-gc leaking \
		-size full \
		-print-allocs . \
		-o intrepid-fstopper-$(shell git describe).elf
	du -h intrepid-fstopper-$(shell git describe).elf

monitor:
	tinygo monitor -baudrate 9600

test:
	go test ./num
