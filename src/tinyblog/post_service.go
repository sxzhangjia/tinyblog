package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	//"time"
	//"log"
)

type PostService struct {
	DBC *MDBC
	C   *mgo.Collection
	S   *Session
}

func (self *PostService) Save(p *Post) *ResMessage {

	t := &TimeData{}

	if p.Title == "" || p.Content == "" {
		return getUserResMessage(false, REQUIRED_DEFAULT, POST_MODE_CODE)
	}

	p.Id_ = bson.NewObjectId()
	p.CreateTime = t
	p.EditTime = t

	//log.Println(t)

	err := self.C.Insert(p)

	return getResMessage(err, SAVE_SUCCESS, POST_MODE_CODE)
}

func (self *PostService) Del(id string) *ResMessage {
	err := self.C.RemoveId(bson.ObjectIdHex(id))
	return getResMessage(err, DEL_SUCCESS, POST_MODE_CODE)
}

func (self *PostService) Discard(id string) *ResMessage {
	err := self.C.UpdateId(bson.ObjectIdHex(id), bson.M{"isdiscard": true})
	return getResMessage(err, UPDATE_SUCCESS, POST_MODE_CODE)
}

func (self *PostService) GetOne(sel *SelectData) *Post {

	p := new(Post)
	self.C.Find(sel.Condition).One(p)
	return p
}

func (self *PostService) GetOneById(id string) (*Post, bool) {

	p := new(Post)
	err := self.C.FindId(bson.ObjectIdHex(id)).One(p)
	if err == nil {
		return p, true
	} else {
		return p, false
	}
}

func (self *PostService) IsExist(title string) bool {
	p := new(Post)
	self.C.Find(bson.M{"title": title}).Select(bson.M{"title": 1}).One(p)
	if p.Title != "" {
		return true
	} else {
		return false
	}
}

func (self *PostService) Count(sel *SelectData) int {

	if sel.Condition == nil {
		n, _ := self.C.Count()
		return n
	} else {
		n, _ := self.C.Find(sel.Condition).Count()
		return n
	}
}

func (self *PostService) GetList(sel *SelectData) []Post {

	pl := make([]Post, sel.Limit)
	q := self.C.Find(sel.Condition)
	q = q.Sort(sel.Sort).Skip(sel.Start).Limit(sel.Limit)
	q.All(&pl)

	//log.Println("----------------------------", pl)

	return pl
}

/*


func (self *PostService) Update(postId int, data interface{}) {
	self.DBC.UpdateSet(POST_TAB, BSONM{"id": postId}, data)
}

func (self *PostService) InsertComment(postId int, content string) {

	post := &Post{}
	self.DBC.SelectOne(POST_TAB, BSONM{"id": postId}, post)

	comment := &Comment{
		Content:    content,
		Auth:       "haha",
		Email:      "e@qq.com",
		host:       "",
		Ip:         "",
		Display:    true,
		CreateTime: time.Now(),
	}

	self.DBC.UpdatePush(POST_TAB, BSONM{"id": postId}, "comment", comment)

}

func (self *PostService) deleteComment(postId int, commentId int) {

	postSel := BSONM{"id": postId}
	commentSel := BSONM{"id": commentId}
	self.DBC.UpdatePull(POST_TAB, postSel, "comment", commentSel)
	self.DBC.UpdateInc(POST_TAB, BSONM{"id": postId}, "commentnum", -1)
}
*/
//func (self *PostService) updateComment( postId int, commentId int, BSONM sel) {

//db.shcool.update({ "_id" : 2, "students.name" : "ajax"},{"$inc" : {"students.0.age" : 1} });
//postSel := BSONM{"id": postId}
//commentSel := BSONM{"id": commentId}
//self.DBC.UpdateSet(POST_TAB, postSel, data)

//}
