package irc

import (
	"code.google.com/p/go-uuid/uuid"
)

type User struct {
	Name     string
	Uname    string
	Rname    string
	Msgs     <-chan *Message
	inMsgs   MChannel
	server   *Server
	handlers []HandlerFunc
	Handlers map[string]HandlerFunc
	old      bool
}

func (s *Server) NewUser(nname, uname string) *User {
	u := &User{
		Name:     nname,
		Uname:    uname,
		inMsgs:   make(chan *Message),
		server:   s,
		Handlers: make(map[string]HandlerFunc),
	}

	u.Msgs = u.inMsgs.Queue()

	nicks := func(m *Message) {
		if m.Type == NICK {
		}
	}
	u.handlers = []HandlerFunc{nicks}

	if u2, ok := s.Users[nname]; ok {
		*u2 = *u
	} else {
		s.Users[nname] = u
	}

	go func() {
		for m := range u.Msgs {
			for _, f := range u.handlers {
				go f(m)
			}
			for _, f := range u.Handlers {
				go f(m)
			}
		}
	}()
	return u
}

func (u *User) addHandler(f HandlerFunc) {
	u.handlers = append(u.handlers, f)
}

func (u *User) AddHandler(name string, f HandlerFunc) string {
	if name == "" {
		name = uuid.New()
	}
	u.Handlers[name] = f
	return name
}

func (u *User) RemHandler(name string) {
	delete(u.Handlers, name)
}

func (this *User) Say(s string) {
	this.server.PrivMsg(this, s)
}

func (this *User) Write(p []byte) (int, error) {
	return this.server.Write(append([]byte("PRIVMSG "+this.Name+" :"), p...))
}

func (this *User) String() string {
	if this == nil {
		return ""
	}
	return this.Name
}

func (this *User) Ident() string {
	if this == nil {
		return ""
	}
	return this.Name
}

func (u *User) Old() bool {
	return u.old
}
