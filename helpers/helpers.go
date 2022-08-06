package helpers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/schema"
)

func NowPlusSeconds(seconds int) time.Time {
	return time.Now().Local().Add(time.Second * time.Duration(seconds))
}

func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func ParseDateIso(date string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05-0700", date)
}

func FormatDateIso(date time.Time) string {
	return date.Format("2006-01-02T15:04:05-0700")
}

func RespondWithError(w http.ResponseWriter, code int, errorCode string, message string) {
	fmt.Printf("ErrorCode: %v, Message: %v", errorCode, message)
	RespondWithJSON(w, code, map[string]interface{}{"success": false, "errorCode": errorCode, "message": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

var Decoder = schema.NewDecoder()
