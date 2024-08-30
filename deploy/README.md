# deploy

배포 세팅들
- docker-compose
- terraform
- kubernetes / helm

## k8s 배포 관련
go template engine (text/template) 을 사용하여 `yaml` 파일 생성

참고
- values/*.yaml 에는 정적인 값만 정의
- IntelliJ - yaml + go template syntax highlighting
  - Helm Values 플러그인 설치
  - *.yaml.tpl 파일 타입을 Helm template file 로 지정
- [scripts/gotpl.go](../scripts/gotpl.go)
  - 커스텀 함수는 funcMap 에 주석과 함께 추가
- [Makefile](../Makefile)
  - `make template-k8s-yaml`
  - [template_k8s_yaml.sh](../scripts/template_k8s_yaml.sh)
- [go text/template package](https://golang.org/pkg/text/template/)
- [sprig go template library](https://masterminds.github.io/sprig/)
