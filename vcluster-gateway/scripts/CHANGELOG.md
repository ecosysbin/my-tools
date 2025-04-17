## 1.1 to 1.21

- NetworkPolicy Egress 暴露 6060 8080 7777
    - shell: [1.21-add-port.sh](1.2.1/must/1.21-add-port.sh)
      commit: [feat（charts）: networkPolicy add 6060 and 8080 ports](https://gitlab.datacanvas.com/AlayaNeW/OSM/vcluster-gateway/commit/9b11b59eb07142fed1798df9be0715118d201157)

## 1.21 to 1.3

- 去掉 ResourceQuota 中的 `requests.storage` 字段，不再限制 Vcluster 总的存储容量
    - shell: [1.3-update-resourcequota.sh](1.3/must/.discarded/1.3-update-resourcequota.sh)
      commit: [feat(chart-dev): delete requests.storage in resourceQuotas](https://gitlab.datacanvas.com/AlayaNeW/OSM/vcluster-gateway/commit/8cd835316de7c6726d860fd524855ccd62e71980)

- 升级 Hooks Plugin 插件到 v0.8.1
    - shell: [1.3-upgrade-hooks-plugin.sh](1.3/must/.discarded/1.3-upgrade-hooks-plugin.sh)
      commit: [feat(chart-dev): Upgrade the hook plugin version to v0.8.1](https://gitlab.datacanvas.com/AlayaNeW/OSM/vcluster-gateway/commit/003debb22437f73b769c46c3c1968c61fc47cea2)
        - 同步 CoreDNS Pod 到宿主机时添加 `dc.com/tenant.source: system` label
          commit: [feat: coredns pod add labels: dc.com/tenant.source: system](https://gitlab.datacanvas.com/AlayaNeW/OSM/vcluster-gateway/commit/9025cbe79cbd4b1500e8e82eb726d56a0604178c)
        - 注入环境变量，支持用户构建镜像
          commit: [feat(plugin): Update Pod Hooks Plugin to support build user image function](https://gitlab.datacanvas.com/AlayaNeW/OSM/vcluster-gateway/commit/96577db699a63563baa3a09b9ed28e45b463474c)