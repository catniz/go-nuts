apiVersion: v1
kind: ConfigMap
metadata:
  name: "{{ .namespace }}-configmap"
  namespace: "{{ .namespace }}"
  labels:
    app: "{{ .namespace }}-configmap"
data:
  MYSQL_USER: "{{ .infra.rdb_username }}"
