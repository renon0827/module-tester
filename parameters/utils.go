package parameters

import (
	"encoding/json"
	"module-tester/constants"
	"module-tester/repository"
	"strings"
)

func dataToString(data interface{}) string {
	if val, err := json.Marshal(data); err == nil {
		return string(val)
	}
	if val, castable := data.(string); castable {
		return val
	}

	return ""
}

func dataSetFromRepository(data interface{}, repo repository.Repository) interface{} {
	var result interface{}
	if m, castable := data.(map[string]interface{}); castable {
		for key := range m {
			m[key] = dataSetFromRepository(m[key], repo)
		}
		result = m
	} else if ary, castable := data.([]interface{}); castable {
		for idx := range ary {
			ary[idx] = dataSetFromRepository(ary[idx], repo)
		}
		result = ary
	} else if v, castable := data.(string); repo != nil && castable &&
		strings.HasPrefix(v, constants.ValiableValuePrefix) && strings.HasSuffix(v, constants.ValiableValueSuffix) {
		v = v[len(constants.ValiableValuePrefix) : len(v)-len(constants.ValiableValueSuffix)]
		if len(v) > 0 {
			repoVal, err := repo.Get(v)
			if err == nil {
				result = repoVal
			}
		}
	} else {
		result = data
	}

	return result
}
