# Go 日志框架使用详解

## 1. log包使用

Go语言内置的log包实现了简单的日志服务。打印日志的方式有三种：调用函数Print系列(Print|Printf|Println）、调用函数Fatal系列（Fatal|Fatalf|Fatalln）、和调用函数Panic系列（Panic|Panicf|Panicln），内置log包没有日志分级的概念。

### 1.1 入门使用

log包中预定义了一个Logger类型的对象，可以直接调用log包中的以下方法进行日志输出。

#### 1.1.1 调用函数Print系列

```go
package main

import (
    "log"
)

func main() {
    log.Print("这是一条很普通的日志。")
    log.Println("这是一条很普通的日志。")
    v := "格式化的"
    log.Printf("这是一条%s日志。\n", v)
}
```

输出：

2023/01/31 14:41:04 这是一条很普通的日志。
2023/01/31 14:41:04 这是一条很普通的日志。
2023/01/31 14:41:04 这是一条格式化的日志。

#### 1.1.2 调用函数Fatal系列

```go
package main

import (
    "log"
)

func main() {
    log.Fatal("这是一条Fatal日志。")
    log.Fatalln("这是一条Fatal日志。")
    v := "格式化的"
    log.Fatalf("这是一条%sFatal日志。", v)
}
```

输出：

2023/01/31 14:45:24 这是一条Fatal日志。
exit status 1

<u>注：Fatal系列打印日志后，调用os.Exit(1)退出，后续程序不再执行。</u>

#### 1.1.3 调用函数Panic系列

```go
package main

import (
    "log"
)
func main() {
    v := "格式化的"
    log.Panic("这是一条会触发panic的日志。")
    log.Panicln("这是一条触发panic的日志。")
    log.Panicf("这是一条触发panic的%s日志。",v)
}
```

输出：

2023/01/31 14:50:35 这是一条会触发panic的日志。
panic: 这是一条会触发panic的日志。

goroutine 1 [running]:
log.Panic({0xc000109f60?, 0x60?, 0x0?})
  D:/lch/go/src/log/log.go:388 +0x65
main.main()
  d:/lch/goCode/src/go_code/chapter07/logtest/main.go:11 +0x48
exit status 2

<u>注：panic系列打印日志后，会抛出一个panic错误，程序异常退出。</u>

### 1.2 配置Logger

#### 1.2.1 配置flag

​    上一节演示了log打印日志的简单示例，默认情况下的logger只会打印日志的日期、时间、message信息，但是在一些情况下，我们需要知道记录该日志的文件名和行号，方便开发人员进行定位。为满足此要求，log提供了**SetFlags（flag int）**方法。Falg可选择的值如下：

![image-20230131151840105](C:\Users\liuchenhong.HOLLYSYS\AppData\Roaming\Typora\typora-user-images\image-20230131151840105.png)

Logger类型中预定了以上可选择的常量。

```go
func main() {
    log.SetFlags(log.Lmicroseconds | log.Ldate |log.Llongfile )
    log.Println("这是一条很普通的日志")
}

```

输出：

2023/01/31 15:22:45.971754 d:/lch/goCode/src/go_code/chapter07/logtest/main.go:15: 这是一条很普通的日志

<u>注：可以通过调用**Flags()**方法查看设置的flag值。</u>

#### 1.2.2 配置前缀

通过调用**SetPrefix(prefix string)**方法设置日志输出的前缀信息。

```go
package main

import (
    "log"
)

func main() {
    log.SetPrefix("[prefix] ")
    log.Println("这是一条很普通的日志")
}
```

输出：

[prefix] 2023/01/31 15:31:49 这是一条很普通的日志

<u>注:可以通过调用Prefix()方法查看设置的前缀信息。</u>

#### 1.2.3 配置日志输出位置

​     通过调用**SetOutput(w io.Writer)**函数设置标准logger输出目的地，默认是io.Stderr输出。

```go
import (
    "log"
    "os"
    "fmt"
)

func main() {
    filePath := "C:/Users/EDY/Desktop/log.txt"
    file,err := os.OpenFile(filePath, os.O_WRONLY | os.O_CREATE,0666)//第三个参数在window系统无用，第二个参数是创建一个新文件
    if err != nil {
        fmt.Println("打开文件出错，err=",err)
        return
    }
    //及时关闭file句柄
    defer file.Close()
    log.SetPrefix("[prefix] ")
    log.SetOutput(file)
    log.Println("这是一条很普通的日志")
}
```

输出：

<img src="C:\Users\liuchenhong.HOLLYSYS\AppData\Roaming\Typora\typora-user-images\image-20230131154156455.png" alt="image-20230131154156455" style="zoom: 80%;" />

### 1.3 自定义logger

log包提供了一个自定义logger的方法：**New(out io.Writer, prefix string, flag int)**。生成logger时，指定输出位置，前缀，flag属性信息。

```go
import (
    "log"
    "os"
    "fmt"
)

func main() {
    filePath := "C:/Users/EDY/Desktop/log.txt"
    file,err := os.OpenFile(filePath, os.O_WRONLY | os.O_CREATE,0666)//第三个参数在window系统无用，第二个参数是创建一个新文件
    if err != nil {
        fmt.Println("打开文件出错，err=",err)
        return
    }
    //及时关闭file句柄
    defer file.Close()
    logger := log.New(file,"[prefix] ",log.Lmicroseconds | log.Ldate |log.Llongfile)
    logger.Println("这是一条很普通的日志!")
}
```

输出：

![image-20230131155052053](C:\Users\liuchenhong.HOLLYSYS\AppData\Roaming\Typora\typora-user-images\image-20230131155052053.png)

### 1.4 总结

1. 内置log操作简单。
2. 内置log输出日志无细粒度的级别划分。

## 2. zerolog包使用

### 2.1 包地址

github仓库：https://github.com/rs/zerolog

### 2.2 入门使用

```go
import (
    "github.com/rs/zerolog/log"
)

func main() {
    log.Print("这是一条日志")
    }
```

输出：

{"level":"debug","time":"2023-01-31T16:49:20+08:00","message":"这是一条日志"}

<u>注：Print()方法打印日志的默认级别为debug,默认的输出位置为os.Stderr</u>

### 2.3 日志级别

![image-20230131164350281](C:\Users\liuchenhong.HOLLYSYS\AppData\Roaming\Typora\typora-user-images\image-20230131164350281.png)

zerolog使用链式api调用，示例如下：

```go
package main

import (
    "github.com/rs/zerolog/log"
)

func main() {
    log.Info().Msg("这是一条Info级别的日志")
    log.Warn().Msg("这是一条Warn级别的日志")
    log.Error().Msg("这是一条Warn级别的日志")
    v := "Panic"
    log.Panic().Msgf("这是一条%s级别的日志",v)
}

输出：
{"level":"info","time":"2023-01-31T16:58:52+08:00","message":"这是一条Info级别的日志"}
{"level":"warn","time":"2023-01-31T16:58:52+08:00","message":"这是一条Warn级别的日志"}
{"level":"error","time":"2023-01-31T16:58:52+08:00","message":"这是一条Warn级别的日志"}
{"level":"panic","time":"2023-01-31T16:58:52+08:00","message":"这是一条Panic级别的日志"}
panic: 这是一条Panic级别的日志

goroutine 1 [running]:
github.com/rs/zerolog.(*Logger).Panic.func1({0xc0000170e0?, 0x0?})
    D:/lch/goCode/src/github.com/rs/zerolog/log.go:376 +0x2d
github.com/rs/zerolog.(*Event).msg(0xc00005c0c0, {0xc0000170e0, 0x20})
    D:/lch/goCode/src/github.com/rs/zerolog/event.go:156 +0x2a5
github.com/rs/zerolog.(*Event).Msgf(0xc00005c0c0, {0x1011be7?, 0x1f?}, {0xc000089f60?, 0x0?, 0x0?})
    D:/lch/goCode/src/github.com/rs/zerolog/event.go:129 +0x4e
main.main()
    d:/lch/goCode/src/go_code/chapter07/zerotest1/main.go:15 +0x105
exit status 2
```

注：日志的输出，必须调用Msg,Msgf或Send方法。

### 2.4 添加字段

在日志打印信息中增加自定义字段和字段值。zerolog添加字段为强类型，可添加的基本数据类型如下：

- `Str`
- `Bool`
- `Int`, `Int8`, `Int16`, `Int32`, `Int64`
- `Uint`, `Uint8`, `Uint16`, `Uint32`, `Uint64`
- `Float32`, `Float64`

可链式添加多个字段，示例如下：

```go
import (
    "github.com/rs/zerolog/log"
)

func main() {
    log.Debug().
        Str("姓名", "lch").
        Float64("score", 833.09).
        Msg("日志信息")
    
    log.Debug().
        Str("Name", "Tom").
        Send()
}

输出：  

{"level":"debug","姓名":"lch","score":833.09,"time":"2023-01-31T17:17:25+08:00","message":"日志信息"}
{"level":"debug","Name":"Tom","time":"2023-01-31T17:17:25+08:00"}
```

### 2.5 全局信息设置

- `log.Logger`: 设置全局的logger。
- `zerolog.SetGlobalLevel`:设置可输出的日志的最低级别。
- `zerolog.DisableSampling`: If argument is `true`, all sampled loggers will stop sampling and issue 100% of their log events.
- `zerolog.TimestampFieldName`: 时间戳字段名。
- `zerolog.LevelFieldName`:日志级别字段名.
- `zerolog.MessageFieldName`: 日志信息字段名.
- `zerolog.ErrorFieldName`: 默认为error,error日志使用Err()方法时的信息字段名.
- `zerolog.TimeFieldFormat`: 时间格式，可选值：  `zerolog.TimeFormatUnix`, `zerolog.TimeFormatUnixMs` ， `zerolog.TimeFormatUnixMicro`.

示例如下：

```go
package main

import (
    //"os"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    //"fmt"
)

func main() {

    //时间戳字段名，默认为time
    zerolog.TimestampFieldName = "t"
    //日志级别字段名,默认为level
    zerolog.LevelFieldName = "l"
    //定义的信息的字段名，默认为message
    zerolog.MessageFieldName = "m"
    //时间戳格式
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    //设置打印日志的最低级别
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
    log.Info().
        Str("姓名", "lch").
        Float64("score", 833.09).
        Msg("日志信息")
    
    log.Debug().
        Str("Name", "Tom").
        Send()
}

输出：

{"l":"info","姓名":"lch","score":833.09,"t":1675157631,"m":"日志信息"}
```

### 2.6  自定义logger

以上几节使用的均为包中预设的Logger实例,除此之外，我们可以使用New(out io.Writer)生成一个Logger类型的对象。

```go
package main

import (
    "os"
    "github.com/rs/zerolog"
)

func main() {
    log := zerolog.New(os.Stderr)
    log.Info().
        Str("姓名", "lch").
        Float64("score", 833.09).
        Msg("日志信息")

}

输出：

{"level":"info","姓名":"lch","score":833.09,"message":"日志信息"}
```

### 2.7 添加文件和行数信息

​    一般要求打印日志信息中包含打印日志所在的文件和行数，便于开发人员定位问题。调用**With()**使用上下文信息，调用**Caller()**方法打印文件名称和行数。

```go
package main

import (
    "os"
    "github.com/rs/zerolog"
)

func main() {
    log := zerolog.New(os.Stderr)
    log = log.With().Caller().Logger()
    log.Info().
        Msg("日志信息")
}

输出：

{"level":"info","caller":"d:/lch/goCode/src/go_code/chapter07/zerotest1/main.go:14","message":"日志信息"}
```

### 2.8 设置多输出

使用**MultiLevelWriter(w ...os.Writer)**方法设置多个输出。

```go
package main

import (
  "github.com/rs/zerolog"
  "os"
  "fmt"
)
func main(){
    //设置输出文件路径
    filePath := "C:/Users/EDY/Desktop/log.txt"
    file,err := os.OpenFile(filePath, os.O_WRONLY | os.O_CREATE,0666)//第三个参数在window系统无用，第二个参数是创建一个新文件
    if err != nil {
        fmt.Println("打开文件出错，err=",err)
        return
    }
    //及时关闭file句柄
    defer file.Close()
    //输出到文件和控制台，输出可有多个
    mutil := zerolog.MultiLevelWriter(os.Stderr,file)
    log := zerolog.New(mutil)
    log = log.With().Caller().Logger()
    log.Info().Timestamp().Msg("hello world")
}

输出：

{"level":"info","time":"2023-02-01T09:21:28+08:00","caller":"d:/lch/goCode/src/go_code/chapter07/zerologtest/main.go:21","message":"hello world"}
```

![image-20230201092458327](C:\Users\liuchenhong.HOLLYSYS\AppData\Roaming\Typora\typora-user-images\image-20230201092458327.png)

