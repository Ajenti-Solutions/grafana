package oauthserver

import (
	"context"
	"testing"

	ac "github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/accesscontrol/acimpl"
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/oauthserver/utils"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/stretchr/testify/require"
)

func setupTestEnv(t *testing.T) *Client {
	t.Helper()

	client := &Client{
		Name:             "my-ext-service",
		ClientID:         "RANDOMID",
		Secret:           "RANDOMSECRET",
		GrantTypes:       "client_credentials,urn:ietf:params:oauth:grant-type:jwt-bearer",
		ServiceAccountID: 2,
		SelfPermissions: []ac.Permission{
			{Action: ac.ActionUsersImpersonate, Scope: ac.ScopeUsersAll},
		},
		SignedInUser: &user.SignedInUser{
			UserID: 2,
			OrgID:  1,
		},
	}
	return client
}

func TestClient_GetScopesOnUser(t *testing.T) {
	testCases := []struct {
		name                   string
		impersonatePermissions []ac.Permission
		initTestEnv            func(*Client)
		expectedScopes         []string
	}{
		{
			name: "should return nil when the service account has no impersonate permissions",
			initTestEnv: func(c *Client) {
				c.SelfPermissions = []ac.Permission{}
			},
			expectedScopes: nil,
		},
		{
			name: "should return the 'profile', 'email' and associated RBAC action",
			initTestEnv: func(c *Client) {
				c.SelfPermissions = []ac.Permission{
					{Action: ac.ActionUsersImpersonate, Scope: ac.ScopeUsersAll},
				}
				c.SignedInUser.Permissions = map[int64]map[string][]string{
					1: {
						ac.ActionUsersImpersonate: {ac.ScopeUsersAll},
					},
				}
				c.ImpersonatePermissions = []ac.Permission{
					{Action: ac.ActionUsersRead, Scope: ScopeGlobalUsersSelf},
				}
			},
			expectedScopes: []string{"profile", "email", ac.ActionUsersRead},
		},
		{
			name: "should return 'entitlements' and associated RBAC action scopes",
			initTestEnv: func(c *Client) {
				c.SelfPermissions = []ac.Permission{
					{Action: ac.ActionUsersImpersonate, Scope: ac.ScopeUsersAll},
				}
				c.SignedInUser.Permissions = map[int64]map[string][]string{
					1: {
						ac.ActionUsersImpersonate: {ac.ScopeUsersAll},
					},
				}
				c.ImpersonatePermissions = []ac.Permission{
					{Action: ac.ActionUsersPermissionsRead, Scope: ScopeUsersSelf},
				}
			},
			expectedScopes: []string{"entitlements", ac.ActionUsersPermissionsRead},
		},
		{
			name: "should return 'groups' and associated RBAC action scopes",
			initTestEnv: func(c *Client) {
				c.SelfPermissions = []ac.Permission{
					{Action: ac.ActionUsersImpersonate, Scope: ac.ScopeUsersAll},
				}
				c.SignedInUser.Permissions = map[int64]map[string][]string{
					1: {
						ac.ActionUsersImpersonate: {ac.ScopeUsersAll},
					},
				}
				c.ImpersonatePermissions = []ac.Permission{
					{Action: ac.ActionTeamsRead, Scope: ScopeTeamsSelf},
				}
			},
			expectedScopes: []string{"groups", ac.ActionTeamsRead},
		},
		{
			name: "should return all scopes",
			initTestEnv: func(c *Client) {
				c.SelfPermissions = []ac.Permission{
					{Action: ac.ActionUsersImpersonate, Scope: ac.ScopeUsersAll},
				}
				c.SignedInUser.Permissions = map[int64]map[string][]string{
					1: {
						ac.ActionUsersImpersonate: {ac.ScopeUsersAll},
					},
				}
				c.ImpersonatePermissions = []ac.Permission{
					{Action: ac.ActionUsersRead, Scope: ScopeGlobalUsersSelf},
					{Action: ac.ActionUsersPermissionsRead, Scope: ScopeUsersSelf},
					{Action: ac.ActionTeamsRead, Scope: ScopeTeamsSelf},
					{Action: dashboards.ActionDashboardsRead, Scope: dashboards.ScopeDashboardsAll},
				}
			},
			expectedScopes: []string{"profile", "email", ac.ActionUsersRead,
				"entitlements", ac.ActionUsersPermissionsRead,
				"groups", ac.ActionTeamsRead,
				"dashboards:read"},
		},
		{
			name: "should return stored scopes when the client's impersonate scopes has already been set",
			initTestEnv: func(c *Client) {
				c.SignedInUser.Permissions = map[int64]map[string][]string{
					1: {
						ac.ActionUsersImpersonate: {ac.ScopeUsersAll},
					},
				}
				c.ImpersonateScopes = []string{"dashboard:create", "profile", "email", "entitlements", "groups"}
			},
			expectedScopes: []string{"profile", "email", "entitlements", "groups", "dashboard:create"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := setupTestEnv(t)
			if tc.initTestEnv != nil {
				tc.initTestEnv(c)
			}
			scopes := c.GetScopesOnUser(context.Background(), acimpl.ProvideAccessControl(setting.NewCfg()), 3)
			require.ElementsMatch(t, tc.expectedScopes, scopes)
		})
	}
}

