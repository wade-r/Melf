package main

import (
	"log"
	"strings"
	"time"
)

// max size of a multi-line message
const MessageMaxSize = 512

type Message map[string]interface{}

// Session represents a context related with a instance_id
type Session struct {
	message   Message
	count     uint
	updatedAt time.Time
}

// existed sessions, key is instance_id
var sessions = make(map[string]*Session)

// sessionNewMessageChan for handle session
var newMessageChan = make(chan Message, 128)

func dispatchMessageToSessions(m Message) {
	newMessageChan <- m
}

// main loop for session processing
// WARN: all operations related with 'sessions' variable should be scheduled here
func sessionManagementLoop() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case m := <-newMessageChan:
			{
				// check existence of 'short_message'
				s, _ := m[kShortMessage].(string)
				if len(s) == 0 {
					log.Println("short_message field not found")
					continue
				}

				// find or create Session
				id := m[kContainerId].(string)
				sess := sessions[id]
				if sess == nil {
					sess = &Session{}
					sessions[id] = sess
				}

				if sess.message != nil {
					// previous message exists
					os := sess.message[kShortMessage].(string)

					if isPartial(s) && len(os) < MessageMaxSize {
						// append partial message
						sess.message[kShortMessage] = os + "\r\n" + s
						sess.count = sess.count + 1
						sess.updatedAt = time.Now()
					} else {
						// replace the message and send last one
						go sendMessage(sess.message)

						sess.message = m
						sess.count = 1
						sess.updatedAt = time.Now()
					}
				} else {
					// previous message does not exist
					sess.message = m
					sess.count = 1
					sess.updatedAt = time.Now()
				}
			}
		case now := <-ticker.C:
			{
				// every 1 second, check existing sessions
				var expires = make([]string, 0)

				for id, sess := range sessions {
					if sess.message == nil {
						// empty session, GC after expired for 5 seconds
						if now.Sub(sess.updatedAt) > 5*time.Second {
							expires = append(expires, id)
						}
					} else {
						// message cached for more than 1 second, send and empty the session
						if now.Sub(sess.updatedAt) > 1*time.Second {
							// send message
							go sendMessage(sess.message)

							// empty the session
							sess.message = nil
							sess.count = 0
							sess.updatedAt = now
						}
					}
				}

				// do remove expired sessions
				for _, id := range expires {
					delete(sessions, id)
				}
			}
		}
	}
}

func isPartial(s string) bool {
	return strings.HasPrefix(s, " ") || strings.HasPrefix(s, "\t") || strings.HasPrefix(s, "\r") || strings.HasPrefix(s, "\n")
}
