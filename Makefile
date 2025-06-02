.PHONY: flash

flash:
	tinygo flash -print-allocs -print-stacks -baudrate 9600 -monitor -target ./timer-1.2-target.json

