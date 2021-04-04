package dynamodbsaas

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/ddelizia/saasaas/pkg/errors"
)

func DeserializeKey(key string) (map[string]*dynamodb.AttributeValue, error) {
	if key == "" {
		return nil, errors.NewTraceable("", "empty key cannot be deserialized", nil)
	}

	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, errors.NewTraceable("", fmt.Sprintf("cannot decode with base64 string [%s]", key), err)
	}

	var result map[string]*dynamodb.AttributeValue
	err = json.Unmarshal([]byte(string(decoded)), &result)

	if err != nil {
		return nil, errors.NewTraceable("", fmt.Sprintf("cannot unmarshal string [%s]", string(decoded)), err)
	}

	return result, nil
}

func SerializeKey(data map[string]*dynamodb.AttributeValue) (string, error) {
	if data == nil {
		return "", nil
	}
	bb, err := json.Marshal(data)
	if err != nil {
		return "", errors.NewTraceable("", "cannot marshal data", err)
	}

	return base64.StdEncoding.EncodeToString(bb), nil
}
