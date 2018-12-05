package engine

import (
  "annotation2md/feather"
  "annotation2md/parser"
  "log"
  "sort"
)

func Start(title string, files []string, vars map[string]string) string {
  markdown := ""
  if title != "" {
    markdown = "# " + title + "\n\n"
  }

  docs := make([]*parser.Doc, 0, 50)

  if len(files) > 0 {
    for _, filename := range files {
      docAnnotations := feather.GetDocAnnotationFormFile(filename, vars)
      if len(docAnnotations) == 0 {
        continue
      }
      for _, docAnnotation := range docAnnotations {
        doc := parser.Annotation2DocStruct(docAnnotation)
        docs = append(docs, doc)
      }
    }

    idList := make([]string, 0, len(docs))
    hasIdDocMap := make(map[string]*parser.Doc)

    for _, doc := range docs {
      if doc.Id != "" {
        if _, ok := hasIdDocMap[doc.Id]; ok {
          log.Fatal("存在相同的Id: " + doc.Id + "，请修改后重试")
          return ""
        }
        idList = append(idList, doc.Id)
        hasIdDocMap[doc.Id] = doc
      }
    }

    // 排序
    sort.Strings(idList)

    // 排序后的doc
    sortDocList := make([]*parser.Doc, 0, len(docs))
    // 有id的doc优先排进来
    for _, id := range idList {
      sortDocList = append(sortDocList, hasIdDocMap[id])
    }
    // 把没有id的继续往下排
    for _, doc := range docs {
      if doc.Id == "" {
        sortDocList = append(sortDocList, doc)
      }
    }

    // 处理排好序的列表
    for _, doc := range sortDocList {
      res := parser.DocStruct2Markdown(doc)
      if res != "" {
        markdown += res + "\n"
      }
    }
  }

  return markdown
}