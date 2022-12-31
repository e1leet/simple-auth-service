package manager

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager_HashPassword(t *testing.T) {
	testcases := []struct {
		password string
		salt     string
		expected string
	}{
		{
			password: "test",
			salt:     "test",
			expected: "37268335dd6931045bdcdf92623ff819a64244b53d0e746d438797349d4da578",
		},
		{
			password: "amazing",
			salt:     "test",
			expected: "7fdbe69446f6434631f403a93eb2ccf745eda12837d0b0ae920fb6255b623fa3",
		},
		{
			password: "test",
			salt:     "wow",
			expected: "9f95811c9583d55ca0bb0c594f88af0e1c4c5b8a86ab175d642482c0e90b8a8e",
		},
		{
			password: "blabla",
			salt:     "blabla",
			expected: "88f2b02f60d21173be8d761120379d97ad427cdf29e31a5884cf8ea5fdbb0583",
		},
	}
	for _, tt := range testcases {
		name := fmt.Sprintf("hash_password__password:%s,salt:%s", tt.password, tt.salt)

		t.Run(name, func(t *testing.T) {
			m := New(tt.salt)
			assert.Equal(t, tt.expected, m.HashPassword(tt.password))
		})
	}
}

func TestManager_CheckPassword(t *testing.T) {
	testcases := []struct {
		password       string
		hashedPassword string
		salt           string
		expected       bool
	}{
		{
			password:       "test",
			hashedPassword: "37268335dd6931045bdcdf92623ff819a64244b53d0e746d438797349d4da578",
			salt:           "test",
			expected:       true,
		},
		{
			password:       "amazing",
			hashedPassword: "7fdbe69446f6434631f403a93eb2ccf745eda12837d0b0ae920fb6255b623fa3",
			salt:           "test",
			expected:       true,
		},
		{
			password:       "blabla",
			hashedPassword: "88f2b02f60d21173be8d761120379d97ad427cdf29e31a5884cf8ea5fdbb0583",
			salt:           "blabla",
			expected:       true,
		},
		{
			password:       "bdaslabla",
			hashedPassword: "dawdasd37hg87hh288jjq09jwd8hg937h49y5-38489jda8jdisakjd3j28r7h4n",
			salt:           "dsnauds",
			expected:       false,
		},
		{
			password:       "a9sjd98jasd",
			hashedPassword: "dausnd983n29jsd8j23289jdisoaijdiwasjd38jdijasdiojfejfeifddiojdio",
			salt:           "djiasdj",
			expected:       false,
		},
		{
			password:       "a9sjd98jasd",
			hashedPassword: "dasdasdwqwadsdwasdwasdd73h273he372d7h3hd732793hd904h0573hhhh0d83",
			salt:           "djiasddj",
			expected:       false,
		},
	}
	for _, tt := range testcases {
		name := fmt.Sprintf("check_password__password:%s,salt:%s,expected:%v", tt.password, tt.salt, tt.expected)
		t.Run(name, func(t *testing.T) {
			m := New(tt.salt)
			assert.Equal(t, tt.expected, m.CheckPassword(tt.password, tt.hashedPassword))
		})
	}
}
