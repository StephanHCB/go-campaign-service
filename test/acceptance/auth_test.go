package acceptance

func tstUnauthenticated() string {
	return ""
}

const validtoken_demosecret_HS256_admin = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJodHRwczovL2dpdGh1Yi5jb20vU3RlcGhhbkhDQi9nby1jYW1wYWlnbi1zZXJ2aWNlL3JvbGVzIjpbImFkbWluIl19.EZC_nxHsZKrNLK6BvFqJrgpqWMv8OnnjpxAwst3b9RA"

const validtoken_demosecret_HS256_noroles = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.6F1Mu1zkAGqlk65ndU8InIVa5N8LIhDuOQYr-V_x8Tk"

func tstAuthAdmin() string {
	return validtoken_demosecret_HS256_admin
}

func tstAuthUser() string {
	return validtoken_demosecret_HS256_noroles
}
