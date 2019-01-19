package main

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/vocab"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/httpsig"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

/* ====APPLICATION====




*/ 

type Application struct() 

func (m *Application) Owns(c context.Context, id *url.URL) (owns bool) { 
	log.Printf("BEGIN Checking ownership: %s" id)
	defer log.Printf("END Checking ownership: %s" id)
	//The server owns the provided ID if the hostname and port
	//of the given ID matches those provided by the user in the configuration file. 

	h1 := id.Host
	h2 := viper.GetString("server.hostname") + ":" + viper.GetString("server.port")
	owns = h1 == h2 || id.Hostname() == "" 

	log.Println("Onwership status:", owns) 
	return
}

func (m *Application) Get(c context.Context, id *url.URL, rw pub.RWType) (pub.PubObject, error) { 
	log.Printf("BEGIN Getting: %s", id) 
	defer log.Printf("END getting %s", id)

	//First check if the server has the given ID 
	if has, err := m.Has(c, id); err != nil {
		return nil, fmt.Errorf("server does not have the given ID:%s", err)
	}
	else if !has {
		return nil, fmt.Errorf("%s not found", id)
	}

	// Grab the request path from the ID

	p := id.Path

	

	switch {
	// Get Actor
	// =========
	case personRegexActor.MatchString(p):
		n := personRegexActor.FindStringSubmatch(p)[1]
		a, ok := personActors[n]
		if !ok {
			return nil, fmt.Errorf("GET person (actor) not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor, nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor, nil
		default:
			panic("Impossible")
		}
	case groupRegexActor.MatchString(p):
		n := groupRegexActor.FindStringSubmatch(p)[1]
		a, ok := groupActors[n]
		if !ok {
			return nil, fmt.Errorf("GET group (actor) not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor, nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor, nil
		default:
			panic("Impossible")
		}

	// Get Inbox
	// =========
	case personRegexInbox.MatchString(p):
		n := personRegexInbox.FindStringSubmatch(p)[1]
		a, ok := personActors[n]
		if !ok {
			return nil, fmt.Errorf("GET person (inbox) not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor.GetInboxOrderedCollection(), nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor.GetInboxOrderedCollection(), nil
		default:
			panic("Impossible")
		}
	case groupRegexInbox.MatchString(p):
		n := groupRegexInbox.FindStringSubmatch(p)[1]
		a, ok := groupActors[n]
		if !ok {
			return nil, fmt.Errorf("GET group (inbox) not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor.GetInboxOrderedCollection(), nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor.GetInboxOrderedCollection(), nil
		default:
			panic("Impossible")
		}

	// Get Outbox
	// ==========
	case personRegexOutbox.MatchString(p):
		n := personRegexOutbox.FindStringSubmatch(p)[1]
		a, ok := personActors[n]
		if !ok {
			return nil, fmt.Errorf("GET person (outbox) not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor.GetOutboxOrderedCollection(), nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor.GetOutboxOrderedCollection(), nil
		default:
			panic("Impossible")
		}
	case groupRegexOutbox.MatchString(p):
		n := groupRegexOutbox.FindStringSubmatch(p)[1]
		a, ok := groupActors[n]
		if !ok {
			return nil, fmt.Errorf("GET group (outbox) not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor.GetOutboxOrderedCollection(), nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor.GetOutboxOrderedCollection(), nil
		default:
			panic("Impossible")
		}

	// Get Following
	// =============
	case personRegexFollowing.MatchString(p):
		n := personRegexFollowing.FindStringSubmatch(p)[1]
		a, ok := personActors[n]
		if !ok {
			return nil, fmt.Errorf("GET person (following) not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor.GetFollowingCollection(), nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor.GetFollowingCollection(), nil
		default:
			panic("Impossible")
		}
	case groupRegexFollowing.MatchString(p):
		n := groupRegexFollowing.FindStringSubmatch(p)[1]
		a, ok := groupActors[n]
		if !ok {
			return nil, fmt.Errorf("GET group (following) not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor.GetFollowingCollection(), nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor.GetFollowingCollection(), nil
		default:
			panic("Impossible")
		}

	// Get Followers
	// =============
	case personRegexFollowers.MatchString(p):
		n := personRegexFollowers.FindStringSubmatch(p)[1]
		a, ok := personActors[n]
		if !ok {
			return nil, fmt.Errorf("GET person (followers) not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor.GetFollowersCollection(), nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor.GetFollowersCollection(), nil
		default:
			panic("Impossible")
		}
	case groupRegexFollowers.MatchString(p):
		n := groupRegexFollowers.FindStringSubmatch(p)[1]
		a, ok := groupActors[n]
		if !ok {
			return nil, fmt.Errorf("GET group (followers) not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor.GetFollowersCollection(), nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor.GetFollowersCollection(), nil
		default:
			panic("Impossible")
		}

	// Get Liked
	// =========
	case personRegexLiked.MatchString(p):
		n := personRegexLiked.FindStringSubmatch(p)[1]
		a, ok := personActors[n]
		if !ok {
			return nil, fmt.Errorf("GET person (liked) not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor.GetLikedCollection(), nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor.GetLikedCollection(), nil
		default:
			panic("Impossible")
		}
	case groupRegexLiked.MatchString(p):
		n := groupRegexLiked.FindStringSubmatch(p)[1]
		a, ok := groupActors[n]
		if !ok {
			return nil, fmt.Errorf("GET group (liked) not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor.GetLikedCollection(), nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor.GetLikedCollection(), nil
		default:
			panic("Impossible")
		}

	// Get Default
	// ===========
	default:
		return nil, fmt.Errorf("GET not found")
	}
}


func (m *Application) GetAsVerifiedUser (c context.Context, id, authdUser *url.URL, rw pub.RWType) (pub.PubObject, error){
	log.Printf("BEGIN Getting verified: %s", id)
	defer log.Printf("END Getting verified: %s", id)

	// Since all the ActivityStream objects are public in a social link
	// Aggregator, we can simply call Get in this case

	return m.Get(c, id, rw)
}

// Has determines if the server already knows about the object
// or activity specified by the given ID
func (m *Application) Has(c context.Context, id *url.URL) (bool, error) { 
	log.Printf("BEGIN Checking has: %s", ID)
	defer log.Printf("END Checking has: %s", id)

	// Grab the request path from the ID
	p := id.Path

		switch {
	// Has Actor
	// =========
	case personRegexActor.MatchString(p):
		n := personRegexActor.FindStringSubmatch(p)[1]
		_, ok := personActors[n]
		return ok, nil
	case groupRegexActor.MatchString(p):
		n := groupRegexActor.FindStringSubmatch(p)[1]
		_, ok := groupActors[n]
		return ok, nil

	// Has Inbox
	// =========
	case personRegexInbox.MatchString(p):
		n := personRegexInbox.FindStringSubmatch(p)[1]
		_, ok := personActors[n]
		return ok, nil
	case groupRegexInbox.MatchString(p):
		n := groupRegexInbox.FindStringSubmatch(p)[1]
		_, ok := groupActors[n]
		return ok, nil

	// Has Outbox
	// ==========
	case personRegexOutbox.MatchString(p):
		n := personRegexOutbox.FindStringSubmatch(p)[1]
		_, ok := personActors[n]
		return ok, nil
	case groupRegexOutbox.MatchString(p):
		n := groupRegexOutbox.FindStringSubmatch(p)[1]
		_, ok := groupActors[n]
		return ok, nil

	// Has Following
	// =============
	case personRegexFollowing.MatchString(p):
		n := personRegexFollowing.FindStringSubmatch(p)[1]
		_, ok := personActors[n]
		return ok, nil
	case groupRegexFollowing.MatchString(p):
		n := groupRegexFollowing.FindStringSubmatch(p)[1]
		_, ok := groupActors[n]
		return ok, nil

	// Has Followers
	// =============
	case personRegexFollowers.MatchString(p):
		n := personRegexFollowers.FindStringSubmatch(p)[1]
		_, ok := personActors[n]
		return ok, nil
	case groupRegexFollowers.MatchString(p):
		n := groupRegexFollowers.FindStringSubmatch(p)[1]
		_, ok := groupActors[n]
		return ok, nil

	// Has Liked
	// =========
	case personRegexLiked.MatchString(p):
		n := personRegexLiked.FindStringSubmatch(p)[1]
		_, ok := personActors[n]
		return ok, nil
	case groupRegexLiked.MatchString(p):
		n := groupRegexLiked.FindStringSubmatch(p)[1]
		_, ok := groupActors[n]
		return ok, nil

	// Has Default
	// ===========
	default:
		return false, nil
	}
}
// Set should write or overwrite the value of the provided object for
// its 'id'.
func (m *Application) Set(c context.Context, o pub.PubObject) error {
	log.Println("BEGIN About to set")
	defer log.Println("END About to set")

	// Serialize the PubObject
	b, err := o.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize:%s", err)
	}
	log.Printf("Setting: %s", b)
	// Get the ID of the object that we are setting
	id := o.GetId()
	if id == nil {
		return fmt.Errorf("id is nil")
	}

	// Grab the request path from the ID
	p := id.Path

	switch {
	// SET Actor
	// =========
	case personRegexActor.MatchString(p):
		log.Println("SET person (actor) initiated")
		n := personRegexActor.FindStringSubmatch(p)[1]
		a, ok := personActors[n]
		if !ok {
			return fmt.Errorf("SET person (liked) not found")
		}
		a.actorMu.Lock()
		defer a.actorMu.Unlock()
		oc, ok := o.(vocab.PersonType)
		if !ok {
			return fmt.Errorf("setting %s but not a PersonType", id)
		}
		a.actor = oc
		return nil
	case groupRegexActor.MatchString(p):
		log.Println("SET group (actor) initiated")
		n := groupRegexActor.FindStringSubmatch(p)[1]
		a, ok := groupActors[n]
		if !ok {
			return fmt.Errorf("SET group (liked) not found")
		}
		a.actorMu.Lock()
		defer a.actorMu.Unlock()
		oc, ok := o.(vocab.GroupType)
		if !ok {
			return fmt.Errorf("setting %s but not a GroupType", id)
		}
		a.actor = oc
		return nil

	// SET Inbox
	// =========
	case personRegexInbox.MatchString(p):
		log.Println("SET person (inbox) initiated")
		n := personRegexInbox.FindStringSubmatch(p)[1]
		a, ok := personActors[n]
		if !ok {
			return fmt.Errorf("SET person (inbox) not found")
		}
		a.actorMu.Lock()
		defer a.actorMu.Unlock()
		oc, ok := o.(vocab.OrderedCollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.actor.SetInboxOrderedCollection(oc)
		return nil
	case groupRegexInbox.MatchString(p):
		log.Println("SET group (inbox) initiated")
		n := groupRegexInbox.FindStringSubmatch(p)[1]
		a, ok := groupActors[n]
		if !ok {
			return fmt.Errorf("SET group (inbox) not found")
		}
		a.actorMu.Lock()
		defer a.actorMu.Unlock()
		oc, ok := o.(vocab.OrderedCollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.actor.SetInboxOrderedCollection(oc)
		return nil

	// SET Outbox
	// ==========
	case personRegexOutbox.MatchString(p):
		log.Println("SET person (outbox) initiated")
		n := personRegexOutbox.FindStringSubmatch(p)[1]
		a, ok := personActors[n]
		if !ok {
			return fmt.Errorf("SET person (outbox) not found")
		}
		a.actorMu.Lock()
		defer a.actorMu.Unlock()
		oc, ok := o.(vocab.OrderedCollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.actor.SetOutboxOrderedCollection(oc)
		return nil
	case groupRegexOutbox.MatchString(p):
		log.Println("SET group (outbox) initiated")
		n := groupRegexOutbox.FindStringSubmatch(p)[1]
		a, ok := groupActors[n]
		if !ok {
			return fmt.Errorf("SET group (outbox) not found")
		}
		a.actorMu.Lock()
		defer a.actorMu.Unlock()
		oc, ok := o.(vocab.OrderedCollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.actor.SetOutboxOrderedCollection(oc)
		return nil

	// SET Following
	// =============
	case personRegexFollowing.MatchString(p):
		log.Println("SET person (following) initiated")
		n := personRegexFollowing.FindStringSubmatch(p)[1]
		a, ok := personActors[n]
		if !ok {
			return fmt.Errorf("SET person (following) not found")
		}
		a.actorMu.Lock()
		defer a.actorMu.Unlock()
		oc, ok := o.(vocab.CollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.actor.SetFollowingCollection(oc)
		return nil
	case groupRegexFollowing.MatchString(p):
		log.Println("SET group (following) initiated")
		n := groupRegexFollowing.FindStringSubmatch(p)[1]
		a, ok := groupActors[n]
		if !ok {
			return fmt.Errorf("SET group (following) not found")
		}
		a.actorMu.Lock()
		defer a.actorMu.Unlock()
		oc, ok := o.(vocab.CollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.actor.SetFollowingCollection(oc)
		return nil

	// SET Followers
	// =============
	case personRegexFollowers.MatchString(p):
		log.Println("SET person (followers) initiated")
		n := personRegexFollowers.FindStringSubmatch(p)[1]
		a, ok := personActors[n]
		if !ok {
			return fmt.Errorf("SET person (followers) not found")
		}
		a.actorMu.Lock()
		defer a.actorMu.Unlock()
		oc, ok := o.(vocab.CollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.actor.SetFollowersCollection(oc)
		return nil
	case groupRegexFollowers.MatchString(p):
		log.Println("SET group (followers) initiated")
		n := groupRegexFollowers.FindStringSubmatch(p)[1]
		a, ok := groupActors[n]
		if !ok {
			return fmt.Errorf("SET group (followers) not found")
		}
		a.actorMu.Lock()
		defer a.actorMu.Unlock()
		oc, ok := o.(vocab.CollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.actor.SetFollowersCollection(oc)
		return nil

	// SET Liked Collection
	// ====================
	case personRegexLiked.MatchString(p):
		log.Println("SET person (liked) initiated")
		n := personRegexLiked.FindStringSubmatch(p)[1]
		a, ok := personActors[n]
		if !ok {
			return fmt.Errorf("SET person (liked) not found")
		}
		a.actorMu.Lock()
		defer a.actorMu.Unlock()
		oc, ok := o.(vocab.CollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.actor.SetLikedCollection(oc)
		return nil
	case groupRegexLiked.MatchString(p):
		log.Println("SET group (liked) initiated")
		n := groupRegexLiked.FindStringSubmatch(p)[1]
		a, ok := groupActors[n]
		if !ok {
			return fmt.Errorf("SET group (liked) not found")
		}
		a.actorMu.Lock()
		defer a.actorMu.Unlock()
		oc, ok := o.(vocab.CollectionType)
		if !ok {
			return fmt.Errorf("setting %s but not an OrderedCollectionType", id)
		}
		a.actor.SetLikedCollection(oc)
		return nil

	// SET Document
	// ============
	case documentRegex.MatchString(p):
		log.Println("SET document initiated")
		return nil

	// SET Create
	// ==========
	case createRegex.MatchString(p):
		log.Println("SET create initiated")
		return nil

	// SET Liked Object
	// ================
	case likedRegex.MatchString(p):
		log.Println("SET liked initiated")
		return nil

	// SET Liked Object
	// ================
	case joinRegex.MatchString(p):
		log.Println("SET join initiated")
		return nil

	// SET Follow
	// ==========
	case followRegex.MatchString(p):
		log.Println("SET follow initiated")
		return nil

	// SET Default
	// ===========
	default:
		return fmt.Errorf("SET not found:", p)
	}
}

// GetInbox returns the OrderedCollection inbox of the actor for this
// context. It is up to the implementation to provide the correct
// collection for the kind of authorization given in the request.
func (m *Application) GetInbox(c context.Context, r *http.Request, rw pub.RWType) (vocab.OrderedCollectionType, error) {
	log.Printf("BEGIN GetInbox: %s", r.URL)
	defer log.Printf("END GetInbox: %s", r.URL)

	// Grab the request path
	p := (*r.URL).Path

	switch {
	// GetInbox Person
	// ===============
	case personRegexInbox.MatchString(p):
		n := personRegexInbox.FindStringSubmatch(p)[1]
		a, ok := personActors[n]
		if !ok {
			return nil, fmt.Errorf("inbox not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor.GetInboxOrderedCollection(), nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor.GetInboxOrderedCollection(), nil
		default:
			panic("Impossible")
		}

	// GetInbox Group
	// ==============
	case groupRegexInbox.MatchString(p):
		n := groupRegexInbox.FindStringSubmatch(p)[1]
		a, ok := groupActors[n]
		if !ok {
			return nil, fmt.Errorf("inbox not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor.GetInboxOrderedCollection(), nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor.GetInboxOrderedCollection(), nil
		default:
			panic("Impossible")
		}

	// GetInbox Default
	// ================
	default:
		return nil, fmt.Errorf("inbox not found")
	}
}

// GetOutbox returns the OrderedCollection inbox of the actor for this
// context. It is up to the implementation to provide the correct
// collection for the kind of authorization given in the request.
func (m *Application) GetOutbox(c context.Context, r *http.Request, rw pub.RWType) (vocab.OrderedCollectionType, error) {
	log.Printf("BEGIN GetOutbox: %s", r.URL)
	defer log.Printf("END GetOutbox: %s", r.URL)

	// Grab the request path
	p := (*r.URL).Path

	switch {
	// GetOutbox Application
	// =====================
	case applicationRegexOutbox.MatchString(p):
		applicationActor.actorMu.RLock()
		defer applicationActor.actorMu.RUnlock()
		return applicationActor.actor.GetOutboxOrderedCollection(), nil

	// GetOutbox Person
	// ================
	case personRegexOutbox.MatchString(p):
		n := personRegexOutbox.FindStringSubmatch(p)[1]
		a, ok := personActors[n]
		if !ok {
			return nil, fmt.Errorf("outbox not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor.GetOutboxOrderedCollection(), nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor.GetOutboxOrderedCollection(), nil
		default:
			panic("Impossible")
		}

	// GetOutbox Group
	// ===============
	case groupRegexOutbox.MatchString(p):
		n := groupRegexOutbox.FindStringSubmatch(p)[1]
		a, ok := groupActors[n]
		if !ok {
			return nil, fmt.Errorf("outbox not found")
		}
		switch rw {
		case pub.Read:
			a.actorMu.RLock()
			defer a.actorMu.RUnlock()
			return a.actor.GetOutboxOrderedCollection(), nil
		case pub.ReadWrite:
			a.actorMu.Lock()
			defer a.actorMu.Unlock()
			return a.actor.GetOutboxOrderedCollection(), nil
		default:
			panic("Impossible")
		}

	// GetOutbox Default
	// =================
	default:
		return nil, fmt.Errorf("outbox not found")
	}
}

// NewId takes in a client id token and returns an ActivityStreams IRI
// id for a new Activity posted to the outbox. The object is provided
// as a Typer so clients can use it to decide how to generate the IRI.
func (m *Application) NewId(c context.Context, t pub.Typer) *url.URL {
	log.Printf("BEGIN New ID")
	defer log.Printf("END New ID")

	switch {
	case vocab.HasTypeObject(t):
		log.Println("ID type Object")
	case vocab.HasTypeLink(t):
		log.Println("ID type Link")
	case vocab.HasTypeActivity(t):
		log.Println("ID type Activity")
	case vocab.HasTypeIntransitiveActivity(t):
		log.Println("ID type IntransitiveActivity")
	case vocab.HasTypeCollection(t):
		log.Println("ID type Collection")
	case vocab.HasTypeOrderedCollection(t):
		log.Println("ID type OrderedCollection")
	case vocab.HasTypeCollectionPage(t):
		log.Println("ID type CollectionPage")
	case vocab.HasTypeOrderedCollectionPage(t):
		log.Println("ID type OrderedCollectionPage")
	case vocab.HasTypeAccept(t):
		log.Println("ID type Accept")
	case vocab.HasTypeTentativeAccept(t):
		log.Println("ID type TentativeAccept")
	case vocab.HasTypeAdd(t):
		log.Println("ID type Add")
	case vocab.HasTypeArrive(t):
		log.Println("ID type Arrive")
	case vocab.HasTypeCreate(t):
		log.Println("ID type Create")
		u, err := url.Parse(fmt.Sprintf("%s/create", baseUrl))
		if err != nil {
			log.Fatal("failed to parse url:", err)
		}
		return u
	case vocab.HasTypeDelete(t):
		log.Println("ID type Delete")
	case vocab.HasTypeFollow(t):
		log.Println("ID type Follow")
		u, err := url.Parse(fmt.Sprintf("%s/follow", baseUrl))
		if err != nil {
			log.Fatal("failed to parse url:", err)
		}
		return u
	case vocab.HasTypeIgnore(t):
		log.Println("ID type Ignore")
	case vocab.HasTypeJoin(t):
		log.Println("ID type Join")
		u, err := url.Parse(fmt.Sprintf("%s/join", baseUrl))
		if err != nil {
			log.Fatal("failed to parse url:", err)
		}
		return u
	case vocab.HasTypeLeave(t):
		log.Println("ID type Leave")
	case vocab.HasTypeLike(t):
		log.Println("ID type Like")
		u, err := url.Parse(fmt.Sprintf("%s/liked", baseUrl))
		if err != nil {
			log.Fatal("failed to parse url:", err)
		}
		return u
	case vocab.HasTypeOffer(t):
		log.Println("ID type Offer")
	case vocab.HasTypeInvite(t):
		log.Println("ID type Invite")
	case vocab.HasTypeReject(t):
		log.Println("ID type Reject")
	case vocab.HasTypeTentativeReject(t):
		log.Println("ID type TentativeReject")
	case vocab.HasTypeRemove(t):
		log.Println("ID type Remove")
	case vocab.HasTypeUndo(t):
		log.Println("ID type Undo")
	case vocab.HasTypeUpdate(t):
		log.Println("ID type Update")
	case vocab.HasTypeView(t):
		log.Println("ID type View")
	case vocab.HasTypeListen(t):
		log.Println("ID type Listen")
	case vocab.HasTypeRead(t):
		log.Println("ID type Read")
	case vocab.HasTypeMove(t):
		log.Println("ID type Move")
	case vocab.HasTypeTravel(t):
		log.Println("ID type Travel")
	case vocab.HasTypeAnnounce(t):
		log.Println("ID type Announce")
	case vocab.HasTypeBlock(t):
		log.Println("ID type Block")
	case vocab.HasTypeFlag(t):
		log.Println("ID type Flag")
	case vocab.HasTypeDislike(t):
		log.Println("ID type Dislike")
	case vocab.HasTypeQuestion(t):
		log.Println("ID type Question")
	case vocab.HasTypeApplication(t):
		log.Println("ID type Application")
	case vocab.HasTypeGroup(t):
		log.Println("ID type Group")
	case vocab.HasTypeOrganization(t):
		log.Println("ID type Organization")
	case vocab.HasTypePerson(t):
		log.Println("ID type Person")
	case vocab.HasTypeService(t):
		log.Println("ID type Service")
	case vocab.HasTypeRelationship(t):
		log.Println("ID type Relationship")
	case vocab.HasTypeArticle(t):
		log.Println("ID type Article")
	case vocab.HasTypeDocument(t):
		log.Println("ID type Document")
		u, err := url.Parse(fmt.Sprintf("%s/document", baseUrl))
		if err != nil {
			log.Fatal("failed to parse url:", err)
		}
		return u
	case vocab.HasTypeAudio(t):
		log.Println("ID type Audio")
	case vocab.HasTypeImage(t):
		log.Println("ID type Image")
	case vocab.HasTypeVideo(t):
		log.Println("ID type Video")
	case vocab.HasTypeNote(t):
		log.Println("ID type Note")
	case vocab.HasTypePage(t):
		log.Println("ID type Page")
	case vocab.HasTypeEvent(t):
		log.Println("ID type Event")
	case vocab.HasTypePlace(t):
		log.Println("ID type Place")
	case vocab.HasTypeProfile(t):
		log.Println("ID type Profile")
	case vocab.HasTypeTombstone(t):
		log.Println("ID type Tombstone")
	case vocab.HasTypeMention(t):
		log.Println("ID type Mention")
	default:
		log.Println("ID type DEFAULT")
	}

	return nil
}

// GetPublicKey fetches the public key for a user based on the public
// key id. It also determines which algorithm to use to verify the
// signature.
func (m *Application) GetPublicKey(c context.Context, publicKeyId string) (crypto.PublicKey, httpsig.Algorithm, *url.URL, error) {
	log.Println("BEGIN Getting Public Key")
	defer log.Println("END Getting Public Key")

	query := "SELECT fingerprint, public_key, name FROM persons;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, httpsig.RSA_SHA256, nil, fmt.Errorf("error querying database:%s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var fingerprint string
		var publicKey string
		var name string
		err = rows.Scan(&fingerprint, &publicKey, name)
		if err != nil {
			return nil, httpsig.RSA_SHA256, nil, err
		}
		if fingerprint == publicKeyId {
			u, err := url.Parse(fmt.Sprintf("%s/activity/person/%s", baseUrl, name))
			if err != nil {
				return nil, httpsig.RSA_SHA256, nil, fmt.Errorf("error parsing url:%s", err)
			}

			return publicKey, httpsig.RSA_SHA256, u, nil
		}
	}

	return nil, httpsig.RSA_SHA256, nil, fmt.Errorf("not implemented: GetPublicKey")
}

// CanAdd returns true if the provided object is allowed to be added to
// the given target collection. Applicable to either or both of the
// SocialAPI and FederateAPI.
func (m *Application) CanAdd(c context.Context, o vocab.ObjectType, t vocab.ObjectType) bool {
	log.Println("BEGIN Checking canAdd")
	defer log.Println("END Checking canAdd")

	// TODO: Implement CanAdd
	return true
}

// CanRemove returns true if the provided object is allowed to be
// removed from the given target collection. Applicable to either or
// both of the SocialAPI and FederateAPI.
func (m *Application) CanRemove(c context.Context, o vocab.ObjectType, t vocab.ObjectType) bool {
	log.Println("BEGIN Checking canRemove")
	defer log.Println("END Checking canRemove")

	// TODO: Implement CanRemove
	return true
}

// FederateAPI
// ===========
//
// FederateAPI is provided by users of the go-fed/activity library and
// designed to handle receiving messages from ActivityPub servers
// through the Federative API.

// OnFollow determines whether to take any automatic reactions in
// response to this follow. Note that if this application does not own
// an object on the activity, then the 'AutomaticAccept' and
// 'AutomaticReject' results will behave as if they were 'DoNothing'.
// FollowResponse instructs how to proceed upon immediately receiving a request
// to follow.
func (m *Application) OnFollow(c context.Context, s *streams.Follow) pub.FollowResponse {
	log.Println("BEGIN FederateAPI OnFollow")
	defer log.Println("END FederateAPI OnFollow")

	return pub.AutomaticAccept
	// return pub.AutomaticReject
	// return pub.DoNothing
}

// Unblocked should return an error if the provided actor ids are not
// able to interact with this particular end user due to being blocked
// or other application-specific logic. This error is passed
// transparently back to the request thread via PostInbox.
//
// If nil error is returned, then the received activity is processed as
// a normal unblocked interaction.
func (m *Application) Unblocked(c context.Context, actorIRIs []*url.URL) error {
	log.Println("BEGIN FederateAPI Unblocked")
	defer log.Println("END FederateAPI Unblocked")

	return nil
}

// FilterForwarding is invoked when a received activity needs to be
// forwarded to specific inboxes owned by this server in order to avoid
// the ghost reply problem. The IRIs provided are collections owned by
// this server that the federate peer requested inbox forwarding to.
//
// Implementors must apply some sort of filtering to the provided IRI
// collections. Without any filtering, the resulting application is
// vulnerable to becoming a spam bot for a malicious federate peer.
// Typical implementations will filter the iris down to be only the
// follower collections owned by the actors targeted in the activity.
func (m *Application) FilterForwarding(c context.Context, activity vocab.ActivityType, iris []*url.URL) ([]*url.URL, error) {
	log.Println("BEGIN FederateAPI FilterFollowing")
	defer log.Println("END FederateAPI FilterFollowing")

	return iris, nil
}

// NewSigner returns a new httpsig.Signer for which deliveries can be
// signed by the actor delivering the Activity. Let me take this moment
// to really level with you, dear anonymous reader-of-documentation. You
// want to use httpsig.RSA_SHA256 as the algorithm. Otherwise, your app
// will not federate correctly and peers will reject the signatures. All
// other known implementations using HTTP Signatures use RSA_SHA256,
// hardcoded just like your implementation will be.
//
// The headers available for inclusion in the signature are:
//     Date
//     User-Agent
func (m *Application) NewSigner() (httpsig.Signer, error) {
	log.Println("BEGIN FederateAPI NewSigner")
	defer log.Println("END FederateAPI NewSigner")

	prefs := []httpsig.Algorithm{httpsig.RSA_SHA256}
	// headersToSign := []string{httpsig.RequestTarget, "date", "digest"}
	headersToSign := []string{httpsig.RequestTarget, "date"}
	signer, _, err := httpsig.NewSigner(prefs, headersToSign, httpsig.Signature)
	if err != nil {
		return nil, fmt.Errorf("error in NewSigner:", err)
	}
	return signer, nil
}

// PrivateKey fetches the private key and its associated public key ID.
// The given URL is the inbox or outbox for the actor whose key is
// needed.
func (m *Application) PrivateKey(boxIRI *url.URL) (privKey crypto.PrivateKey, pubKeyId string, err error) {
	log.Println("BEGIN FederateAPI PrivateKey")
	defer log.Println("END FederateAPI PrivateKey")

	p := boxIRI.Path

	fetchPrivateKey := func(n string) (crypto.PrivateKey, string, error) {
		query := "SELECT fingerprint, private_key FROM persons WHERE name=$1;"
		var fingerprint string
		var privateKey string
		err := db.QueryRow(query, n).Scan(&fingerprint, &privateKey)
		if err != nil {
			return nil, "", fmt.Errorf("error querying database:%s", err)
		}
		block, _ := pem.Decode([]byte(privateKey))
		if block == nil {
			return nil, "", errors.New("failed to parse PEM block containing the key")
		}
		priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, "", fmt.Errorf("error reading in private key:%s", err)
		}
		return priKey, fingerprint, nil
	}

	switch {
	case personRegexInbox.MatchString(p):
		n := personRegexInbox.FindStringSubmatch(p)[1]
		return fetchPrivateKey(n)
	case personRegexOutbox.MatchString(p):
		n := personRegexOutbox.FindStringSubmatch(p)[1]
		return fetchPrivateKey(n)
	default:
		return nil, "", fmt.Errorf("PrivateKey regex failed to match")
	}
}

