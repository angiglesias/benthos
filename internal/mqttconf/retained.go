package mqttconf

import (
	"encoding/json"
	"fmt"

	"github.com/Jeffail/benthos/v3/internal/docs"
	"gopkg.in/yaml.v3"
)

func DefaultRetained() Retained {
	return Retained{isBool: true}
}

// Custom Retained to allow interpolation
type Retained struct {
	isBool  bool
	strVal  string
	boolVal bool
}

func NewRetained(val interface{}) (Retained, error) {
	switch value := val.(type) {
	case bool:
		return Retained{isBool: true, boolVal: value}, nil
	case string:
		return Retained{strVal: value}, nil
	default:
		return Retained{}, fmt.Errorf("value of type %T is not accepted", val)
	}
}

func (r Retained) IsBool() bool { return r.isBool }

func (r Retained) BoolVal() bool {
	if !r.isBool {
		return false
	}
	return r.boolVal
}

func (r Retained) StrVal() string {
	if r.isBool {
		return ""
	}
	return r.strVal
}

func (r Retained) MarshalJSON() ([]byte, error) {
	if r.isBool {
		return json.Marshal(r.boolVal)
	}
	return json.Marshal(r.strVal)
}

func (r *Retained) UnmarshalJSON(data []byte) error {
	// try to unmarshal json to bool
	if err := json.Unmarshal(data, &r.boolVal); err == nil {
		r.isBool = true
		return nil
	}
	// fallback to unmarshal a string
	r.isBool = false
	return json.Unmarshal(data, &r.strVal)
}

func (r Retained) MarshalYAML() (interface{}, error) {
	node := &yaml.Node{Kind: yaml.ScalarNode}
	if r.isBool {
		return node, node.Encode(r.boolVal)
	}
	return node, node.Encode(r.strVal)
}

func (r *Retained) UnmarshalYAML(value *yaml.Node) error {
	// try to unmarshal yaml node to bool
	if err := value.Decode(&r.boolVal); err == nil {
		r.isBool = true
		return nil
	}
	// fallback to unmarshal a string
	r.isBool = false
	return value.Decode(&r.strVal)
}

// RetainedFieldSpec defines a retained message flag
func RetainedFieldSpec() docs.FieldSpec {
	return docs.FieldCommon("retained", "Set message as retained on the topic.").IsInterpolated()
}
