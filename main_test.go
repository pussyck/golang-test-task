package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) SMembers(key string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRedisClient) HGetAll(key string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRedisClient) Set(key string, value interface{}) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockRedisClient) SAdd(key string, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockRedisClient) SInter(keys ...string) ([]string, error) {
	args := m.Called(keys)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisClient) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func TestProcessJSON(t *testing.T) {
	mockRedisClient := new(MockRedisClient)

	data := []byte(`[
		{"global_id": 1, "Mode": "test", "ID": 123},
		{"global_id": 2, "Mode": "test", "ID": 456}
	]`)

	mockRedisClient.On("Set", "1", mock.Anything).Return(nil)
	mockRedisClient.On("SAdd", "index:mode:test", "1").Return(nil)
	mockRedisClient.On("SAdd", "index:id:123", "1").Return(nil)

	mockRedisClient.On("Set", "2", mock.Anything).Return(nil)
	mockRedisClient.On("SAdd", "index:mode:test", "2").Return(nil)
	mockRedisClient.On("SAdd", "index:id:456", "2").Return(nil)

	err := storage.ProcessJSON(data, mockRedisClient)

	assert.NoError(t, err)
	mockRedisClient.AssertExpectations(t)
}

func TestProcessFile(t *testing.T) {
	mockRedisClient := new(MockRedisClient)

	data := []byte(`[
		{"global_id": 1, "Mode": "test", "ID": 123}
	]`)

	file := bytes.NewReader(data)

	mockRedisClient.On("Set", "1", mock.Anything).Return(nil)
	mockRedisClient.On("SAdd", "index:mode:test", "1").Return(nil)
	mockRedisClient.On("SAdd", "index:id:123", "1").Return(nil)

	err := storage.ProcessFile(file, mockRedisClient)

	assert.NoError(t, err)
	mockRedisClient.AssertExpectations(t)
}

func TestProcessURL(t *testing.T) {
	mockRedisClient := new(MockRedisClient)

	// Создание тестового HTTP-сервера
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[
			{"global_id": 1, "Mode": "test", "ID": 123}
		]`))
	}))
	defer server.Close()

	mockRedisClient.On("Set", "1", mock.Anything).Return(nil)
	mockRedisClient.On("SAdd", "index:mode:test", "1").Return(nil)
	mockRedisClient.On("SAdd", "index:id:123", "1").Return(nil)

	err := storage.ProcessURL(server.URL, mockRedisClient)

	assert.NoError(t, err)
	mockRedisClient.AssertExpectations(t)
}

func TestProcessZip(t *testing.T) {
	mockRedisClient := new(MockRedisClient)

	// Создание ZIP-архива в памяти
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	fileWriter, _ := zipWriter.Create("data.json")
	fileWriter.Write([]byte(`[
		{"global_id": 1, "Mode": "test", "ID": 123}
	]`))
	zipWriter.Close()

	mockRedisClient.On("Set", "1", mock.Anything).Return(nil)
	mockRedisClient.On("SAdd", "index:mode:test", "1").Return(nil)
	mockRedisClient.On("SAdd", "index:id:123", "1").Return(nil)

	err := storage.ProcessZip(buf.Bytes(), mockRedisClient)

	assert.NoError(t, err)
	mockRedisClient.AssertExpectations(t)
}

func TestSearchData(t *testing.T) {
	mockRedisClient := new(MockRedisClient)

	// Ожидаем вызов SInter с переменным количеством строк
	mockRedisClient.On("SInter", mock.AnythingOfType("[]string")).Return([]string{"1", "2"}, nil)
	mockRedisClient.On("Get", "1").Return(`{"global_id": 1, "Mode": "test", "ID": 123}`, nil)
	mockRedisClient.On("Get", "2").Return(`{"global_id": 2, "Mode": "test", "ID": 456}`, nil)

	// Вызов функции
	records, err := storage.SearchData("", "test", "", mockRedisClient)

	// Проверка результата
	assert.NoError(t, err)
	assert.Len(t, records, 2)

	var record1 map[string]interface{}
	var record2 map[string]interface{}
	json.Unmarshal([]byte(`{"global_id": 1, "Mode": "test", "ID": 123}`), &record1)
	json.Unmarshal([]byte(`{"global_id": 2, "Mode": "test", "ID": 456}`), &record2)

	assert.Equal(t, record1, records[0])
	assert.Equal(t, record2, records[1])
	mockRedisClient.AssertExpectations(t)
}

func TestSearchDataWithGlobalID(t *testing.T) {
	mockRedisClient := new(MockRedisClient)

	mockRedisClient.On("Get", "1").Return(`{"global_id": 1, "Mode": "test", "ID": 123}`, nil)

	records, err := storage.SearchData("1", "", "", mockRedisClient)

	assert.NoError(t, err)
	assert.Len(t, records, 1)

	var record map[string]interface{}
	json.Unmarshal([]byte(`{"global_id": 1, "Mode": "test", "ID": 123}`), &record)

	assert.Equal(t, record, records[0])
	mockRedisClient.AssertExpectations(t)
}

func TestSearchDataNoResults(t *testing.T) {
	mockRedisClient := new(MockRedisClient)

	mockRedisClient.On("SInter", mock.AnythingOfType("[]string")).Return([]string{}, nil)

	records, err := storage.SearchData("", "test", "", mockRedisClient)

	assert.NoError(t, err)
	assert.Len(t, records, 0)
	mockRedisClient.AssertExpectations(t)
}
