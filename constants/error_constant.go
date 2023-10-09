package constants

const ERR_100001 = 100001
const ERR_100002 = 100002
const ERR_100003 = 100003
const ERR_400 = 400
func GetErrorConstant() map[int]string {
	return map[int]string{
		ERR_100001: "Internal Query Error",
		ERR_100002: "Panic Nil Value",
		ERR_100003: "Overflow",
		ERR_400: "Bad Request",
		0:"Unknown error",
	}
}
