package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"math/rand"
	"strconv"
	"time"
)

type SessionData struct {
	U  *User
	ch chan bool
}

func CreateUUID() string {

	nano := time.Now().UnixNano()
	rand.Seed(nano)
	randNum := rand.Int63()

	m := md5.New()
	io.WriteString(m, strconv.FormatInt(randNum, 16))

	return hex.EncodeToString(m.Sum(nil))
}

type Session struct {
	ExpireHour time.Duration
	Data       map[string]*SessionData
}

func (self *Session) del(id string) {
	time.Sleep(self.ExpireHour * time.Hour)
	delete(self.Data, id)
}

func (self *Session) New(sd *SessionData) string {
	id := CreateUUID()
	self.Data[id] = sd
	go self.del(id)
	return id
}

func (self *Session) Destroy(k string) {

}

func (self *Session) Set(k string, v *SessionData) {
	self.Data[k] = v
}

func (self *Session) Get(k string) *SessionData {
	return self.Data[k]
}

func (self *Session) IsLogin(uuid string) bool {
	u := self.Get(uuid)
	if u != nil {
		return true
	} else {
		return false
	}
}

func (self *Session) GetPowerCode(uuid string) int {
	usrData := self.Data[uuid]
	if usrData == nil {
		return -1
	} else {
		return usrData.U.PowerCode
	}
}

func (self *Session) Auth(uuid string) {

}
