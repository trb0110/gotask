package gotask

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"sync"
	"time"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"errors"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const SIGNING_KEY = "6818ae239b509ab943a267f896d014e48c673d8d990b045c50f2eac017c8611b"

type Token_Sync struct {
	sync.Mutex
	Tokens map[string]Token
}

type Token struct  {
	username		string
	userid			string
	Role			string
	tokenCreation 	time.Time
	tokenLastCalled time.Time
}


var TokensMap Token_Sync


func InitializeTokens(){
	TokensMap = Token_Sync{
		Tokens: make(map [string] Token),
	}
	//fmt.Println(TokensMap)
}


func (post *Token_Sync) Get(key string) (Token, bool) {
	post.Lock()
	defer post.Unlock()
	tok,exists := post.Tokens[key]
	return tok,exists
}

func (post *Token_Sync)DeleteAll() {
	post.Lock()
	defer post.Unlock()
	for key, _ := range post.Tokens {
		delete(post.Tokens, key)
	}
}

func (post *Token_Sync)Set(username string, tokenKey string,userid string,role string) {
	post.Lock()
	defer post.Unlock()
	t:= time.Now()
	fmt.Println("Time of token creation" + t.String())
	token := Token{username,userid,role,t,t}
	post.Tokens[tokenKey] = token
}

func RandString(length int) string {
	return StringWithCharset(length, charset)
}
func StringWithCharset(length int, charset string) string {

	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
func GenerateToken(user *User)(tokenK string,err error){

	randString := RandString(10)
	//fmt.Println(randString)
	hashedPassword , _ := bcrypt.GenerateFromPassword([]byte(SIGNING_KEY),8)
	hashedUser , _ := bcrypt.GenerateFromPassword([]byte(user.Username+randString),8)

	fullHesh := string(hashedPassword) +string(hashedUser)

	hashedFullHesh , _ :=  bcrypt.GenerateFromPassword([]byte(fullHesh),8)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true

	tokenString, err := token.SignedString(hashedFullHesh)
	if err != nil {
		return "", errors.New("Failed to sign token")
	}

	TokensMap.Set(user.Username, string(tokenString),user.UserID,user.Role)
	//fmt.Println(TokensMap)
	return string(tokenString), nil
}

func LookupTokenKey(tokenKey string)(exists bool){

	//Check if provided security token exists in map
	_, exists = TokensMap.Get(tokenKey)

	return exists
}
