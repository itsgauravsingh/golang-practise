package repos

func UserIsValid(uname, pwd string) bool {
	_uname, _pwd, _isValid := "admin", "admin123", false
	if _uname == uname && _pwd == pwd {
		_isValid = true
	} else {
		_isValid = false
	}
	return _isValid
}
