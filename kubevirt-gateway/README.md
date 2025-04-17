<!-- 
Copyright 2023 The Zetyun.GCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. 
-->

# kubevirt-gateway Gateway

> Makefile使用说明

```shell
$ make help # 查看帮助
```

> 本地测试

1. 使用`make run`命令直接运行在本地8083端口, 可以使用`CONFIGPATH`变量指定配置文件，默认路径为`./config.yaml`

> 生成或修改swagger配置文件

1. 先执行`make install-swag`命令安装`swag`
2. 按照文档写好API注释，可参考当前代码或者[swag官方文档](https://github.com/swaggo/swag/blob/master/README_zh-CN.md#%E5%BF%AB%E9%80%9F%E5%BC%80%E5%A7%8B)
3. 写好注释后，执行`make swag`命令，可以生成或更新文档，然后执行`make run`命令就可以本地查看文档了
