package main


/*:doc
@Id 1.1.1
@Title 一个标题2
@Method get
@Router /list
@Raw
${skip as abddd} | int | 跳过
@Query
a | int | test
b | int | test
@Result
list.${followUpStatus} | int | 跟进状态 1:无需求 2:有需求 3:已预约 4:已签约 9:未接通
${skip as a}? | int | 跳过
a | int | bbb
${skip}? | int | 跳过
*/
func main(){}

/*:doc
@Id 1.1.2
@Title 一个标题2
@Method get
@Router /list
@Raw
${skip as abddd} | int | 跳过
@Query
a | int | test
b | int | test
@Result
list.${followUpStatus} | int | 跟进状态 1:无需求 2:有需求 3:已预约 4:已签约 9:未接通
${skip as a}? | int | 跳过
a | int | bbb
${skip}? | int | 跳过
*/