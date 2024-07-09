package auth

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// https://datatracker.ietf.org/doc/html/draft-ietf-oauth-security-topics#refresh_token_protection

const (
	ScopeRefresh = "refresh"
	ScopeAccess  = "access"
)

type ScopeStrings []string

var _ json.Marshaler = (ScopeStrings)(nil)
var _ json.Unmarshaler = (*ScopeStrings)(nil)

func (s ScopeStrings) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.Join(s, " "))
}
func (s *ScopeStrings) UnmarshalJSON(b []byte) error {
	str := ""
	if len(b) == 0 {
		*s = nil
		return nil
	}
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	*s = strings.Split(str, " ")
	return nil
}

type Claims struct {
	jwt.RegisteredClaims

	// The `scope` (Scope) claim. See https://www.rfc-editor.org/rfc/rfc8693.html#name-scope-scopes-claim
	Scope ScopeStrings `json:"scope,omitempty"`
}

func NewClaims() *Claims {
	t := jwt.NewNumericDate(time.Now())
	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  t,
			NotBefore: t,
		},
	}
}

func (c *Claims) WithLifetime(duration time.Duration) *Claims {
	return c.WithExpirationTime(time.Now().Add(duration))
}

// https://www.iana.org/assignments/jwt/jwt.xhtml

// The "iss" (issuer) claim identifies the principal that issued the
// JWT.  The processing of this claim is generally application specific.
// The "iss" value is a case-sensitive string containing a StringOrURI
// value.  Use of this claim is OPTIONAL.
func (c *Claims) WithIssuer(iss string) *Claims {
	c.Issuer = iss
	return c
}

// The "sub" (subject) claim identifies the principal that is the
// subject of the JWT.  The claims in a JWT are normally statements
// about the subject.  The subject value MUST either be scoped to be
// locally unique in the context of the issuer or be globally unique.
// The processing of this claim is generally application specific.  The
// "sub" value is a case-sensitive string containing a StringOrURI
// value.  Use of this claim is OPTIONAL.
func (c *Claims) WithSubject(sub string) *Claims {
	c.Subject = sub
	return c
}

// The "aud" (audience) claim identifies the recipients that the JWT is
// intended for.  Each principal intended to process the JWT MUST
// identify itself with a value in the audience claim.  If the principal
// processing the claim does not identify itself with a value in the
// "aud" claim when this claim is present, then the JWT MUST be
// rejected.  In the general case, the "aud" value is an array of case-
// sensitive strings, each containing a StringOrURI value.  In the
// special case when the JWT has one audience, the "aud" value MAY be a
// single case-sensitive string containing a StringOrURI value.  The
// interpretation of audience values is generally application specific.
// Use of this claim is OPTIONAL.
func (c *Claims) WithAudience(aud []string) *Claims {
	c.Audience = aud
	return c
}

// The "exp" (expiration time) claim identifies the expiration time on
// or after which the JWT MUST NOT be accepted for processing.  The
// processing of the "exp" claim requires that the current date/time
// MUST be before the expiration date/time listed in the "exp" claim.
// Implementers MAY provide for some small leeway, usually no more than
// a few minutes, to account for clock skew.  Its value MUST be a number
// containing a NumericDate value.  Use of this claim is OPTIONAL.
func (c *Claims) WithExpirationTime(exp time.Time) *Claims {
	c.ExpiresAt = jwt.NewNumericDate(exp)
	return c
}

// The "nbf" (not before) claim identifies the time before which the JWT
// MUST NOT be accepted for processing.  The processing of the "nbf"
// claim requires that the current date/time MUST be after or equal to
// the not-before date/time listed in the "nbf" claim.  Implementers MAY
// provide for some small leeway, usually no more than a few minutes, to
// account for clock skew.  Its value MUST be a number containing a
// NumericDate value.  Use of this claim is OPTIONAL.
func (c *Claims) WithNotBeforeTime(nbf time.Time) *Claims {
	c.NotBefore = jwt.NewNumericDate(nbf)
	return c
}

// The "iat" (issued at) claim identifies the time at which the JWT was
// issued.  This claim can be used to determine the age of the JWT.  Its
// value MUST be a number containing a NumericDate value.  Use of this
// claim is OPTIONAL.
func (c *Claims) WithIssuedAtTime(iat time.Time) *Claims {
	c.IssuedAt = jwt.NewNumericDate(iat)
	return c
}

// The "jti" (JWT ID) claim provides a unique identifier for the JWT.
// The identifier value MUST be assigned in a manner that ensures that
// there is a negligible probability that the same value will be
// accidentally assigned to a different data object; if the application
// uses multiple issuers, collisions MUST be prevented among values
// produced by different issuers as well.  The "jti" claim can be used
// to prevent the JWT from being replayed.  The "jti" value is a case-
// sensitive string.  Use of this claim is OPTIONAL.
func (c *Claims) WithJWTID(jti string) *Claims {
	c.ID = jti
	return c
}

// The value of the scope claim is a JSON string containing a space-separated
// list of scopes associated with the token, in the format described in Section
// 3.3 of [RFC6749].
//
// https://www.rfc-editor.org/rfc/rfc8693.html#name-scope-scopes-claim
func (c *Claims) WithScopes(scope ...string) *Claims {
	c.Scope = scope
	return c
}
