package feather

import (
  "bytes"
  "io/ioutil"
  "log"
  "os"
  "regexp"
  "strings"
  "unicode"
)

// 根据目录找到文件名
func GetFileNamesFromDir(dirname string) []string {
  filenames := make([]string, 0, 10)
  dirs, _ := ioutil.ReadDir(dirname)

  // 整理所有的文件名
  for _, dir := range dirs {
    if !dir.IsDir() {
      d := dirname + "/" + dir.Name()
      filenames = append(filenames, d)
    }
  }
  return filenames
}

// 获取注释中的文档，如果有变量在这里做替换
func GetDocAnnotationFormFile(filename string, vars map[string]string) []string {
  fileBytes, err := ioutil.ReadFile(filename)
  annotationDocs := make([]string, 0, 20)

  if err != nil {
    return annotationDocs
  }

  // 找到注释文档，即：以/*:doc开头，以*/结尾的这段注释
  reg := regexp.MustCompile(`/\*:doc((\s|.)*?)\*/`)
  docAnnotationMatchs := reg.FindAllSubmatch(fileBytes, -1)
  varsreg := regexp.MustCompile(`\${\S+}\??(\s|.)*?\n`)
  varskey := regexp.MustCompile(`\${\S+}\??`)

  // 把匹配结果中的注释都整理出来
  for _, s := range docAnnotationMatchs {
    doc := varsreg.ReplaceAllFunc(s[1], func(b []byte) []byte {
      findString := varskey.FindString(string(b))
      if v, ok := vars[findString]; ok {
        return []byte(v + "\n")
      } else {
        log.Fatal("找不到相关变量：", findString)
        return b
      }
    })
    replaceDoc := varsreg.ReplaceAllFunc(s[1], func(b []byte) []byte {
      findString := varskey.FindString(string(b))
      if v, ok := vars[findString]; ok {
        reg := regexp.MustCompile(`^[^|]*?\|`)
        r := reg.ReplaceAllString(v, "|")
        return []byte(findString + " " + r + "\n")
      }
      return b
    })
    // 替换原文档
    fileBytes = bytes.Replace(fileBytes, s[1], replaceDoc, -1)
    annotationDocs = append(annotationDocs, string(doc))
  }

  if len(annotationDocs) > 0 {
    ioutil.WriteFile(filename, fileBytes, os.ModePerm)
  }

  return annotationDocs
}

// 获取变量文件的变量
func GetKeyVarsFromFile(filename string) map[string]string {
  fileBytes, err := ioutil.ReadFile(filename)
  vars := make(map[string]string)

  if err != nil {
    return vars
  }

  // 找到@开头的每一段
  reg := regexp.MustCompile(`\$(\S+)([^\\$]+)`)
  regResult := reg.FindAllStringSubmatch(string(fileBytes), -1)
  for _, res := range regResult {
    key := res[1]
    val := strings.TrimRightFunc(res[2], unicode.IsSpace)
    vars["${"+key+"}"] = key + val
    vars["${"+key+"}?"] = key + "?" + val
  }

  return vars
}
