# mc-cli
## 安装
* 编译好的binary[下载地址](https://github.com/dingxin-tech/mc-cli/releases/download/v0.0.1-alpha/mc)
* 通过 go 进行编译
  ```bash
   git clone https://github.com/dingxin-tech/mc-cli.git
   go build -o mc
   ```

## 配置
默认配置文件在 ~/.mc.yaml，可通过 --config 指定配置文件位置
格式为
```yaml
end_point: 
project_name: 
access_id: 
access_key: 
```

## 功能
支持且仅支持提交sql
```bash
mc query "select * from table"
```
