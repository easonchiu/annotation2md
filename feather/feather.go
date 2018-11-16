package feather

import (
  "io/ioutil"
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

// 获取注释中的文档
func GetDocAnnotationFormFile(filename string) []string {
  fileBytes, err := ioutil.ReadFile(filename)
  annotationDocs := make([]string, 0, 20)

  if err != nil {
    return annotationDocs
  }

  // 找到注释文档，即：以/*:doc开头，以*/结尾的这段注释
  reg := regexp.MustCompile(`/\*:doc((\s|.)*?)\*/`)
  docAnnotationMatchs := reg.FindAllSubmatch(fileBytes, -1)

  // 把匹配结果中的注释都整理出来
  for _, s := range docAnnotationMatchs {
    annotationDocs = append(annotationDocs, string(s[1]))
  }

  return annotationDocs
}
