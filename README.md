# Passbase sample app

Very simple application which has a protected endpoint at `/verified`. This endpoint only returns a 200 if the authenticated user has also had their identity verified through Passbase.

Authentication to the web app is achieved via Basic auth for simplicity. In reality, this could be anything (OIDC, LDAP, etc).

The only necessary thing is to map a user id (in the applications domain) to a Passbase identity verification key (from passbase's domain).
