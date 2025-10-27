package biligo

// ExportIdentity 导出当前的身份信息
func ExportIdentity() Identity {
	return identity
}

// ImportIdentity 导入身份信息，包括 Cookie、RefreshToken 和 Uid
func ImportIdentity(id Identity) error {
	identity.Cookie = id.Cookie
	identity.RefreshToken = id.RefreshToken
	identity.Uid = id.Uid
	return cookie.set(identity.Cookie)
}

// ExportCookie 导出当前的 Cookie
func ExportCookie() string {
	return identity.Cookie
}

// ImporCookie 导入新的 Cookie
func ImporCookie(cookieStr string) error {
	identity.Cookie = cookieStr
	return cookie.set(identity.Cookie)
}

// ExportRefreshToken 导出当前的 RefreshToken
func ExportRefreshToken() string {
	return identity.RefreshToken
}

// ImporRefreshToken 导入新的 RefreshToken
func ImporRefreshToken(refreshToken string) {
	identity.RefreshToken = refreshToken
}
