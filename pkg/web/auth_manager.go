package web

/*
	The Authentication Manager informs the Session Manager that the session associated
	with the token is to be connected to a specific desktop client. The Session Manager
	then informs each service in the session that it must connect directly to the client.
	The user can then interact with the session.
*/

type AuthenticationManager interface {
}
