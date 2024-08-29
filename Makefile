include scripts/make_env.sh

build:
	chmod +x scripts/build.sh
	scripts/build.sh $(TARGET) $(TAG)

register-githooks:
# if you want to ignore the hooks, use git commit --no-verify
	chmod u+x githooks/*
	git config core.hooksPath githooks # go >= 2.9

.PHONY: build register-githooks
