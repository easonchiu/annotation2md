package feather

import (
  "io/ioutil"
  "log"
  "regexp"
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
  varsreg := regexp.MustCompile(`\${\S+}\??`)

  // 把匹配结果中的注释都整理出来
  for _, s := range docAnnotationMatchs {
    doc := s[1]
    doc = varsreg.ReplaceAllFunc(doc, func(b []byte) []byte {
      if v, ok := vars[string(b)]; ok {
        return []byte(v)
      } else {
        log.Fatal("找不到相关变量：", string(b))
        return b
      }
    })
    annotationDocs = append(annotationDocs, string(doc))
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
    val := res[2]
    vars["${"+key+"}"] = key + val
    vars["${"+key+"}?"] = key + "?" + val
  }

  return vars
}
