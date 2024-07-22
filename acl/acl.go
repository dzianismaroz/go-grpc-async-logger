package acl

import (
	"encoding/json"
	"errors"
	"fmt"
)

// To build ACL list and provide authentication by SID anf requested path.
type (
	sid              = string   // SID ( identity of user / module/ path of system).
	objectIdentities = []string // Transient holder for raw parsed data.
	objectIdentity   = string   // Path.
	pathSet          = map[objectIdentity]struct{}

	ACLS struct {
		content map[sid]pathSet // Paths allowed for certain client.
	}
)

// Authenticate certain @sid (user) by @path requested.
// We only read ACL permissions from different goroutines ( so it's safe concurrent).
func (a *ACLS) Authenticate(sid string, path string) (bool, error) {
	objsId, ok := a.content[sid]
	if !ok {
		return false, errors.New("not permitted")
	}
	_, allowed := objsId[path]
	if !allowed {
		return false, nil
	}
	return true, nil
}

// Build map of maps for faster resolving ACL perm.
// (It should  works faster than iteration on every slice element).
// 1. Request by SID . IF presented then :
// 2. Get request path from map and reply final result.
func BuildACL(raw string) (*ACLS, error) {
	var resMap map[sid]objectIdentities // holder of raw content.
	err := json.Unmarshal([]byte(raw), &resMap)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ACL config: %s", err.Error())
	}
	result := make(map[sid]pathSet, len(resMap))
	for k, v := range resMap {
		temp := make(pathSet, len(v))
		for i := 0; i < len(v); i++ {
			temp[v[i]] = struct{}{}
		}
		result[k] = temp
	}
	return &ACLS{content: result}, nil
}
