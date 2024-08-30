include scripts/make_env.sh

# ex) make build TARGET=server TAG=latest
build:
	@bash scripts/build.sh $(TARGET) $(TAG)

register-githooks:
# if you want to ignore the hooks, use git commit --no-verify
	chmod u+x githooks/*
	git config core.hooksPath githooks # go >= 2.9

# ex1) make template-k8s-yaml ENV=dev
# ex2) make template-k8s-yaml ENV=dev TARGET=redis.yaml
template-k8s-yaml:
	@bash scripts/template_k8s_yaml.sh $(ENV) $(TARGET)

.PHONY: build
