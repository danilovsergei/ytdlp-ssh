package gnome

import (
	"errors"
	"fmt"
	"log"

	"github.com/godbus/dbus/v5"
	gnome_keyring "github.com/ppacher/go-dbus-keyring"
)

const chromeSaveStorage = "chromium Safe Storage"

func GetChromeKeyringPassword() ([]byte, error) {
	// what is d-bus @https://dbus.freedesktop.org/
	var chromeSecret []byte
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	svc, err := gnome_keyring.GetSecretService(conn)
	if err != nil {
		return nil, err
	}
	session, err := svc.OpenSession()
	if err != nil {
		return nil, err
	}
	defer func() {
		session.Close()
	}()
	collections, err := svc.GetAllCollections()
	if err != nil {
		return nil, err
	}
	for _, col := range collections {
		items, err := col.GetAllItems()
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			label, err := item.GetLabel()
			if err != nil {
				log.Fatalln(err)
				continue
			}
			// TODO Check it when running gnome
			if label == chromeSaveStorage {
				se, err := item.GetSecret(session.Path())
				if err != nil {
					log.Fatalln(err)
					return nil, err
				}
				fmt.Println(string(se.Value))
				chromeSecret = se.Value
			}
		}
	}
	if chromeSecret == nil {
		return nil, errors.New("dbus secret is empty")
	}
	return chromeSecret, nil
}
