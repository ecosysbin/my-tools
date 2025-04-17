package datasource

const kubeconfigScript = `
SERVICE_ACCOUNT_DIR="/var/run/secrets/kubernetes.io/serviceaccount"
KUBERNETES_SERVICE_SCHEME=$(case $KUBERNETES_SERVICE_PORT in 80|8080|8081) echo "http";; *) echo "https"; esac)
KUBERNETES_SERVER_URL="$KUBERNETES_SERVICE_SCHEME"://"$KUBERNETES_SERVICE_HOST":"$KUBERNETES_SERVICE_PORT"
KUBERNETES_CLUSTER_CA_FILE="$SERVICE_ACCOUNT_DIR"/ca.crt
KUBERNETES_NAMESPACE=$(cat "$SERVICE_ACCOUNT_DIR"/namespace)
KUBERNETES_USER_TOKEN=$(cat "$SERVICE_ACCOUNT_DIR"/token)
KUBERNETES_CONTEXT="inCluster"

mkdir -p "$HOME"/.kube
cat << EOF > "$HOME"/.kube/config
apiVersion: v1
kind: Config
preferences: {}
current-context: $KUBERNETES_CONTEXT
clusters:
- cluster:
    server: $KUBERNETES_SERVER_URL
    certificate-authority: $KUBERNETES_CLUSTER_CA_FILE
  name: inCluster
users:
- name: podServiceAccount
  user:
    token: $KUBERNETES_USER_TOKEN
contexts:
- context:
    cluster: inCluster
    user: podServiceAccount
    namespace: $KUBERNETES_NAMESPACE
  name: $KUBERNETES_CONTEXT
EOF
`

const _aExample = `
apiVersion: v1
kind: Config
preferences: {}
current-context: inCluster
clusters:
- cluster:
    server: https://192.168.194.129:443
    certificate-authority: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
  name: inCluster
users:
- name: podServiceAccount
  user:
    token: eyJhbGciOiJSUzI1NiIsImtpZCI6IjlNYnhuand4Rl9wZjZtcGZSV25qVkJCN1MzcVBhaFBCc1E0U3F6aV9rbVUifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiLCJrM3MiXSwiZXhwIjoxNzUxNTMyNjM1LCJpYXQiOjE3MTk5OTY2MzUsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0IiwicG9kIjp7Im5hbWUiOiJnbyIsInVpZCI6IjczOTg4YWZkLTAyYzAtNDhkZi04NWQ3LTlmNjVhM2M1NzM4NiJ9LCJzZXJ2aWNlYWNjb3VudCI6eyJuYW1lIjoiZGVmYXVsdCIsInVpZCI6IjExM2U2MTg4LTcyNGItNGYwMS05NDQwLTdhOGE2YzYwNTMyZiJ9LCJ3YXJuYWZ0ZXIiOjE3MjAwMDAyNDJ9LCJuYmYiOjE3MTk5OTY2MzUsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.Lv4k_NjMIQ_w7qFqEUy39aZFb8eWjyZwioz-8PiGV83padp_LJnnGB1Z7XmvWDY9qWjwwb34fLOZsn2q2hHVaWhUWDX41AzZIDMy3rx3lvrVx7meGhnLvD8M56Q7c69nKsGV36-VOjaNy54S1y25mziY3ae6rliWgOA_Cdxt5VHQ4m0ER4mXG9c9UYj0Fr6_bJmtypQSKFM7jUdxGzgBE1-ly-jqG4eKIhni-qYGvSF3a9pZ8MjR86LsxWeXURDX6h9bSNciXgc7k6BdRgtuP59Yhuwyn99109eBjz1LZPuPrXYSXqo9FwP0GhyoAhJj85GJAa-pcs4OtnQPJNhqDw
contexts:
- context:
    cluster: inCluster
    user: podServiceAccount
    namespace: default
  name: inCluster
`
