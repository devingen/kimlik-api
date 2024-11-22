package dto

// GetUserInfoResponse reflects the OpenID UserInfo Response data.
// See: https://openid.net/specs/openid-connect-core-1_0.html#UserInfoResponse
type GetUserInfoResponse struct {
	Sub               string `json:"sub"`
	Name              string `json:"name"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	Picture           string `json:"picture"`
}
