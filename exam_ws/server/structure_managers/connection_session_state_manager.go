package structure_managers

import (
	"exam_ws/structures"
	"log"
	"net"
	"strings"
	"sync"
)

type ConnectionSessionStateManager struct {
	states map[string]*structures.ConnectionSessionState
	mutex  sync.Mutex
}

func (sm *ConnectionSessionStateManager) RemoveConnect(conn net.Conn, hash string, user_role string) {
	current_client_addr := conn.RemoteAddr().String()
	sesstion_state, _ := sm.GetConnectionSessionState(hash)
	var session_state_client_addr string

	if user_role == "examiner" {
		session_state_client_addr = sesstion_state.GetExaminerSockAddr()
		if session_state_client_addr == current_client_addr {
			sm.SetExaminerSockAddr(hash, "")
			log.Printf("Пользователь с ролью %s был отключен и удален из статуса в сесии %s", user_role, hash)
		} else {
			log.Printf("Пользователь был отключен")
		}
	} else {
		session_state_client_addr = sesstion_state.GetUserSockAddr()
		if session_state_client_addr == current_client_addr {
			sm.SetUserSockAddr(hash, "")
			log.Printf("Пользователь с ролью %s был отключен и удален из статуса в сесии %s", user_role, hash)
		} else {
			log.Printf("Пользователь был отключен")
		}
	}

	conn.Close()
}

func NewConnectionSessionStateManager() *ConnectionSessionStateManager {
	return &ConnectionSessionStateManager{
		states: make(map[string]*structures.ConnectionSessionState),
	}
}

func (sm *ConnectionSessionStateManager) NewConnectionSessionState(conn net.Conn, msg string) bool {
	var exsists bool
	var hash string

	if strings.Contains(msg, "/") {
		msg_parts := strings.Split(msg, "/")
		hash := msg_parts[1]
		sm.mutex.Lock()
		_, exsists = sm.states[hash]
		sm.mutex.Unlock()
	} else {
		_, exsists = sm.states[msg]
		hash = msg
	}

	if exsists {
		return false
	} else {
		sm.mutex.Lock()
		defer sm.mutex.Unlock()

		sm.states[hash] = structures.NewConnectionSessionState(hash)
		return true
	}
}

func (sm *ConnectionSessionStateManager) GetConnectionSessionState(hash string) (*structures.ConnectionSessionState, bool) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	conn_session_state, exists := sm.states[hash]

	return conn_session_state, exists
}

func (sm *ConnectionSessionStateManager) RemoveConnectionSessionState(hash string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	delete(sm.states, hash)
}

func (sm *ConnectionSessionStateManager) SetExaminerSockAddr(hash string, addr string) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	state, _ := sm.states[hash]

	if !state.SessionHaveExaminerSockAddr() {
		state.SetExaminerSockAddr(addr)
		return true
	}

	return false
}

func (sm *ConnectionSessionStateManager) SetUserSockAddr(hash string, addr string) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	state, _ := sm.states[hash]

	if !state.SessionHaveUserSockAddr() {
		state.SetUserSockAddr(addr)
		return true
	}

	return false
}

func (sm *ConnectionSessionStateManager) GetAllConnectionSessionStates() map[string]structures.ConnectionSessionState {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	result := make(map[string]structures.ConnectionSessionState)
	for hash, state := range sm.states {
		result[hash] = *state.Copy()
	}
	return result
}
