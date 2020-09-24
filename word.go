package wordcounter

import (
	"fmt"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

// SearchKeyword: When the bot receives a text message which has SearchKeyword,
// the bot will reply a list of word counting.
// MaxResultCount is the max number of users who are listed in the word counting list.
// Hours is the duration which the bot records messages.
type Config struct {
	SearchKeyword  string
	MaxResultCount int
	Hours          int
}

// wechat message record, only text content.
type msg struct {
	time    time.Time
	contact string
	count   int
}

// When the message count is above startCleanCount, remove outdated messages.
const startCleanCount = 3000

type roomData struct {
	roomName string
	data     []msg
}

var rooms []*roomData

func getRoom(roomName string) *roomData {
	for _, r := range rooms {
		if r.roomName == roomName {
			return r
		}
	}
	return nil
}

type contactResult struct {
	name  string
	count int
}

type contactResultSlice []contactResult

// from big to small
func (s contactResultSlice) Less(i, j int) bool { return s[i].count > s[j].count }
func (s contactResultSlice) Swap(i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}
func (s contactResultSlice) Len() int { return len(s) }

func getResult(room *roomData, duration time.Duration, resultCount int) string {
	var result string
	countMap := make(map[string]int)
	for _, msg := range room.data {
		if time.Now().Sub(msg.time) < duration {
			countMap[msg.contact] += msg.count
		}
	}

	var sortSlice contactResultSlice
	for name, count := range countMap {
		sortSlice = append(sortSlice, contactResult{
			name:  name,
			count: count,
		})
	}
	sort.Sort(sortSlice)
	number := 0
	for _, res := range sortSlice {
		result += fmt.Sprintf("%s: %d\n", res.name, res.count)
		number++
		if number >= resultCount {
			break
		}
	}
	return result
}

var mux sync.Mutex

func New(config Config) *wechaty.Plugin {
	plugin := wechaty.NewPlugin()
	plugin.OnMessage(func(context *wechaty.Context, message *user.Message) {
		if message.Type() == schemas.MessageTypeText {
			mux.Lock()
			// Search
			if strings.Contains(message.Text(), config.SearchKeyword) {
				room := getRoom(message.Room().ID())
				if room != nil {
					result := getResult(room, time.Duration(config.Hours)*time.Hour, config.MaxResultCount)
					_, err := message.Room().Say(result)
					if err != nil {
						log.Print(err)
					}
				}
				mux.Unlock()
				return
			}

			// record words
			roomName := message.Room().ID()
			room := getRoom(roomName)
			if room == nil {
				room = new(roomData)
				room.roomName = roomName
				rooms = append(rooms, room)
			}

			name := message.From().Name()
			count := len([]rune(message.Text()))
			room.data = append(room.data, msg{
				time:    time.Now(),
				contact: name,
				count:   count,
			})

			// clean old messages
			if len(room.data) > startCleanCount {
				var newData []msg
				for _, msg := range room.data {
					if time.Now().Sub(msg.time) < time.Duration(config.Hours)*time.Hour {
						newData = append(newData, msg)
					}
				}
				room.data = newData
			}

			mux.Unlock()
		}
	})
	return plugin
}
