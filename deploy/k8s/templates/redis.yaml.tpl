apiVersion: v1
kind: ConfigMap
metadata:
  name: redis-config
  namespace: "{{ .namespace }}"
  labels:
    app: redis
data:
  redis-config: |
    maxmemory 1gb
    maxmemory-policy allkeys-lru

---

apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis
  namespace: "{{ .namespace }}"
  labels:
    app: redis
spec:
  serviceName: redis
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
        - name: redis
          image: redis:latest
          command:
            - redis-server
            - "/redis-master/redis.conf"
            - "--requirepass"
            - "$(REDIS_PASSWORD)"
          env:
            - name: MASTER
              value: "true"
            - name: REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: "{{ .namespace }}-secret"
                  key: REDIS_PASSWORD
          ports:
            - containerPort: 6379
              name: redis
          volumeMounts:
            - mountPath: /redis-master-data
              name: data
            - mountPath: /redis-master
              name: config
      volumes:
        - name: data
          emptyDir: { }
        - name: config
          configMap:
            name: redis-config
            items:
              - key: redis-config
                path: redis.conf
# mount data (dump.rdb) in pv if needed

---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: redis
  name: redis
  namespace: default
spec:
  selector:
    app: redis
  ports:
    - port: 6379
      protocol: TCP
      targetPort: 6379
      nodePort: 30030
  type: NodePort