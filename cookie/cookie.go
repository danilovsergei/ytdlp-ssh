//	Based on amazing code and research in
//
// https://github.com/teocci/go-chrome-cookies
package cookie

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"ytlpd-ssh/cookie/gnome"
	"ytlpd-ssh/cookie/kwallet"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/pbkdf2"
)

type cookie struct {
	Host         string
	Path         string
	KeyName      string
	Value        string
	IsSecure     bool
	IsHTTPOnly   bool
	HasExpire    bool
	IsPersistent bool
	CreateDate   time.Time
	ExpireDate   time.Time
}

type KeyringType string

const keyringKwallet KeyringType = "KWALLET"
const keyringGnome KeyringType = "GNOME"
const keyringText KeyringType = "TEXT"

var keyring = keyringKwallet

type cookies struct {
	cookies map[string][]cookie
}

const QueryChromiumCookie = `SELECT name, encrypted_value, host_key, path, creation_utc, expires_utc, is_secure, is_httponly, has_expires, is_persistent FROM cookies`

type OutputFormat int

// Finds chrome profile folder by provided email
func findChromeProfile(email string) (string, error) {
	// TODO write function to
	return "", nil
}

func CopyToLocalPath(src, dst string) error {
	locals, _ := filepath.Glob("*")
	for _, v := range locals {
		if v == dst {
			err := os.Remove(dst)
			if err != nil {
				return err
			}
		}
	}
	sourceFile, err := os.ReadFile(src)
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = os.WriteFile(dst, sourceFile, 0777)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return err
}

func (c *cookies) CopyDB(chromeCookiesDb string) error {
	return CopyToLocalPath(chromeCookiesDb, "/tmp/chrome_cookies")
}

func (c *cookies) ChromeParse() error {
	c.CopyDB(cookieStorePath)

	c.cookies = make(map[string][]cookie)
	cookieDB, err := sql.Open("sqlite3", "/tmp/chrome_cookies")
	if err != nil {
		return err
	}
	defer func() {
		if err := cookieDB.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	rows, err := cookieDB.Query(QueryChromiumCookie)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	for rows.Next() {
		var (
			key, host, path                               string
			isSecure, isHTTPOnly, hasExpire, isPersistent int
			createDate, expireDate                        int64
			value, encryptValue                           []byte
		)
		err = rows.Scan(&key, &encryptValue, &host, &path, &createDate, &expireDate, &isSecure, &isHTTPOnly, &hasExpire, &isPersistent)
		if err != nil {
			log.Fatalln(err)
		}
		cookie := cookie{
			Host:       host,
			IsSecure:   IntToBool(isSecure),
			HasExpire:  IntToBool(hasExpire),
			ExpireDate: TimeEpochFormat(expireDate),
			Path:       path,
			KeyName:    key,
		}
		value, err = decryptCookie(encryptValue)
		if err != nil {
			log.Fatalln(err)
		}
		cookie.Value = string(value)
		c.cookies[host] = append(c.cookies[host], cookie)
	}
	return nil
}

func decryptCookie(encryptedCookie []byte) ([]byte, error) {
	v10Key := deriveKey([]byte("peanuts"))
	emptyKey := deriveKey([]byte(""))
	if len(encryptedCookie) == 0 {
		return []byte(""), nil
	}

	version := string(encryptedCookie[:3])
	encryptedCookieValue := encryptedCookie[3:]
	if version == "v10" {
		return aes128CBCDecrypt(encryptedCookieValue, [][]byte{v10Key, emptyKey})
	} else if version == "v11" {
		var keyringPass []byte
		var err error

		if keyring == keyringGnome {
			keyringPass, err = gnome.GetChromeKeyringPassword()
		} else if keyring == keyringKwallet {
			keyringPass, err = kwallet.GetChromeKeyringPassword()
		}
		if err != nil {
			return nil, err
		}
		v11Key := deriveKey(keyringPass)
		return aes128CBCDecrypt(encryptedCookieValue, [][]byte{v11Key, emptyKey})
	}
	// Probably its possible to just proceed here
	// Fail for now to see what cookies are part of it
	log.Fatalf("Failed to detect cookie version %s \n", version)
	return nil, nil
}

func deriveKey(chromeSecret []byte) []byte {
	return pbkdf2.Key(chromeSecret, []byte("saltysalt"), 1, 16, sha1.New)
}

func aes128CBCDecrypt(encryptedCookie []byte, secretKeys [][]byte) ([]byte, error) {
	var initializationVector = []byte{32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32}
	for _, key := range secretKeys {
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		dst := make([]byte, len(encryptedCookie))
		mode := cipher.NewCBCDecrypter(block, initializationVector)
		mode.CryptBlocks(dst, encryptedCookie)
		dst = PKCS5UnPadding(dst)
		if len(dst) > 0 {
			return dst, nil
		}
	}
	return nil, errors.New("failed to decrypt a cookie. Provided keys could be wrong")
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpad := int(src[length-1])
	if length >= unpad {
		return src[:(length - unpad)]
	} else {
		return []byte{}
	}
}

// func DPApi(data []byte) ([]byte, error) {
// 	return nil, nil
// }

func IntToBool(a int) bool {
	switch a {
	case 0, -1:
		return false
	}
	return true
}
func TimeEpochFormat(epoch int64) time.Time {
	maxTime := int64(99633311740000000)
	if epoch > maxTime {
		return time.Date(2049, 1, 1, 1, 1, 1, 1, time.Local)
	}
	t := time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)
	d := time.Duration(epoch)
	for i := 0; i < 1000; i++ {
		t = t.Add(d)
	}
	return t
}

// Outputs cookies in HTTP Cookie File format required by ytdlp
func (c *cookies) outPut() string {
	out := strings.Builder{}
	appendString := func(strs ...string) {
		for _, str := range strs {
			out.WriteString(str)
			out.WriteString("\n")
		}
	}

	appendString(
		"# Netscape HTTP Cookie File",
		"# https://curl.haxx.se/rfc/cookie_spec.html",
		"# This is a generated file! Do not edit.")
	for host, value := range c.cookies {
		if host == ".youtube.com" {
			for _, cookie := range value {
				hasSubdomains := strings.ToUpper(fmt.Sprintf("%t", strings.HasPrefix(cookie.Host, ".")))
				isSecure := strings.ToUpper(fmt.Sprintf("%t", cookie.IsSecure))
				var expireDate = "0"
				if cookie.HasExpire {
					expireDate = fmt.Sprint(cookie.ExpireDate.Unix())
				}

				arr := strings.Join([]string{
					cookie.Host,
					hasSubdomains,
					cookie.Path,
					isSecure,
					expireDate,
					cookie.KeyName,
					cookie.Value,
				}, "\t")
				appendString(arr)
			}
		}
	}
	return out.String()
}

func ParseCookies() string {
	c := cookies{}
	c.ChromeParse()
	return c.outPut()
}
