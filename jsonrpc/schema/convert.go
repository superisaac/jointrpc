package schema

// util functions
func convertMap(node map[string]interface{}, attrName string, optional bool) (map[string]interface{}, bool) {
	if v, ok := node[attrName]; ok {
		if aMap, ok := v.(map[string]interface{}); ok {
			return aMap, ok
		}
	} else if optional {
		return make(map[string]interface{}), true
	}
	return nil, false
}

func convertList(node map[string]interface{}, attrName string, optional bool) ([]interface{}, bool) {
	if v, ok := node[attrName]; ok {
		if aList, ok := v.([]interface{}); ok {
			return aList, ok
		}
	} else if optional {
		return [](interface{}){}, true
	}
	return nil, false
}

func convertMapOfMap(node map[string](interface{}), attrName string, optional bool) (map[string](map[string]interface{}), bool) {
	if mm, ok := convertMap(node, attrName, optional); ok {
		resMap := make(map[string](map[string]interface{}))
		for name, value := range mm {
			mv, ok := value.(map[string]interface{})
			if !ok {
				return nil, false
			}
			resMap[name] = mv
		}
		return resMap, true
	}
	return nil, false
}

func convertListOfMap(node map[string]interface{}, attrName string, optional bool) ([](map[string]interface{}), bool) {
	if v, ok := node[attrName]; ok {
		if aList, ok := v.([]interface{}); ok {
			arr := make([](map[string]interface{}), 0)
			for _, item := range aList {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					return nil, false
				}
				arr = append(arr, itemMap)
			}
			return arr, true
		}
	} else if optional {
		return [](map[string]interface{}){}, true
	}
	return nil, false
}

func convertListOfString(node map[string]interface{}, attrName string, optional bool) ([]string, bool) {
	if v, ok := node[attrName]; ok {
		if aList, ok := v.([]interface{}); ok {
			arr := make([]string, 0)
			for _, item := range aList {
				strItem, ok := item.(string)
				if !ok {
					return nil, false
				}
				arr = append(arr, strItem)
			}
			return arr, true
		}
	} else if optional {
		return []string{}, true
	}
	return nil, false
}
