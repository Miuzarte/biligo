package biligo

// ExportIdentity 返回当前的身份信息
func ExportIdentity() Identity {
	return identity
}

// StoreIdentity 存储身份信息，包括 Cookie、RefreshToken 和 Uid
func StoreIdentity(id Identity) error {
	identity.Cookie = id.Cookie
	identity.RefreshToken = id.RefreshToken
	identity.Uid = id.Uid
	return cookie.store(identity.Cookie)
}

// ExportCookie 返回当前存储的 Cookie
func ExportCookie() string {
	return identity.Cookie
}

// StoreCookie 存储新的 Cookie
func StoreCookie(cookieStr string) error {
	identity.Cookie = cookieStr
	return cookie.store(identity.Cookie)
}

// ExportRefreshToken 返回当前存储的 RefreshToken
func ExportRefreshToken() string {
	return identity.RefreshToken
}

// StoreRefreshToken 存储新的 RefreshToken
func StoreRefreshToken(refreshToken string) {
	identity.RefreshToken = refreshToken
}
