## go map for gorutines

### usage

```golang
var xMap map[int]int
var yMap = make(map[string][]string)
var xMapValue int
var yMapValue []string

xxMap := NewMap(xMap)
xxMap.Add(1, 1)
xxMap.Query(1, &xMapValue)
xxMap.Delete(1)
xxMap.Query(1, &xMapValue)

yyMap := NewMap(yMap)
yyMap.Add("slice", []string{"it's", "a", "slice"})
yyMap.Query("slice", &yMapValue)
yyMap.Delete("slice")
yyMap.Query("slice", &yMapValue)
```