package db

// TEMPORARY local-only shim to work around a pre-existing stale generated
// mock. This file is deleted after verification and is never committed.

import sphinxdb "github.com/stakwork/sphinx-tribes/db"

func (_m *Database) GetBountyByUnlockCode(code string) (sphinxdb.NewBounty, error) {
	ret := _m.Called(code)

	var r0 sphinxdb.NewBounty
	if rf, ok := ret.Get(0).(func(string) sphinxdb.NewBounty); ok {
		r0 = rf(code)
	} else {
		r0 = ret.Get(0).(sphinxdb.NewBounty)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
