package main

import (
	"sync/atomic"
	"time"
)

type User struct {
	id int64
	name string
}
func NewUser(id int64, name string) *User{
	return &User{id,name}
}
func (u *User)GetID() int64 {
	return u.id
}
func (u *User)GetName() string {
	return u.name
}

type Article struct {
	id     int64
	title string
	link string
	posterID int64
	timeStampMS int64
	votes int64
}
func NewArticle(id, posterID int64, title, link string) *Article{
	ts := int64(time.Now().Nanosecond()/1E6)
	 return &Article{id ,title, link, posterID,
	 	ts, 0}
}
func (a *Article) GetVote() int64{
	return a.votes
}
func (a *Article) AddVote() {
	atomic.AddInt64(&(a.votes), int64(1))
}
func (a *Article) GetTime() int64{
	return a.timeStampMS
}