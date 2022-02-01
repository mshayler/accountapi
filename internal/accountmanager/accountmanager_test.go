package accountmanager

import (
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testUser = "TestUser"
	testPass = "TestPass"
)

var (
	testlgr, err = logging.GetLogger("Account Test Logger")
)

func TestNewManager(t *testing.T) {
	testCases := map[string]struct {
		lgr         *logging.Logger
		expectError bool
	}{
		"New Manager Fail": {
			nil,
			true,
		},
		"New Manager Success": {
			lgr:         testlgr,
			expectError: false,
		},
	}

	for name, v := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tc := v

			c, err := NewManager(tc.lgr)
			if tc.expectError {
				assert.NotNil(t, "Error should not be nil.")
				assert.Nil(t, c, "Manager should be nil.")
			} else {
				assert.Nil(t, err, "Error should be nil. Got: %v", err)
				assert.NotNil(t, c, "Manager should not be nil.")
			}
		})

	}
}

func TestCreateAccount(t *testing.T) {
	testCases := map[string]struct {
		user          string
		pass          string
		removeAccount bool
		expectError   bool
	}{
		"Add Account": {
			user:        testUser,
			pass:        testPass,
			expectError: false,
		},
		"Add Account without Username": {
			user:        "",
			pass:        testPass,
			expectError: true,
		},
		"Add Account without Password": {
			user:        testUser,
			pass:        "",
			expectError: true,
		},
		"Fail Duplicate Account": {
			user:          testUser,
			pass:          testPass,
			removeAccount: true,
			expectError:   true,
		},
	}

	// Create an account Manager
	c, err := NewManager(testlgr)
	if err != nil {
		testlgr.Fatal("Failed to create account Manager.")
	}

	for name, v := range testCases {
		tc := v

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := c.CreateAccount(tc.user, tc.pass)
			if tc.expectError {
				assert.NotNil(t, err, "Error should not be nil.")
			} else {
				assert.Nil(t, err, "Error should be nil. Got: %v", err)
			}
			if tc.removeAccount {
				c.Persistence.DeleteAccount(tc.user)
			}
		})
	}
}

func TestLoginAccount(t *testing.T) {
	testCases := map[string]struct {
		user            string
		pass            string
		token           string
		successfulLogin bool
		expectError     bool
	}{
		"Login without Username": {
			user:            "",
			pass:            testPass,
			successfulLogin: false,
			expectError:     true,
		},
		"Login without Password": {
			user:            testUser,
			pass:            "",
			successfulLogin: false,
			expectError:     true,
		},
		"Login with Incorrect Password": {
			user:            testUser,
			pass:            "fakepassword",
			successfulLogin: false,
			expectError:     true,
		},
		"Login Successful with Token": {
			user:            testUser,
			pass:            testPass,
			successfulLogin: true,
			token:           "d3e1bd0b6137b55c507c2383013c0ef6f556e5d1",
			expectError:     false,
		},
	}

	// Create an account Manager
	c, err := NewManager(testlgr)
	if err != nil {
		testlgr.Fatal("Failed to create account Manager.")
	}

	// Verify can't login without account made
	_, err = c.LoginAccount(testUser, testPass)
	assert.NotNil(t, err, "Error should not be nil for login for non-existing account")

	// create the user now
	_, err = c.CreateAccount(testUser, testPass)
	assert.Nil(t, err, "Creation of account was an error.")

	for name, v := range testCases {
		tc := v

		res, err := c.LoginAccount(tc.user, tc.pass)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// This could be cleaner, implement interface to skip it
			if tc.expectError {
				assert.NotNil(t, err, "Error should not be nil.")
			} else {
				assert.Nil(t, err, "Error should be nil. Got: %v", err)
			}

			if tc.successfulLogin {
				assert.NotEmpty(t, res, "successful login but empty token")
			}

		})
	}
	c.Persistence.DeleteAccount(testUser)
}

func TestVerifyAccount(t *testing.T) {

	// Create an account Manager
	c, err := NewManager(testlgr)
	if err != nil {
		testlgr.Fatal("Failed to create account Manager.")
	}

	// create the user now
	_, err = c.CreateAccount(testUser, testPass)
	assert.Nil(t, err, "Creation of account was an error.")

	// login with correct details
	tkn, err := c.LoginAccount(testUser, testPass)
	assert.Nil(t, err, "Error when logging into account.")

	// Verify with fake token
	_, err = c.VerifyAccount(testUser, "faketoken")
	assert.NotNil(t, err, "Should be error when using fake token")

	// Verify with real token
	res, err := c.VerifyAccount(testUser, tkn)
	assert.Nil(t, err, "Should not be an error with real token")
	assert.True(t, res, "Result should be a true login.")

	c.Persistence.DeleteAccount(testUser)

}
