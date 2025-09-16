package structures

import (
	"log"
	"sync"
)

type ConnectionSessionState struct {
	Hash             string
	ExaminerSockAddr string
	UserSockAddr     string
	mutex            sync.Mutex
}

func NewConnectionSessionState(hash string) *ConnectionSessionState {
	log.Printf("Создана новая сессия: %s", hash)
	return &ConnectionSessionState{
		Hash: hash,
	}
}

func (conn_session_state *ConnectionSessionState) SetExaminerSockAddr(addr string) {
	conn_session_state.mutex.Lock()
	defer conn_session_state.mutex.Unlock()
	conn_session_state.ExaminerSockAddr = addr
}

func (conn_session_state *ConnectionSessionState) SetUserSockAddr(addr string) {
	conn_session_state.mutex.Lock()
	defer conn_session_state.mutex.Unlock()
	conn_session_state.UserSockAddr = addr
}

func (conn_session_state *ConnectionSessionState) GetExaminerSockAddr() string {
	conn_session_state.mutex.Lock()
	defer conn_session_state.mutex.Unlock()
	return conn_session_state.ExaminerSockAddr
}

func (conn_session_state *ConnectionSessionState) GetUserSockAddr() string {
	conn_session_state.mutex.Lock()
	defer conn_session_state.mutex.Unlock()
	return conn_session_state.UserSockAddr
}

func (conn_session_state *ConnectionSessionState) SessionHaveExaminerSockAddr() bool {
	conn_session_state.mutex.Lock()
	defer conn_session_state.mutex.Unlock()

	if conn_session_state.ExaminerSockAddr != "" {
		return true
	} else {
		return false
	}

}

func (conn_session_state *ConnectionSessionState) SessionHaveUserSockAddr() bool {
	conn_session_state.mutex.Lock()
	defer conn_session_state.mutex.Unlock()

	if conn_session_state.UserSockAddr != "" {
		return true
	} else {
		return false
	}

}

func (conn_session_state *ConnectionSessionState) GetHash() string {
	conn_session_state.mutex.Lock()
	defer conn_session_state.mutex.Unlock()
	return conn_session_state.Hash
}

func (conn_session_state *ConnectionSessionState) Copy() *ConnectionSessionState {
	conn_session_state.mutex.Lock()
	defer conn_session_state.mutex.Unlock()

	return &ConnectionSessionState{
		Hash:             conn_session_state.Hash,
		ExaminerSockAddr: conn_session_state.ExaminerSockAddr,
		UserSockAddr:     conn_session_state.UserSockAddr,
	}
}
