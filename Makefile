SOURCE := $(shell git rev-parse --show-toplevel)

include $(SOURCE)/build/make/dev.mk
