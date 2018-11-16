# annotation2md

### 通过代码注释生成markdown形式的api文档

1. 下载代码
2. $ go install
3. $ annotation2md --dir=注释所在的目录 --title=文档标题 --outfile=文档名


**例：**

```
/*:doc
@Id 9.9.9
@Title 查找一个列表
@Method get
@Router /users/{which}/list
@Query
skip? | int | 10 | skip
limit? | int | 10 | limit
@Param
which | string | 随便一个什么参数
@Result
foo | string | foo
bar | int | bar
*/
```

**将生成为：**

## 9.9.9. 查找一个列表

|URL|Method|
|-|-|
|/users/{which}/list|GET|

**Param**

|Name|Type|Required|Default|Description|
|-|-|-|-|-|
|which|string|true||随便一个什么参数|

**Query**

|Name|Type|Required|Default|Description|
|-|-|-|-|-|
|skip|int|false|10|skip|
|limit|int|false|10|limit|

**Result**

|Name|Type|Description|
|-|-|-|
|foo|string|foo|
|bar|int|bar|