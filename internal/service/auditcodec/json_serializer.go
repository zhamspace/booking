package auditcodec

import (
	"encoding/json"
	"fmt"
)

type JSONSerializer struct{}

func (JSONSerializer) Serialize(resource string, object any) (json.RawMessage, error) {
	if object == nil {
		return nil, fmt.Errorf("audit serializer: object is required")
	}

	switch resource {
	default:
		switch v := object.(type) {
		case map[string]string:
			return json.Marshal(v)
		case map[string]any:
			return json.Marshal(v)
		default:
			return json.Marshal(v)
		}
	}
}

func (JSONSerializer) Deserialize(resource string, data json.RawMessage) (any, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("audit serializer: object data is required")
	}

	switch resource {
	default:
		var v map[string]any
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return v, nil
	}
}
