package parser

import (
  "fmt"
  "log"
  "regexp"
  "strings"
)

// 传参条目的结构体
type RequestData struct {
  Name        string
  Required    bool
  Type        string
  Default     string
  Description string
}

// 返回条目的结构体
type ResultData struct {
  Name        string
  Type        string
  Description string
}

// 文档结构体
type Doc struct {
  Id          string
  Title       string
  Description string
  Method      string
  Router      string
  Params      []*RequestData // url中的参数
  Queries     []*RequestData // ?后面的参数
  Raws        []*RequestData // body中的raw参数
  Result      []*ResultData  // 返回的数据
}

// 注释代码转换成doc struct
func Annotation2DocStruct(annotation string) *Doc {
  doc := Doc{
    Id:          "",
    Title:       "",
    Description: "",
    Method:      "get",
    Router:      "",
    Params:      make([]*RequestData, 0, 10),
    Queries:     make([]*RequestData, 0, 10),
    Raws:        make([]*RequestData, 0, 10),
    Result:      make([]*ResultData, 0, 10),
  }

  // 找到@开头的每一段
  reg := regexp.MustCompile(`@(\S+)([^@]+)`)
  regResult := reg.FindAllStringSubmatch(annotation, -1)

  for _, item := range regResult {
    key := strings.TrimSpace(item[1])
    value := strings.TrimSpace(item[2])
    switch key {
    case "Id":
      doc.Id = value
    case "Title":
      doc.Title = value
    case "Description":
      fallthrough
    case "Desc":
      doc.Description = value
    case "Method":
      doc.Method = value
    case "Router":
      doc.Router = value
    case "Param":
      doc.Params = append(doc.Params, ParseRequestData(value)...)
    case "Query":
      doc.Queries = append(doc.Queries, ParseRequestData(value)...)
    case "Raw":
      doc.Raws = append(doc.Raws, ParseRequestData(value)...)
    case "Result":
      doc.Result = append(doc.Result, ParseResultData(value)...)
    }
  }

  return &doc
}

// 解析RequestData
func ParseRequestData(origin string) []*RequestData {
  reg := regexp.MustCompile(`\n`)
  lines := reg.Split(origin, -1)
  data := make([]*RequestData, 0, len(lines))

  for _, line := range lines {
    var (
      name        string
      required    bool
      typ         string
      def         string
      description string
    )

    reg := regexp.MustCompile(`\|`)
    split := reg.Split(line, -1)

    if len(split) == 0 {
      return nil
    }

    // 去空格
    for i, s := range split {
      split[i] = strings.TrimSpace(s)
    }

    switch len(split) {
    case 4: // name type default description
      name = split[0]
      typ = split[1]
      def = split[2]
      description = split[3]
    case 3: // name type description
      name = split[0]
      typ = split[1]
      description = split[2]
    case 2: // name type
      name = split[0]
      typ = split[1]
    case 1: // name
      name = split[0]
    default:
      return nil
    }

    // 判断key是否必填
    reg = regexp.MustCompile(`\?$`)
    required = !reg.MatchString(name)
    name = reg.Split(name, -1)[0]

    data = append(data, &RequestData{
      Name:        name,
      Required:    required,
      Type:        typ,
      Default:     def,
      Description: description,
    })
  }
  return data
}

// 解析
func ParseResultData(origin string) []*ResultData {

  reg := regexp.MustCompile(`\n`)
  lines := reg.Split(origin, -1)
  data := make([]*ResultData, 0, len(lines))
  for _, line := range lines {
    var (
      name        string
      typ         string
      description string
    )

    reg := regexp.MustCompile(`\|`)
    split := reg.Split(line, -1)

    if len(split) == 0 {
      break
    }

    // 去空格
    for i, s := range split {
      split[i] = strings.TrimSpace(s)
    }

    name = split[0]
    if len(split) > 1 {
      typ = split[1]
    }
    if len(split) > 2 {
      description = split[2]
    }

    data = append(data, &ResultData{name, typ, description})
  }
  return data
}

// Doc struct解析成markdown格式的string
func DocStruct2Markdown(data *Doc) string {
  md := ""
  if data.Title == "" || data.Router == "" {
    log.Fatal("@Title或@Router不能为空，请检查")
    return md
  }

  // 标题
  if data.Id != "" {
    md += "## " + data.Id + ". " + data.Title + "\n"
  } else {
    md += "## " + data.Title + "\n"
  }

  // 描述
  if data.Description != "" {
    md += "**" + data.Description + "**\n"
  }

  // 路由和方法
  md += "\n|URL|Method|\n|-|-|\n"
  md += "|" + data.Router + "|" + strings.ToUpper(data.Method) + "|\n"

  // param
  if len(data.Params) > 0 {
    md += "\n**Param**\n\n"
    md += "|Name|Type|Required|Default|Description|\n|-|-|-|-|-|\n"
    for _, p := range data.Params {
      md += fmt.Sprintf("|%s|%s|%v|%s|%s|\n",
        p.Name, p.Type, p.Required, p.Default, p.Description)
    }
  }

  // query
  if len(data.Queries) > 0 {
    md += "\n**Query**\n\n"
    md += "|Name|Type|Required|Default|Description|\n|-|-|-|-|-|\n"
    for _, p := range data.Queries {
      md += fmt.Sprintf("|%s|%s|%v|%s|%s|\n",
        p.Name, p.Type, p.Required, p.Default, p.Description)
    }
  }

  // raw
  if len(data.Raws) > 0 {
    md += "\n**Raw**\n\n"
    md += "|Name|Type|Required|Default|Description|\n|-|-|-|-|-|\n"
    for _, p := range data.Raws {
      md += fmt.Sprintf("|%s|%s|%v|%s|%s|\n",
        p.Name, p.Type, p.Required, p.Default, p.Description)
    }
  }

  // result
  if len(data.Result) > 0 {
    md += "\n**Result**\n\n"
    md += "|Name|Type|Description|\n|-|-|-|\n"
    for _, p := range data.Result {
      md += fmt.Sprintf("|%s|%s|%s|\n",
        p.Name, p.Type, p.Description)
    }
  }

  return md
}
