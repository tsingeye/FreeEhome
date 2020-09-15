package config

import "sync"

var StreamSession *sessions

type sessions struct {
	SessionMap map[string]*SessionInfo
	sync.RWMutex
}

type SessionInfo struct {
	Session    string
	DeviceID   string
	ChannelID  string
	SessionURL map[string]string
}

func newStreamSession() *sessions {
	return &sessions{
		SessionMap: make(map[string]*SessionInfo),
	}
}

func (s *sessions) Set(channelID string, sessionInfo *SessionInfo) {
	s.Lock()
	defer s.Unlock()
	s.SessionMap[channelID] = sessionInfo
}

func (s *sessions) Get(channelID string) *SessionInfo {
	s.RLock()
	defer s.RUnlock()
	return s.SessionMap[channelID]
}

func (s *sessions) GetAll() map[string]*SessionInfo {
	s.RLock()
	defer s.RUnlock()
	return s.SessionMap
}

//过滤查询session对应的channelID
func (s *sessions) Filter(session string) (string, string) {
	for _, v := range s.GetAll() {
		if v.Session == session {
			return v.DeviceID, v.ChannelID
		}
	}

	return "", ""
}

func (s *sessions) Delete(channelID string) {
	s.Lock()
	defer s.Unlock()
	delete(s.SessionMap, channelID)
}

func init() {
	StreamSession = newStreamSession()
}
