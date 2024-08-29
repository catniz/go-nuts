
build:
	@dockerfile=$$([ -f cmd/$(TARGET)/Dockerfile ] && echo cmd/$(TARGET)/Dockerfile || echo build/common.Dockerfile); \
	echo "Building $(TARGET) image with $$dockerfile"; \
	docker build -t $(TARGET):latest -f $$dockerfile . --build-arg TARGET=$(TARGET)

register-githooks:
# if you want to ignore the hooks, use git commit --no-verify
	chmod u+x githooks/*
	git config core.hooksPath githooks # go >= 2.9

.PHONY: build register-githooks
