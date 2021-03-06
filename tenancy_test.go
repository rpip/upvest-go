package upvest

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// Tests API call to create, retrieve and delete user
func TestUserCRUD(t *testing.T) {
	uid, _ := uuid.NewUUID()
	username := fmt.Sprintf("upvest_test_%s", uid.String())
	user, err := tenancyTestClient.User.Create(username, randomString(12))
	if err != nil {
		t.Errorf("CREATE User returned error: %v", err)
	}
	if user.Username != username {
		t.Errorf("Expected User username %v, got %v", username, user.Username)
	}
	if user.RecoveryKit == "" {
		t.Errorf("Expected User recovery kit to be set, got nil")
	}

	// retrieve the user
	user1, err := tenancyTestClient.User.Get(user.Username)
	if err != nil {
		t.Errorf("GET User returned error: %v", err)
	}

	if user.Username != user1.Username {
		t.Errorf("Expected User username %v, got %v", user.Username, user1.Username)
	}

	// delete user
	_ = tenancyTestClient.User.Delete(user.Username)
	usr, err := tenancyTestClient.User.Get(user.Username)
	aerr := err.(*Error)

	if aerr.StatusCode != 404 {
		t.Errorf("Expected user not found, got %s", usr.Username)
	}
}

// Tests an API call to get list of users
func TestListUsers(t *testing.T) {
	expected := 10
	users, err := tenancyTestClient.User.List()
	if err != nil {
		t.Errorf("List Users returned error: %v", err)
	}
	if len(users.Values) < expected {
		t.Errorf("Expected greater than %d users, got %d", expected, len(users.Values))
	}
}

// Tests an API call to get list of specific number of users
func TestListNUsers(t *testing.T) {
	expected := 10
	users, err := tenancyTestClient.User.ListN(expected)
	if err != nil {
		t.Errorf("List Users returned error: %v", err)
	}

	if len(users.Values) != expected {
		t.Errorf("Expected greater than %d users, got %d", expected, len(users.Values))
	}
}

// Tests an API call to update a user's password
func TestChangePassword(t *testing.T) {
	user, pw := createTestUser()
	username := user.Username
	params := &ChangePasswordParams{
		OldPassword: pw,
		NewPassword: randomString(12),
	}
	user, _ = tenancyTestClient.User.ChangePassword(username, params)
	if user.Username != username {
		t.Errorf("Expected username %s, got %+v", username, user)
	}
}

// Test to retrive all assets
func TestListAssets(t *testing.T) {
	assets, err := tenancyTestClient.Asset.List()

	if err != nil {
		t.Errorf("List assets returned error: %v", err)
	}

	asset1 := assets.Values[0]
	assertions := []bool{
		asset1.ID == ethARWeaveAsset.ID,
		asset1.Name == ethARWeaveAsset.Name,
		asset1.Symbol == ethARWeaveAsset.Symbol,
		asset1.Exponent == ethARWeaveAsset.Exponent,
		asset1.Protocol == ethARWeaveAsset.Protocol,
	}

	for _, isValid := range assertions {
		if !isValid {
			t.Errorf("Asset structure does not match expected")
		}
	}
}
