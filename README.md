# go-wechaty: Word Counter 插件 （WIP）

统计限定小时内群成员发言词数。仅统计文字类消息。

该插件运行于 [go-wechaty](https://github.com/wechaty/go-wechaty)

注意：本项目使用了[插件机制](https://github.com/FINCTIVE/go-wechaty/tree/plugin)，插件分支代码目前并未合并到项目主分支中。

## 使用方式

```go
import "github.com/FINCTIVE/wordcounter"

func main() {
    var bot = wechaty.NewWechaty()
    bot.Use(wordcounter.New(wordcounter.Config{
            SearchKeyword:  "#排名",
            MaxResultCount: 10,
            Hours:          6,
        }))
    }
    // starting bot ...
}
```

## 配置

- `SearchKeyword`：当群聊消息内出现SearchKeyword，机器人发送排名消息。

- `MaxResultCount`：排名消息最大人数。

- `Hours`：统计范围。