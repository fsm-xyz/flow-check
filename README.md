# README

监控家庭V网使用情况，及时预警

## 方案

### API

通过调用接口方式，监控

#### go

`当前`: 手动请求接口，得到token，然后服务端不断访问服务器(go), 当token失败的时候，用接口更新token

`目标`: 全部自动化，全部接口自动化请求

#### node

`当前`: 未实现

用nodejs直接复用前端代码

### headless

使用headless浏览器模拟

#### rust

`当前`: rust的http接口复用go api服务的

`目标`：全部rust实现

#### go

`当前`: 已完成

`目标`：已完成

### 内嵌页面

`当前`: 未实现

直接嵌入官方页面，插入js代码(tauri, web)

## 资料

页面: https://sn.clientaccess.10086.cn/html5/sx/vfamilyN/index.html

接口: https://sn.clientaccess.10086.cn/html5/indivbusi/familyNew/getStatus
