# annotation2md

### 通过代码注释生成markdown形式的api文档

1. 下载代码
2. $ go install
3. $ annotation2md --dir=注释所在的目录 --title=文档标题 --outfile=文档名 --vars=声明变量的文档

另：如果需要json格式的文档，可在第三步的参数最后加上--json，即可

#### 生成文档

**例：**

```
/*:doc
@Id 9.9.9
@Title 查找一个列表
@Method get
@Router /users/{which}/list
@Headers
Content-Type | application/json
Authorization | xxxxxx | JWT
@Query
skip? | int | 10 | skip
limit? | int | 10 | limit
@Params
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

**Headers**

|Key|Value|Description|
|-|-|-|
|Content-Type|application/json||
|Authorization|xxxxxx|JWT|

**Params**

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

**注释里的@：**

- **@Id** 文档的id，必须唯一，但不一定要写
- **@Title** 文档标题，必须
- **@Method** 接口的请求方式，必须
- **@Router** 接口地址，必须
- **@Headers** 头信息
- **@Params** 接口地址中的参数
- **@Query** 接口地址中的?后面参数
- **@Body** body中的参数
- **@Result** 返回数据

参数可以在@后方用空格隔开写，也可以换行写，但是每个参数都是必须换行的



#### dir

一般我们将注释写在controller层，所以dir一般指向controller的目录，程序会找该目录下的所有文件（不包含子目录），并将`/*:doc`开头到`*/`结束的这段注释解析成文档

#### title

文档的标题

#### outfile

导出的目录，不需要加后缀，程序会自动加上`.md`后缀，且可以是目录加文件名的形式，比如`public/wiki`

#### vars

变量声明文件的地址，例`conf/.docvars`

**变量的写法：**

例： 
```
$skip | int | 10 | 跳过条目数
$limit | int | 10 | 条目数量
$count | int | 条目总数
```
用2个竖线隔开，第一个为变量名，以`$`开头，中间的中变量类型，最后是描述，一行一个变量

**使用：**

有了变量之后，我们可以将开头的例子文档改为
```
/*:doc
...
@Query
${skip}?
${limit}?
...
*/
```

也可以对变量使用别名，比如我在变量声明文件中声明了一个
```
$someThingStatus | int | 1. 状态A 2. 状态B 3. 状态C
```

在使用时
```
/*:doc
...
@Query
${someThingStatus as status}?
...
*/
```
这样在生成时，声明文件取的是someThingStatus，而生成的文档为status