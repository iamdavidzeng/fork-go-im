package helper

import (
	"encoding/json"
	"log"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

func JsonToMap(str []byte) map[string]interface{} {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(str, &jsonMap)
	if err != nil {
		panic(err)
	}
	return jsonMap
}

func HttpBuildQuery(queryData url.Values) string {
	return queryData.Encode()
}

func HashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	plainPwds := []byte(plainPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwds)
	if err != nil {
		return false
	}
	return true
}

func ProduceChannelName(fId, tId string) (cA, cB string) {
	cA = "channel_" + fId + "_" + tId
	cB = "channel_" + tId + "_" + fId
	return cA, cB
}

func ProduceChannelGroupName(tId string) string {
	return "channel_" + tId
}
