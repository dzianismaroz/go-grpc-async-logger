package acl

import (
	"encoding/json"
	"fmt"
	"path"

	. "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Build ACL list and provide authentication by principal SID anf requested path.
type (
	sid              = string                      // SID ( identity of user / module/ path of system).
	objectIdentities = []string                    // Transient holder for raw parsed data.
	objectIdentity   = string                      // Path.
	pathSet          = map[objectIdentity]struct{} // Persist paths as set for more handy access.

	ACLS struct {
		content map[sid]pathSet // Paths allowed for certain client.
	}
)

// Authenticate certain @sid (user) by @path requested.
// We only read ACL permissions from different goroutines ( so it's safe concurrent).
func (a *ACLS) Authenticate(sid []string, path_ string) error {
	if len(sid) == 0 {
		return status.Error(Unauthenticated, "empty sid")
	}
	objsId, ok := a.content[sid[0]]
	if !ok {
		return status.Error(Unauthenticated, "not found")
	}
	if _, allowed := objsId[path_]; !allowed {
		// Now try wildcard.
		if _, allowed := objsId[path.Dir(path_)+"/*"]; !allowed {
			return status.Error(Unauthenticated, "permission denied")
		}
	}
	return nil
}

// Build map of maps for faster resolving ACL perm.
// (It should  works faster than iteration on every slice element).
// 1. Request by SID . IF presented then :
// 2. Get request path from map and reply final result.
func BuildACL(raw string) (*ACLS, error) {
	var resMap map[sid]objectIdentities // holder of raw content.
	err := json.Unmarshal([]byte(raw), &resMap)
	if err != nil { //Seems invalid map provided as ACL list.
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
