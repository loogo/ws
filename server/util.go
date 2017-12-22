package server

func contains(slice []*User, item *User) bool {
	for _, val := range slice {
		if val != nil && val.Code == item.Code {
			return true
		}
	}
	return false
}
