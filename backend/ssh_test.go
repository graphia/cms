package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublicKey_User(t *testing.T) {

	db.Drop("User")
	db.Drop("PublicKey")

	_ = createUser(ds)
	expected, _ := getUserByUsername(ds.Username)

	pkWithUser := PublicKey{UserID: expected.ID}
	pkWithoutUser := PublicKey{UserID: 999}

	_ = db.Save(&pkWithUser)
	db.One("UserID", expected.ID, &pkWithUser)

	type args struct {
		pk PublicKey
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errMsg  string
		want    User
	}{
		{
			name: "Existing User",
			args: args{
				pk: pkWithUser,
			},
			want: expected,
		},
		{
			name: "Missing User",
			args: args{
				pk: pkWithoutUser,
			},
			wantErr: true,
			errMsg:  "Could not find user 999",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			actual, err := tt.args.pk.User()

			if tt.wantErr && err == nil {
				t.Fatal("Error expected, none found")
			}

			if tt.wantErr {
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}

			assert.Equal(t, tt.want, actual)

		})
	}
}
