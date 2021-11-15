package api

import (
	"github.com/golang/mock/gomock"
	mockdb "github.com/simplebank/db/mock"
	db "github.com/simplebank/db/sqlc"
	"github.com/simplebank/db/util"
	"testing"
)

func TestGetAccountAPI(t *testing.T) {
	account := RandomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// build stubs
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

}



func RandomAccount() db.Account{
	return db.Account{
		ID: util.RandomInt(0, 100),
		Owner: util.RandomOwner(),
		Balance: util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}