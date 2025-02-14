package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

type GeneralResult struct {
	Kind   string `json:"kind"`
	ErrMsg string `json:"errMsg"`
}

type TestCouponsResponse struct {
	ErrMsg string            `json:"errMsg"`
	Data   []TestCouponsData `json:"data"`
}

type TestCouponsData struct {
	Name        string `json:"name"`
	Amount      int32  `json:"amount"`
	Left        int32  `json:"left"`
	Stock       int32  `json:"stock"`
	Description string `json:"description"`
}

// 任务1
func TestRegisterUser(t *testing.T) {
	var result GeneralResult
	router := setupRouter()
	// 创建客户
	w := httptest.NewRecorder()
	jsonStr := []byte(`{"username": "customer-test", "password": "123", "kind": "customer"}`)
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	result = GeneralResult{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "", result.ErrMsg)
	// 不填kind类型——默认创建客户
	w = httptest.NewRecorder()
	jsonStr = []byte(`{"username": "customer-test1", "password": "123"}`)
	req, _ = http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	result = GeneralResult{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "", result.ErrMsg)
	// 创建商家
	w = httptest.NewRecorder()
	jsonStr = []byte(`{"username": "saler-test", "password": "123456", "kind": "saler"}`)
	req, _ = http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	result = GeneralResult{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "", result.ErrMsg)
	// 检查新用户是否创建成功
	assert.True(t, isUserExist("customer-test"))
	assert.True(t, isUserExist("customer-test1"))
	assert.True(t, isUserExist("saler-test"))
	// 重复创建新用户
	w = httptest.NewRecorder()
	jsonStr = []byte(`{"username": "customer-test", "password": "123", "kind": "customer"}`)
	req, _ = http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	result = GeneralResult{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "用户已存在", result.ErrMsg)
	// 使用错误的kind类型，创建新用户
	jsonStr = []byte(`{"username": "customer-test", "password": "123", "kind": "customerxxx"}`)
	req, _ = http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	result = GeneralResult{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "错误kind类型", result.ErrMsg)
	// 使用空用户名创建用户
	jsonStr = []byte(`{"username": "", "password": "123", "kind": "customerxxx"}`)
	req, _ = http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	result = GeneralResult{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "空用户名或密码", result.ErrMsg)
	// 使用空的密码创建用户
	jsonStr = []byte(`{"username": "customer-test", "password": "", "kind": "customerxxx"}`)
	req, _ = http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	result = GeneralResult{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "空用户名或密码", result.ErrMsg)
}

// 任务1
func TestUserLogin(t *testing.T) {
	var result GeneralResult
	var header string
	router := setupRouter()
	// 空用户和空密码认证
	w := httptest.NewRecorder()
	jsonStr := []byte(`{"username": "", "password": ""}`)
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
	result = GeneralResult{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "空用户或空密码", result.ErrMsg)
	// 不存在的用户认证
	w = httptest.NewRecorder()
	jsonStr = []byte(`{"username": "no-exist", "password": "123"}`)
	req, _ = http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
	result = GeneralResult{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "用户不存在", result.ErrMsg)
	// 错误密码认证
	w = httptest.NewRecorder()
	jsonStr = []byte(`{"username": "customer-test", "password": "1234"}`)
	req, _ = http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
	result = GeneralResult{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "错误密码", result.ErrMsg)
	// 客户认证
	w = httptest.NewRecorder()
	jsonStr = []byte(`{"username": "customer-test", "password": "123"}`)
	req, _ = http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	result = GeneralResult{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "customer", result.Kind)
	assert.Equal(t, "", result.ErrMsg)
	header = w.Header().Get("Authorization")
	assert.NotEqual(t, "", header) // 非空检查
	// 检查是否Authorization
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/validate", nil)
	req.Header.Set("Authorization", header)
	router.ServeHTTP(w, req)
	result = GeneralResult{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "valid", result.ErrMsg)
	// 商家认证
	w = httptest.NewRecorder()
	jsonStr = []byte(`{"username": "saler-test", "password": "123456"}`)
	req, _ = http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	result = GeneralResult{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "saler", result.Kind)
	assert.Equal(t, "", result.ErrMsg)
	header = w.Header().Get("Authorization")
	assert.NotEqual(t, "", header) // 非空检查
	// 检查是否Authorization
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/validate", nil)
	req.Header.Set("Authorization", header)
	router.ServeHTTP(w, req)
	result = GeneralResult{}
	assert.Nil(t, json.NewDecoder(w.Body).Decode(&result))
	assert.Equal(t, "valid", result.ErrMsg)
}

// 任务2
func testCreateCoupons(t *testing.T) {

}

// 任务2
func testGetCouponsInformation(t *testing.T) {

}

// 任务3
// @Return
// Param1: 返回1，若异常状态为5xx；否则返回0
// Param2: 返回响应所耗时间
// 待后续并发进程调用
func testPatchCoupons(t *testing.T) (int, int64) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("Patch", "/api/users/:username/coupons/:name", nil)
	router.ServeHTTP(w, req)
	// 结束
	defer func() {
		err := w.Result().Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	var err5XX int
	var time int64 = 0
	if w.Result().StatusCode/500 > 0 {
		err5XX = 1
	} else {
		err5XX = 0
	}
	assert.Equal(t, 504, w.Code)
	return err5XX, time
}
