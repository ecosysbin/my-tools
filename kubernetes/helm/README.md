# 添加helm仓库
helm repo add kubeovn https://kubeovn.github.io/kube-ovn/
helm repo add vcluster https://charts.loft.sh

# 查询仓库列表
helm repo list

# 安装一个软件包
helm install kube-ovn kubeovn/kube-ovn

# 安装软件包，带参数（使用了--namesapce，查询时也得带上相同的ns查询）
helm install kube-ovn kubeovn/kube-ovn -f ./values.yaml
helm install openebs openebs/openebs -f ./values.yaml --namespace kube-public
 
# 使用已下载的helm软件包进行安装
1. 下载软件包
helm pull ingress-nginx/ingress-nginx
2. 解压，并修改values.yaml文件
3. 安装
helm install ingress-nginx ingress-nginx -f ingress-nginx/values.yaml --namespace ingress

# 下载软件包
helm pull kubeovn/kube-ovn

# 查询已安装的helm软件包
helm list --namespace kube-public

# 升级版本
helm upgrade vc-kubeflow vcluster/vcluster -n vc-kubeflow -f ./package/values.yaml