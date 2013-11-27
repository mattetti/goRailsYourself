// The crypto package ports some of Rails' crypto functions that
// can be useful outside of Rails. (mainly used to share signed/encryoted messages/cookies/sessions).
package crypto

type MsgSerializer interface {
	Serialize(v interface{}) (string, error)
	Unserialize(data string, v interface{}) error
}
