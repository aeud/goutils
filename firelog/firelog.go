package firelog

import (
	"fmt"
	firego "gopkg.in/zabawaba99/firego.v1"
	"log"
	"strings"
	"sync"
)

type FirebaseMessage struct {
	key     string
	message string
}

func NewFirebaseMessage(k, c string) *FirebaseMessage {
	m := new(FirebaseMessage)
	m.key = k
	m.message = c
	return m
}

func (m *FirebaseMessage) Output() {
	log.Printf("Log %v: %v\n", strings.Replace(m.key, ".", "/", -1), m.message)
}

type FirebaseService struct {
	stack     chan *FirebaseMessage
	wg        *sync.WaitGroup
	authToken string
	endpoint  string
	ref       string
}

func NewFirebaseService(endpoint, authToken, ref string) *FirebaseService {
	s := new(FirebaseService)
	s.stack = make(chan *FirebaseMessage)
	s.wg = new(sync.WaitGroup)
	s.authToken = authToken
	s.endpoint = endpoint
	s.ref = ref
	return s
}

func (s *FirebaseService) Run() {
	for {
		s.Write(<-s.stack)
	}
}

func (s *FirebaseService) Deamon() {
	go s.Run()
}

func (s *FirebaseService) Write(m *FirebaseMessage) {
	url := fmt.Sprintf("%v/%v/%v", s.endpoint, s.ref, strings.Replace(m.key, ".", "/", -1))
	f := firego.New(url, nil)
	f.Auth(s.authToken)
	v := m.message
	// log.Printf("Log %v: %v\n", url, v)
	if err := f.Set(v); err != nil {
		log.Fatalf("Firebase Auth: %v\n", err)
	}
	s.wg.Done()
}

func (s *FirebaseService) Push(m *FirebaseMessage) {
	s.wg.Add(1)
	s.stack <- m
}

func (s *FirebaseService) AsyncPush(m *FirebaseMessage) {
	m.Output()
	s.wg.Add(1)
	go func(s *FirebaseService) {
		s.stack <- m
	}(s)
}

func (s *FirebaseService) Wait() {
	s.wg.Wait()
}