func TestClient_GetScopes(t *testing.T) {
	testCases := []struct {
		name                   string
		impersonatePermissions []ac.Permission
		initTestEnv            func(*Client)
		expectedScopes         []string
	}{
		{
			name: "should return default scopes when the signed in user is nil",
			initTestEnv: func(c *Client) {
				c.SignedInUser = nil
			},
			expectedScopes: []string{"profile", "email", "entitlements", "groups"},
		},
		{
			name: "should return additional scopes from signed in user's permissions",
			initTestEnv: func(c *Client) {
				c.SignedInUser.Permissions = map[int64]map[string][]string{
					1: {
						dashboards.ActionDashboardsRead: {dashboards.ScopeDashboardsAll},
					},
				}
			},
			expectedScopes: []string{"profile", "email", "entitlements", "groups", "dashboards:read"},
		},
		{
			name: "should return default scopes when the signed in user has no permissions",
			initTestEnv: func(c *Client) {
				c.SignedInUser.Permissions = map[int64]map[string][]string{}
			},
			expectedScopes: []string{"profile", "email", "entitlements", "groups"},
		},
		{
			name: "should return stored scopes when the client's scopes has already been set",
			initTestEnv: func(c *Client) {
				c.Scopes = []string{"profile", "email", "entitlements", "groups"}
			},
			expectedScopes: []string{"profile", "email", "entitlements", "groups"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := setupTestEnv(t)
			if tc.initTestEnv != nil {
				tc.initTestEnv(c)
			}
			scopes := c.GetScopes()
			require.ElementsMatch(t, tc.expectedScopes, scopes)
		})
	}
}

func TestClient_ToDTO(t *testing.T) {
	client := &Client{
		ID:          1,
		Name:        "my-ext-service",
		ClientID:    "test",
		Secret:      "testsecret",
		RedirectURI: "http://localhost:3000",
		GrantTypes:  "client_credentials,urn:ietf:params:oauth:grant-type:jwt-bearer",
		Audiences:   "https://example.org,https://second.example.org",
		PublicPem:   []byte("pem_encoded_public_key"),
	}

	dto := client.ToDTO()

	require.Equal(t, client.ClientID, dto.ID)
	require.Equal(t, client.Name, dto.Name)
	require.Equal(t, client.RedirectURI, dto.RedirectURI)
	require.Equal(t, client.GrantTypes, dto.GrantTypes)
	require.Equal(t, client.Audiences, dto.Audiences)
	require.Equal(t, client.PublicPem, []byte(dto.KeyResult.PublicPem))
	require.Empty(t, dto.KeyResult.PrivatePem)
	require.Empty(t, dto.KeyResult.URL)
	require.False(t, dto.KeyResult.Generated)
	require.Equal(t, client.Secret, dto.Secret)
}

