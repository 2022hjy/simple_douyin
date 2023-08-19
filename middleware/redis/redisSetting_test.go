package redis

import (
	"testing"
)

func TestSetValueWithRandomExp(t *testing.T) {
	InitRedis()

	key := "test-key-random"
	value := "test-value-random"

	err := SetValueWithRandomExp(Clients.Test, key, value)
	if err != nil {
		t.Errorf("Error setting value: %v", err)
	}

	val, err := GetValue(Clients.Test, key)
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	} else if val != value {
		t.Errorf("Expected value %s, got %s", value, val)
	}
}

func TestSetValue(t *testing.T) {
	InitRedis()

	key := "test-key-set"
	value := "test-value-set"

	err := SetValue(Clients.Test, key, value)
	if err != nil {
		t.Errorf("Error setting value: %v", err)
	}

	val, err := GetValue(Clients.Test, key)
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	} else if val != value {
		t.Errorf("Expected value %s, got %s", value, val)
	}
}

func TestGetKeysAndUpdateExpiration(t *testing.T) {
	InitRedis()

	key := "test-key-get-update"
	value := "test-value-get-update"

	err := SetValue(Clients.Test, key, value)
	if err != nil {
		t.Errorf("Error setting value: %v", err)
	}

	valInterface, err := GetKeysAndUpdateExpiration(Clients.Test, key)
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}

	val, ok := valInterface.(string)
	if !ok {
		t.Errorf("Expected string value, got %T", valInterface)
	} else if val != value {
		t.Errorf("Expected value %s, got %s", value, val)
	}
}

func TestDeleteKey(t *testing.T) {
	InitRedis()

	key := "test-key-delete"
	value := "test-value-delete"

	err := SetValue(Clients.Test, key, value)
	if err != nil {
		t.Errorf("Error setting value: %v", err)
	}

	err = DeleteKey(Clients.Test, key)
	if err != nil {
		t.Errorf("Error deleting key: %v", err)
	}

	_, err = GetValue(Clients.Test, key)
	if err != NilError {
		t.Errorf("Expected key to be deleted, but got value: %v", err)
	}
}

func TestIsKeyExist(t *testing.T) {
	InitRedis()

	key := "test-key-exist"
	value := "test-value-exist"

	err := SetValue(Clients.Test, key, value)
	if err != nil {
		t.Errorf("Error setting value: %v", err)
	}

	err, exists := isKeyExist(Clients.Test, key)
	if err != nil {
		t.Errorf("Error checking key existence: %v", err)
	} else if !exists {
		t.Errorf("Expected key to exist, but it doesn't")
	}
}
