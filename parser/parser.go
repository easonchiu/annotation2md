package parser

import (
  "fmt"
  "log"
  "regexp"
  "strings"
)

// 头的结构体
type HeadersData struct {
  Key         string
  Value       string
  Description string
}

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
  Required    bool
  Description string
}

// 文档结构体
type Doc struct {
  Id          string
  Title       string
  Description string
  Method      string
  Router      string
  Headers     []*HeadersData // 头
  Params      []*RequestData // url中的参数
  Query       []*RequestData // ?后面的参数
  Body        []*RequestData // body
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
    Headers:     make([]*HeadersData, 0, 10),
    Params:      make([]*RequestData, 0, 10),
    Query:       make([]*RequestData, 0, 10),
    Body:        make([]*RequestData, 0, 10),
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
    case "Headers":
      doc.Headers = append(doc.Headers, ParseHeadersData(value)...)
    case "Params":
      doc.Params = append(doc.Params, ParseRequestData(value)...)
    case "Query":
      doc.Query = append(doc.Query, ParseRequestData(value)...)
    case "Body":
      doc.Body = append(doc.Body, ParseRequestData(value)...)
    case "Result":
      doc.Result = append(doc.Result, ParseResultData(value)...)
    }
  }

  return &doc
}

// 解析HeadersData
func ParseHeadersData(origin string) []*HeadersData {
  reg := regexp.MustCompile(`\n`)
  lines := reg.Split(origin, -1)
  data := make([]*HeadersData, 0, len(lines))

  for _, line := range lines {
    var (
      key         string
      value       string
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
    case 3: // key value description
      key = split[0]
      value = split[1]
      description = split[2]
    case 2: // key value
      key = split[0]
      value = split[1]
    default:
      return nil
    }

    data = append(data, &HeadersData{
      Key:         key,
      Value:       value,
      Description: description,
    })
  }
  return data
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

// 解析result
func ParseResultData(origin string) []*ResultData {
  reg := regexp.MustCompile(`\n`)
  lines := reg.Split(origin, -1)
  data := make([]*ResultData, 0, len(lines))
  for _, line := range lines {
    var (
      name        string
      typ         string
      required    bool
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

    // 判断key是否必填
    reg = regexp.MustCompile(`\?$`)
    required = !reg.MatchString(name)
    name = reg.Split(name, -1)[0]

    data = append(data, &ResultData{name, typ, required, description})
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

  // headers
  if len(data.Headers) > 0 {
    md += "\n**Headers**\n\n"
    md += "|Key|Value|Description|\n|-|-|-|\n"
    for _, p := range data.Headers {
      md += fmt.Sprintf("|%s|%s|%v|\n",
        p.Key, p.Value, p.Description)
    }
  }

  // params
  if len(data.Params) > 0 {
    md += "\n**Params**\n\n"
    md += "|Name|Type|Required|Default|Description|\n|-|-|-|-|-|\n"
    for _, p := range data.Params {
      md += fmt.Sprintf("|%s|%s|%v|%s|%s|\n",
        p.Name, p.Type, p.Required, p.Default, p.Description)
    }
  }

  // query
  if len(data.Query) > 0 {
    md += "\n**Query**\n\n"
    md += "|Name|Type|Required|Default|Description|\n|-|-|-|-|-|\n"
    for _, p := range data.Query {
      md += fmt.Sprintf("|%s|%s|%v|%s|%s|\n",
        p.Name, p.Type, p.Required, p.Default, p.Description)
    }
  }

  // body
  if len(data.Body) > 0 {
    md += "\n**Body**\n\n"
    md += "|Name|Type|Required|Default|Description|\n|-|-|-|-|-|\n"
    for _, p := range data.Body {
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

// 替换冒号
func ReplaceMH(str string) string {
  reg := regexp.MustCompile(`"`)
  return reg.ReplaceAllString(str, `\"`)
}

// Doc struct解析成json格式的string
func DocStruct2JSON(data *Doc) string {
  json := "{"
  if data.Title == "" || data.Router == "" {
    log.Fatal("@Title或@Router不能为空，请检查")
    return ""
  }

  // 标题
  if data.Id != "" {
    json += `"title": "` + ReplaceMH(data.Id) + ". " + ReplaceMH(data.Title) + `"`
  } else {
    json += `"title": "` + ReplaceMH(data.Title) + `"`
  }

  // 描述
  if data.Description != "" {
    json += `,"description": "` + ReplaceMH(data.Description) + `"`
  }

  // 路由和方法
  json += `,"url": "` + ReplaceMH(data.Router) + `"`
  json += `,"method": "` + ReplaceMH(strings.ToUpper(data.Method)) + `"`

  // headers
  if len(data.Headers) > 0 {
    json += `,"headers":[`
    for i, p := range data.Headers {
      json += fmt.Sprintf(`{"key":"%s","value":"%s","description":"%s"}`,
        ReplaceMH(p.Key), ReplaceMH(p.Value), ReplaceMH(p.Description))
      if i < len(data.Headers)-1 {
        json += ","
      }
    }
    json += `]`
  }

  // params
  if len(data.Params) > 0 {
    json += `,"params":[`
    for i, p := range data.Params {
      json += fmt.Sprintf(`{"name":"%s","type":"%s","required":%v,"default":"%s","description":"%s"}`,
        ReplaceMH(p.Name), ReplaceMH(p.Type), p.Required, ReplaceMH(p.Default), ReplaceMH(p.Description))
      if i < len(data.Params)-1 {
        json += ","
      }
    }
    json += `]`
  }

  // query
  if len(data.Query) > 0 {
    json += `,"query":[`
    for i, p := range data.Query {
      json += fmt.Sprintf(`{"name":"%s","type":"%s","required":%v,"default":"%s","description":"%s"}`,
        ReplaceMH(p.Name), ReplaceMH(p.Type), p.Required, ReplaceMH(p.Default), ReplaceMH(p.Description))
      if i < len(data.Query)-1 {
        json += ","
      }
    }
    json += `]`
  }

  // body
  if len(data.Body) > 0 {
    json += `,"body":[`
    for i, p := range data.Body {
      json += fmt.Sprintf(`{"name":"%s","type":"%s","required":%v,"default":"%s","description":"%s"}`,
        ReplaceMH(p.Name), ReplaceMH(p.Type), p.Required, ReplaceMH(p.Default), ReplaceMH(p.Description))
      if i < len(data.Body)-1 {
        json += ","
      }
    }
    json += `]`
  }

  // result
  if len(data.Result) > 0 {
    json += `,"result":[`
    for i, p := range data.Result {
      json += fmt.Sprintf(`{"name":"%s","type":"%s","required":%v,"description":"%s"}`,
        ReplaceMH(p.Name), ReplaceMH(p.Type), p.Required, ReplaceMH(p.Description))
      if i < len(data.Result)-1 {
        json += ","
      }
    }
    json += `]`
  }

  json += "}"
  return json
}
