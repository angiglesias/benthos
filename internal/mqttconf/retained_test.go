package mqttconf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type configExample struct {
	URLs     []string `json:"urls" yaml:"urls"`
	Retained Retained `json:"retained" yaml:"retained"`
}

func TestYamlMarshal(t *testing.T) {
	rBool := Retained{isBool: true, boolVal: true}
	rStr := Retained{strVal: `${! meta("mqtt_retained") }`}

	ser, err := yaml.Marshal(rBool)
	require.NoError(t, err, "Unexpected error marshaling yaml value")
	assert.YAMLEq(t, "true", string(ser), "Boolean value was incorrectly serialized")
	ser, err = yaml.Marshal(rStr)
	require.NoError(t, err, "Unexpected error marshaling yaml value")
	assert.YAMLEq(t, `${! meta("mqtt_retained") }`, string(ser), "String value was incorrectly serialized")
}

func TestYamlUnmarshal(t *testing.T) {
	var cfgBool, cfgStr configExample

	err := yaml.Unmarshal([]byte(boolYaml), &cfgBool)
	require.NoError(t, err, "Unexpecterd error unmarshaling yaml document")
	assert.True(t, cfgBool.Retained.isBool, "Unmarshalled value should be a boolean")
	assert.True(t, cfgBool.Retained.boolVal, "Incorrect umarshalled value")

	err = yaml.Unmarshal([]byte(strYaml), &cfgStr)
	require.NoError(t, err, "Unexpecterd error unmarshaling yaml document")
	assert.False(t, cfgStr.Retained.isBool, "Unmarshalled value should be a string")
	assert.EqualValues(t, `${! meta("mqtt_retained") }`, cfgStr.Retained.strVal, "Incorrect umarshalled value")
}

const (
	boolYaml = `
urls:
  - nats://localhost:4222
  - nats://remote:4222
retained: true
`

	strYaml = `
urls:
  - nats://localhost:4222
  - nats://remote:4222
retained: ${! meta("mqtt_retained") }
`
)
