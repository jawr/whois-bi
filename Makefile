SOURCE := $(shell git rev-parse --show-toplevel)

include $(SOURCE)/build/make/dev.mk
include $(SOURCE)/build/make/nuke.mk
