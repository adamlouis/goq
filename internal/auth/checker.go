package auth

type UPChecker interface {
	Check(username, password string) bool
}

func NewConstUPChecker(username, password string) UPChecker {
	return &rc{
		username: username,
		password: password,
	}
}

type rc struct {
	username, password string
}

func (r *rc) Check(username, password string) bool {
	return username == r.username && password == r.password
}

type KChecker interface {
	Check(k string) bool
}

func NewConstKChecker(k string) KChecker {
	return &kc{k: k}
}

type kc struct {
	k string
}

func (s *kc) Check(k string) bool {
	return k == s.k
}
