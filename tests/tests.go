package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ITU-DevOps-N/go-minitwit"

	"testing"

	model "github.com/ITU-DevOps-N/go-minitwit/models"
	"github.com/stretchr/testify/assert"
)

type TestSuite struct {
	suite.Suite
	openDB func() gorm.Dialector
}

func TestWholeTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// This will run before each test in the suite
func (suite *TestSuite) SetupTestDB() {
	db := sqlite.Open("file::memory:?cache=shared")
	suite.openDB = func() gorm.Dialector {
		return db
	}
	main.SetupDB(suite.openDB().Name())
}

//Helper functions




// // This is an example test that will always succeed
func (suite *TestSuite) TestCreateUser() {

	// Assert
	suite.Equal(http.StatusNoContent, "")
}



func test_register(suite *TestSuite){
	//
}





