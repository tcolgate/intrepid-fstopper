.PHONY: flash build

flash:
	tinygo flash \
		-baudrate 9600 \
		-monitor -target ./timer-1.2-target.json \
		-scheduler=tasks \
		-no-debug

build:
	tinygo build \
		-baudrate 9600 \
		-monitor -target ./timer-1.2-target.json \
		-scheduler=tasks \
		-no-debug .

