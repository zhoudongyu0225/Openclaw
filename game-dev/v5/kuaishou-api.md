# 快手弹幕玩法 API 接口总览

## 一、基础能力（上线必须）

| 接口 | 说明 | 文档 |
|------|------|------|
| /openapi/developer/live/smallPlay/bind | 绑定、解绑、查绑定状态 | 基础能力文档 |
| /openapi/developer/live/smallPlay/round | 开启、结束对局 | 基础能力文档 |
| /openapi/developer/live/smallPlay/gift | 礼物置顶 | 基础能力文档 |

## 二、加固能力（必须）

| 接口 | 说明 | 文档 |
|------|------|------|
| /openapi/developer/live/data/interactive/pushdata/query | 推送消息回查 | 加固能力文档 |
| /openapi/developer/live/data/interactive/ack/receive | cp客户端收到消息ack | 加固能力文档 |
| /openapi/developer/live/data/interactive/ack/show | cp客户端展示消息ack | 加固能力文档 |

## 三、扩展能力

| 接口 | 说明 | 文档 |
|------|------|------|
| /openapi/developer/live/data/interactive/start | 开启快捷加入战队面板 | 快捷加入战队组件 |
| /openapi/developer/live/data/interactive/action/chat | 连线接口 | 主播双人连屏 |
| /openapi/developer/live/smallPlay/gift | 礼物效果提示 | 礼物效果提示 |
| /openapi/developer/live/smallPlay/audienceInfo | 间内用户信息查询 | 用户数据开放 |
| /openapi/developer/live/data/interactive/action/dataSyncControl | 游戏数据同步控制 | 游戏数据同步 |
| /openapi/developer/live/data/interactive/action/dataSyncReport | 游戏数据同步报告 | 游戏数据同步 |

## 四、关键负责人

- 刘长鑫 - 绑定、对局
- 尹云飞 - 礼物置顶
- 贾文武 - 加固能力
- 李永杰 - 快捷加入、连屏
- 宋家强 - 用户数据
- 许京乐 - 游戏数据同步
