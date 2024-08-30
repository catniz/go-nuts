apiVersion: v1
kind: Secret
metadata:
  name: "{{ .namespace }}-secret"
  namespace: "{{ .namespace }}"
  labels:
    app: "{{ .namespace }}-secret"
type: Opaque
data:
  MYSQL_PASSWORD: "{{ .infra.rdb_password | b64enc }}"
  REDIS_PASSWORD: "{{ .redis.password | b64enc }}"
