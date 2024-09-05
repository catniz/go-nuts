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

# ex) make ent-new NAME=SystemConfig SCHEMA=bullish
# ex) make ent-new NAME=User,SystemConfig SCHEMA=bullish
ent-new:
	go run ./pkg/ent/entc/new/main.go \
		-target ./pkg/ent/schema \
		-template ./pkg/ent/entc/new/entnew.tmpl \
		-tables $(NAME) \
		-schema $(SCHEMA)

# ex) make ent-generate
ent-generate:
	go run ./pkg/ent/entc/gen/main.go \
		-target ./pkg/ent/schema \
		-templatedir ./pkg/ent/entc/gen/ \
		-feature sql/upsert,sql/execquery,sql/modifier

.PHONY: build
