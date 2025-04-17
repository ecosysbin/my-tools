# MUST

## Pass the vcluster id

~~~bash
sh 1.3-combined-script.sh <vcluster-id>
~~~

## Set Environment Variables

~~~bash
export VCLUSTER_ID="<vcluster-id>"
export IMAGE_NAME="harbor.dev01.zetyun.cn/dev02-2-datacenter/vcluster/plugins/vcluster-plugin-storagemanager:develop-20241122180459"
export SIZE="<size>"
export VOLUMENAME="<volume-name>"
export TENANT_ID="<tenant-id>" 
export APIKEY="<apikey>"

sh 1.3-enable-storage-plugin.sh
~~~

# OPTIONAL

## Pass the vcluster id

~~~bash
sh 1.3-enable-sync-crd.sh <vcluster-id>
~~~