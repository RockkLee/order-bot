package jwtutil

import (
	"order-bot-mgmt-svc/internal/models"
	"reflect"
	"testing"
	"time"
)

var (
	secret = []byte("test")
	now    = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	claims = models.Claims{
		Sub:   "userIdTest",
		Email: "userEmailTest",
		Exp:   now.Add(30 * time.Minute).Unix(),
		Iat:   now.Unix(),
		Typ:   "access",
	}
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VySWRUZXN0IiwiZW1haWwiOiJ1c2VyRW1haWxUZXN0IiwiZXhwIjo5NDY2ODY2MDAsImlhdCI6OTQ2Njg0ODAwLCJ0eXAiOiJhY2Nlc3MifQ.c1w10izE7WDMW_Sb4KgIYZ--tCKDFLhNNzXnhfgIpJQ"
)

func TestSignJWT(t *testing.T) {
	type args struct {
		secret []byte
		claims models.Claims
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "happy path - access token",
			args: args{
				secret: secret,
				claims: claims,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SignJWT(tt.args.secret, tt.args.claims)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log("got: ", got)
			// if got != tt.want {
			// 	t.Errorf("SignJWT() got = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestParseJWT(t *testing.T) {
	type args struct {
		secret []byte
		token  string
	}
	tests := []struct {
		name    string
		args    args
		want    models.Claims
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				secret: secret,
				token:  token,
			},
			want:    claims,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJWT(tt.args.secret, tt.args.token, now)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseJWT() got = %v, want %v", got, tt.want)
			}
		})
	}
}
