### elasticsearch操作

#### 创建索引

```bash
PUT magnet
{
  "settings": {
    "analysis": {
      "analyzer": {
        "my_analyzer": {
          "tokenizer": "ik_smart",
          "filter": ["lowercase", "my_filter"]
        }
      },
      "filter": {
        "my_filter": {
          "type": "word_delimiter",
          "preserve_original": true
        }
      }
    }
  },
  "mappings" : {
    "properties" : {
      "addedTime" : {
        "type" : "date",
        "format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
      },
      "files" : {
        "type" : "nested",
        "properties" : {
          "name" : {
            "type" : "text",
            "analyzer": "my_analyzer",
            "fields" : {
              "keyword" : {
                "type" : "keyword",
                "ignore_above" : 256
              }
            }
          },
          "size" : {
            "type" : "long"
          }
        }
      },
      "magnet" : {
        "type" : "text",
        "fields" : {
          "keyword" : {
            "type" : "keyword",
            "ignore_above" : 256
          }
        }
      },
      "name" : {
        "type" : "text",
        "analyzer": "my_analyzer",
        "fields" : {
          "keyword" : {
            "type" : "keyword",
            "ignore_above" : 256
          }
        }
      },
      "size" : {
        "type" : "long"
      },
      "torrent" : {
        "type" : "text",
        "fields" : {
          "keyword" : {
            "type" : "keyword",
            "ignore_above" : 256
          }
        }
      },
      "url" : {
        "type" : "text",
        "fields" : {
          "keyword" : {
            "type" : "keyword",
            "ignore_above" : 256
          }
        }
      }
    }
  }
}
```

#### 查询索引映射

```bash
GET magnet/_mapping
```

#### 删除索引

```bash
DELETE magnet
```

#### 复制索引

```bash
POST _reindex
{
  "source": {
    "index": "magnetbackup"
  },
  "dest": {
    "index": "magnet"
  }
}
```

#### 查询文档数量

```bash
GET magnet/_count
```

#### 查询全部文档

```bash
GET magnet/_search
{
  "query": {
    "match_all": {}
  }
}
```

##### should查询

```bash
GET magnet/_search
{
  "query": {
    "bool": {
      "should": [
        { "match": { "name": "alexis" } },
        {
          "nested": {
            "path": "files",
            "query": {
              "match": { "files.name": "anna" }
            }
          }
        }
      ]
    }
  }
}
```
#### 通配符查询

```bash
GET magnet/_search
{
  "query": {
    "wildcard": {"torrent.keyword": "*rarbgprx*"}
  }
}
```
#### 测试分词

```bash
POST magnet/_analyze
{
  "analyzer": "my_analyzer",
  "text":     "SexArt.23.03.17.Alexis.Crystal.And.Ryana.Fondness.XXX.SD.MP4-KLEENEX"
}
```

---

### elasticsearch数据导入导出

#### 导出数据

```bash
elasticdump --input=http://localhost:9200/magnet --output=magnet.json --type=data
```

#### 导入数据 

```bash
elasticdump --input=magnet.json --output=http://localhost:9200/magnet --type=data
```

---

### 运维

* 编译

```bash
go build -o releases/sukebei src/main.go
```

* 运行

```bash
nohup releases/sukebei > logs/sukebei.log 2>&1 &
```

* 查看日志

```bash
tail -f logs/sukebei.log
```

* 查看进程

```bash
ps aux | grep releases/sukebei
```
