excavator -- lib for converting data from one type to another.

## Supported types
Now this convertions are implemented:

* interface{} -> basic types (int, string, whatever)
* map[string]interface{} -> struct
* map[string]interface{} -> map[string]supportedType
* struct -> map[string]supportedType
* []interface{} -> []supportedType

Everything working recursivly, so you can convert `map[string]interface{}` to 
`map[string][]struct` if data is convertable.
