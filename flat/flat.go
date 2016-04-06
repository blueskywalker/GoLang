
package main

import (
        "fmt"

)

func flat(list []interface{}) []interface{} {
	var ret []interface{}

	for _,v := range list {
		switch t := v.(type) {
		default:
			ret = append(ret,t)
		case []interface{}:
			for _, e :=range flat(t) {
				ret = append(ret,e)
			}
		}

	}
	return ret
}

func main()  {

        letters := []string{"a", "b", "c", "d","e"}
        fmt.Println(letters)

        var data []interface{}
	for _, value := range letters {
		var tmp []interface{}
		tmp = append(tmp,value)
		if len(data) > 0 {
			tmp = append(tmp, data)
		}
		data = append(data,tmp)
	}
	fmt.Println(data)
	fmt.Println(flat(data))

}
