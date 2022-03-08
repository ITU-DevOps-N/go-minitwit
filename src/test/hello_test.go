package hello_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// type TestSuite struct {
// 	suite.Suite
// 	openDB func() gorm.Dialector
// }

// func init() {
// 	// initialize the router
// 	router := gin.Default()
// 	// Handlers for testing
// 	router.POST("/register", main.SignUp)

// }

// func TestWholeTestSuite(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// // This will run before each test in the suite
// func (suite *TestSuite) SetupTestDB() {
// 	db := sqlite.Open("file::memory:?cache=shared")
// 	suite.openDB = func() gorm.Dialector {
// 		return db
// 	}
// 	main.SetupDB(suite.openDB().Name())
// }

// //Helper functions

// // // This is an example test that will always succeed
// func (suite *TestSuite) TestCreateUser() {

// 	// Act
// 	body, _ := json.Marshal(gin.H{
// 		"username": "Yennefer of Vengerberg",
// 		"email":    "yennefer@aretuza.wr",
// 		"pwd":      "chaosmaster",
// 	})
// 	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
// 	fmt.Println(req)
// 	// Assert
// 	suite.Equal(http.StatusNoContent, "")
// }

// func test_register(suite *TestSuite) {
// 	//
// }
