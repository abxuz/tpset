package mercury

import (
	"context"
	"testing"
)

func TestLogin(t *testing.T) {
	s := NewService("http://100.64.38.17:6080")
	err := s.Login(context.Background(), "admin", "a0e946e18a1313debcda76bb81ff67a022517e8846e411dfbbddc030556f2a00316cf116de9a6403c48381be57d80504e5b1e4c0abd2d648027e1341ac98136320ed643986fc021be2016ee0baad05f308cda274e094e260e7ab6030dea73b70a7965f946e6fec5e51f9999d7ae8693a83c877efa918224041d7d4911af9aa7f")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s.stok, s.cookies)
}

func TestListGroup(t *testing.T) {
	s := NewService("http://100.64.38.17:6080")
	err := s.Login(context.Background(), "admin", "a0e946e18a1313debcda76bb81ff67a022517e8846e411dfbbddc030556f2a00316cf116de9a6403c48381be57d80504e5b1e4c0abd2d648027e1341ac98136320ed643986fc021be2016ee0baad05f308cda274e094e260e7ab6030dea73b70a7965f946e6fec5e51f9999d7ae8693a83c877efa918224041d7d4911af9aa7f")
	if err != nil {
		t.Fatal(err)
	}
	list, err := s.ListGroup(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}