func Test_ParsePublicKeyPem(t *testing.T) {
	testCases := []struct {
		name         string
		publicKeyPem string
		wantErr      bool
	}{
		{
			name:         "should return error when the public key pem is empty",
			publicKeyPem: "",
			wantErr:      true,
		},
		{
			name: "should return error when the public key pem is invalid",
			publicKeyPem: `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAxP72NEnQF3o3eFFMtFqyloW9oLhTydxXS2dA2NolMvXewO77
UvJw54wkOdrJrJO2BIw+XBrrb+13+koRUnwa2DNsh+SWG0PEe/31mt0zJrCmNM37
iJYIu3KZR2aRlierVY5gyrIniBIZ9blQspI6SRY9xmo6Wdh0VCRnsCV5sMlaqerI
snLpYOjGtMmL0rFuW2jKrAzpbq7L99IDgPbiH7tluaQkGIxoc29S4wjwg0NgQONT
GkfJEeXQIkxOHNM5WGb8mvjX4U0jMdXvC4WUWcS+KpcIV7ee8uEs2xDz++N4HYAS
ty37sY8wwW22QUW9I7XlSC4rsC88Ft5ar8yLsQIDAQABAoIBAAQ1yTv+mFmKGYGj
JiskFZVBNDdpPRQvNvfj8+c2iU08ozc3HEyuZQKT1InefsknCoCwIRyNkDrPBc2F
8/cR8y5r8e25EUqxoPM/7xXxVIinBZRTEyU9BKCB71vu6Z1eiWs9jNzEIDNopKCj
ZmG8nY2Gkckp58eYCEtskEE72c0RBPg8ZTBdc1cLqbNVUjkLvR5e98ruDz6b+wyH
FnztZ0k48zM047Ior69OwFRBg+S7d6cgMMmcq4X2pg3xgQMs0Se/4+pmvBf9JPSB
kl3qpVAkzM1IFdrmpFtBzeaqYNj3uU6Bm7NxEiqjAoeDxO231ziSdzIPtXIy5eRl
9WMZCqkCgYEA1ZOaT77aa54zgjAwjNB2Poo3yoUtYJz+yNCR0CPM4MzCas3PR4XI
WUXo+RNofWvRJF88aAVX7+J0UTnRr25rN12NDbo3aZhX2YNDGBe3hgB/FOAI5UAh
9SaU070PFeGzqlu/xWdx5GFk/kiNUNLX/X4xgUGPTiwY4LQeI9lffzkCgYEA7CA7
VHaNPGVsaNKMJVjrZeYOxNBsrH99IEgaP76DC+EVR2JYVzrNxmN6ZlRxD4CRTcyd
oquTFoFFw26gJIJAYF8MtusOD3PArnpdCRSoELezYdtVhS0yx8TSHGVC9MWSSt7O
IdjzEFpA99HPkYFjXUiWXjfCTK7ofI0RXC6a+DkCgYEAoQb8nYuEGwfYRhwXPtQd
kuGbVvI6WFGGN9opVgjn+8Xl/6jU01QmzkhLcyAS9B1KPmYfoT4GIzNWB7fURLS3
2bKLGwJ/rPnTooe5Gn0nPb06E38mtdI4yCEirNIqgZD+aT9rw2ZPFKXqA16oTXvq
pZFzucS4S3Qr/Z9P6i+GNOECgYBkvPuS9WEcO0kdD3arGFyVhKkYXrN+hIWlmB1a
xLS0BLtHUTXPQU85LIez0KLLslZLkthN5lVCbLSOxEueR9OfSe3qvC2ref7icWHv
1dg+CaGGRkUeJEJd6CKb6re+Jexb9OKMnjpU56yADgs4ULNLwQQl/jPu81BMkwKt
CVUkQQKBgFvbuUmYtP3aqV/Kt036Q6aB6Xwg29u2XFTe4BgW7c55teebtVmGA/zc
GMwRsF4rWCcScmHGcSKlW9L6S6OxmkYjDDRhimKyHgoiQ9tawWag2XCeOlyJ+hkc
/qwwKxScuFIi2xwT+aAmR70Xk11qXTft+DaEcHdxOOZD8gA0Gxr3
-----END RSA PRIVATE KEY-----`,
			wantErr: true,
		},
		{
			name: "should parse the public key if it is in PKCS1 format",
			publicKeyPem: `-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAy06MeS06Ea7zGKfOM8kosxuUBMNhrWKWMvW4Jq1IXG+lyTfann2+
kI1rKeWAQ9YbxNzLynahoKN47EQ6mqM1Yj5v9iKWtSvCMKHWBuqrG5ksaEQaAVsA
PDg8aOQrI1MSW9Hoc1CummcWX+HKNPVwIzG3sCboENFzEG8GrJgoNHZgmyOYEMMD
2WCdfY0I9Dm0/uuNMAcyMuVhRhOtT3j91zCXvDju2+M2EejApMkV9r7FqGmNH5Hw
8u43nWXnWc4UYXEItE8nPxuqsZia2mdB5MSIdKu8a7ytFcQ+tiK6vempnxHZytEL
6NDX8DLydHbEsLUn6hc76ODVkr/wRiuYdQIDAQAB
-----END RSA PUBLIC KEY-----`,
			wantErr: false,
		},
		{
			name: "should parse the public key if it is in PKCS8 format",
			publicKeyPem: `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEbsGtoGJTopAIbhqy49/vyCJuDot+
mgGaC8vUIigFQVsVB+v/HZ4yG1Rcvysig+tyNk1dZQpozpFc2dGmzHlGhw==
-----END PUBLIC KEY-----`,
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := utils.ParsePublicKeyPem([]byte(tc.publicKeyPem))
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
