package govel

import (
	"net/http"
)

// Expire causes the session to expire.
func (s *Session) Expire() {
	s.MaxAge(-1)
}

// Get returns the value for the given key from the session.
func (s *Session) Get(key string) interface{} {
	return s.session.Values[key]
}

// Set sets the value for the given key in the session.
func (s *Session) Set(key string, value interface{}) {
	s.session.Values[key] = value
}

// Delete removes the value for the given key from the session.
func (s *Session) Delete(key string) {
	delete(s.session.Values, key)
}

// IsNew returns true if the session is new.
func (s *Session) IsNew() bool {
	return s.session.IsNew
}

// SetFlash sets the value for the given key in the session with the "_flash" suffix.
func (s *Session) SetFlash(key string, value interface{}) {
	s.Set(key+"_flash", value)
}

// GetFlash returns the flash value for the given key and deletes it.
func (s *Session) GetFlash(key string) interface{} {
	value := s.Get(key + "_flash")

	if value != nil {
		s.Delete(key + "_flash")
		return value
	}

	return nil
}

// Domain sets the domain for the session.
func (s *Session) Domain(domain string) {
	s.session.Options.Domain = domain
}

// Path sets the path for the session.
func (s *Session) Path(path string) {
	s.session.Options.Path = path
}

// MaxAge sets the max age for the session.
func (s *Session) MaxAge(age int) {
	s.session.Options.MaxAge = age
}

// Secure sets the secure flag for the session.
func (s *Session) Secure(secure bool) {
	s.session.Options.Secure = secure
	s.session.Options.HttpOnly = true
}

// SameSite sets the same site flag for the session.
func (s *Session) SameSite(sameSite http.SameSite) {
	s.session.Options.SameSite = sameSite
}
