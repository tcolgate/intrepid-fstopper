.PHONY: flash build

flash:
	tinygo flash \
		-x \
		-baudrate 9600 \
		-monitor -target ./timer-1.2-target.json \
		-scheduler=tasks \
		-no-debug

build:
	tinygo build \
		-x \
		-baudrate 9600 \
		-monitor -target ./timer-1.2-target.json \
		-scheduler=tasks \
		-no-debug \
		-o intrep-ftimer.hex
	du -h intrep-ftimer.hex

monitor:
	tinygo monitor -baudrate 9600
