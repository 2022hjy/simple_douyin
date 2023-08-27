package redis

import (
	"testing"
	"time"
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

func TestGetKeysAndUpdateExpiration_v1(t *testing.T) {
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

// 这个函数对批量插入以相同key的前提下，后续放入 Value 会覆盖前面的 Value，因此如果需要批量插入，传入的 Value 应该是一个数组
//func TestGetKeysAndUpdateExpiration_v2(t *testing.T) {
//	InitRedis()
//
//	key := "test-key-get-update"
//	value_1 := "test-value_1-get-update-1"
//	value_2 := "test-value_2-get-update-2"
//	value_3 := "test-value_3-get-update-3"
//	value_4 := "test-value_4-get-update-4"
//
//	err := SetValueWithRandomExp(Clients.Test, key, value_1)
//	if err != nil {
//		t.Errorf("Error setting value_1: %v", err)
//	}
//	err2 := SetValueWithRandomExp(Clients.Test, key, value_2)
//	if err2 != nil {
//		t.Errorf("Error setting value_2: %v", err2)
//	}
//	err3 := SetValueWithRandomExp(Clients.Test, key, value_3)
//	if err3 != nil {
//		t.Errorf("Error setting value_3: %v", err3)
//	}
//	err4 := SetValueWithRandomExp(Clients.Test, key, value_4)
//	if err4 != nil {
//		t.Errorf("Error setting value_4: %v", err4)
//	}
//
//	valInterface, err := GetKeysAndUpdateExpiration(Clients.Test, key)
//	log.Printf("valInterface:%v\n", valInterface)
//
//	if err != nil {
//		t.Errorf("Error getting value_1: %v", err)
//	}
//
//	val, ok := valInterface.(string)
//	if !ok {
//		t.Errorf("Expected string value_1, got %T", valInterface)
//	} else if val != value_1 {
//		t.Errorf("Expected value_1 %s, got %s", value_1, val)
//	}
//}

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

	exists, err := IsKeyExist(Clients.Test, key)
	if err != nil {
		t.Errorf("Error checking key existence: %v", err)
	} else if !exists {
		t.Errorf("Expected key to exist, but it doesn't")
	}
}

func TestSetHashWithExpiration(t *testing.T) {
	InitRedis()
	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	err := SetHashWithExpiration(Clients.Test, "test-hash", data, 2*time.Minute)
	if err != nil {
		t.Errorf("Error setting hash: %v", err)
	}
	if err != nil {
		t.Errorf("Error getting hash: %v", err)
	}
}
