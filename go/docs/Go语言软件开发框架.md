# 开发语言选型：
Go

# 项目代码结构(主要参考架构： https://github.com/evrone/go-clean-template):
cmd/	
    main.go              //  --help  帮助提示  --version 版本， 加载配置，项目启动。使用cabra框架进行管理
config/                  //  配置
    lang.go              //  多语言国际化
    i18n
        zh.toml
        en.toml
    config.go
    config.yml
internal
    app/                 // 应用层（项目启动，初始化，暴露API服务）
        app.go
        migrate.go       // 数据库表迁移、升级
    controller/          // 服务的handler层，默认只提供REST HTTP ( GIN 框架 ), 需要其他handler，通过目录进行区分
        v1
            route.go
    entity/              // 实体
        translation.go   
    usecase/             // 用户接口
        repo/            // 存储
            translation_sqllite.go
        webapi/          // web api服务，通过web api访问访问其他服务
            translation_google.go
        interfaces.go    // 接口定义
        translation.go   // 业务逻辑
    constants/           // 常量
    utils/               // 工具包
migrations               // 数据库升级迁移脚本
pkg/                     // 外部依赖框架
    httpserver/          
        server.go        
    sqllite/
        sqllite.go       
    logger/
        logger.go        
tests/
vendor/                   // 外部项目软件依赖管理
Makefile                  // 编译脚本
go.mod                     
go.sum
README.md                 // 引用了哪些外部依赖（直接引用）
CHANGELOG.md              // 版本变更记录

# 开源软件引入：
cobra
gin
sqlite
gorm
zerolog
gi18n
cleanenv
gin-swagger

# 数据库选型：sqlite3

# Api设计:
遵循RestAPI规范

场景1.  没有前端访问接口

​    /{子系统名称}/v1/.../xxx                             // 使用小写字母，中划线-分词，
​    例如：GET /hsm-env/v1/apps/run          // 运行节点所有容器应用 

场景2.  有前端访问接口

​    /api/{子系统名称}/v1/.../xxx                      // 使用小写字母，中划线-分词，
​    例如：GET /api/hsm-install/v1/run         // 部署应用 

特别说明：

 1. 对于有前端访问的接口，后端操作失败后，返回给前端的message中的内容都是要可展示的。像 “参数校验错误”，“服务访问异常”。后端对于一些不可控的异常，如数据库访问，需要后端接口对异常进行处理（如打印日志）后 ，返回错误信息（如：message: “服务访问异常”）到前端. 

    

请求Body体/返回体数据格式: json               // 默认格式，需要其他格式，单独讨论

请求返回码：
    返回码	Code	     msg	         含义
    200	                     Success            请求成功
    404	                     Not Found	    无效的URL

返回数据格式：

   成功：{code:0,message:”success”,result:{data}}

   失败:   {code:1,message:”err msg”,result:""}

```
func (r *translationRoutes) history(c *gin.Context) {
	translations, err := r.t.History()
	if err != nil {
		r.l.Error(err, "http - v1 - history")
		// errorResponse(c, http.StatusInternalServerError, "database problems")
		// 报错
		c.JSON(http.StatusOK, Error(err.Error()))
		return
	}
	// 成功
	c.JSON(http.StatusOK, Ok(translations))
}

func Ok(data interface{}) Resp {
	return Resp{Code: 0, Message: "success", Result: data}
}

func Error(err string) Resp {
	return Resp{Code: 1, Message: err, Result: nil}
}

type Resp struct {
	Code    int         `json:"code" `
	Message string      `json:"message" `
	Result  interface{} `json:"result" `
}
```




# 打包编译：
Makefile


# 开发自测试：
单元测试： 主要业务代码要有单元测试用例
API测试:   API接口测试从手动测试逐渐过渡到自动化测试


# 代码开发流程：
代码走查：版本合并。各子系统根据计划自己控制节奏。
代码review会议：发现问题讨论。

# go语言开发技能分享
开发全员参与不定期总结分享