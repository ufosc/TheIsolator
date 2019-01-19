package main

import (
	"regexp"
)

//Regexes for request path matching

//=================================

// Actor

var applicationRegexActor = regexp.MustCompile("^/activity/application/?$")
var personRegexActor = regexp.MustCompile("^/activity/person/[^/]+)/?$")
var groupRegexActor = regexp.MustCompile("^/activity/group/([^/]+)/?$")

// Inbox
var applicationRegexInbox = regexp.MustCompile("^/activity/application/inbox/?$")
var personRegexInbox = regexp.MustCompile("^/activity/person/([^/]+)/inbox/?$")
var groupRegexInbox = regexp.MustCompile("^/activity/group/([^]+)/inbox/?$")

// Outbox

var applicationRegexOutbox = regexp.MustCompile("^/activity/application/outbox/?$")
var personRegexOutbox = regexp.MustCompile("^/activity/person/([^/]+)/outbox/?$")
var groupRegexOutbox = regexp.MustCompile("^/activity/group/([^/]+)/outbox/?$")

// Following
var applicationRegexFollowing = regexp.MustCompile("^/activity/application/following/?$")
var personRegexFollowing = regexp.MustCompile("^/activity/person/([^/]+)/following/?$")
var groupRegexFollowing = regexp.MustCompile("^/activity/group/([^/]+)/following/?$")

// Followers
var applicationRegexFollowers = regexp.MustCompile("^/activity/application/followers/?$")
var personRegexFollowers = regexp.MustCompile("^/activity/person/([^/]+)/followers/?$")
var groupRegexFollowers = regexp.MustCompile("^/activity/group/([^/]+)/followers/?$")

// Liked
var applicationRegexLiked = regexp.MustCompile("^/activity/application/liked/?$")
var personRegexLiked = regexp.MustCompile("^/activity/person/([^/]+)/liked/?$")
var groupRegexLiked = regexp.MustCompile("^/activity/group/([^/]+)/liked/?$")

var documentRegex = regexp.MustCompile("^/document/?$")
var createRegex = regexp.MustCompile("^/create/?$")
var likedRegex = regexp.MustCompile("^/liked/?$")
var joinRegex = regexp.MustCompile("^/join/?$")
var followRegex = regexp.MustCompile("^/follow/?$")
