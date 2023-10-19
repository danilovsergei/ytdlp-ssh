package kwallet

import (
	"github.com/godbus/dbus"
)

const (
	busName      = "org.kde.kwalletd5"
	objectPath   = "/modules/kwalletd5"
	serviceName  = "kdewallet"
	appId        = "keyring"
	chromeFolder = "Chrome Keys"
	chromeKey    = "Chrome Safe Storage"
)

type KWalletDbus struct{ dbus.BusObject }

// Raw command to get kde network wallet name: kdewallet
// dbus-send --session --print-reply=literal --dest=org.kde.kwalletd5 /modules/kwalletd5 'org.kde.KWallet.networkWallet'
// Raw command to get kde keyring password
// kwallet-query --read-password 'Chrome Safe Storage --folder "Chrome Keys" kdewallet'
func GetChromeKeyringPassword() ([]byte, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	kwallet := KWalletDbus{conn.Object(busName, objectPath)}

	handle, err := kwallet.open()
	if err != nil {
		return nil, err
	}

	password, err := kwallet.readPassword(handle)
	if err != nil {
		return nil, err
	}
	return password, nil
}

func (k *KWalletDbus) open() (int32, error) {
	call := k.Call("org.kde.KWallet.open", 0, serviceName, int64(0), appId)
	if call.Err != nil {
		return 0, call.Err
	}
	return call.Body[0].(int32), call.Err
}

func (k *KWalletDbus) readPassword(handle int32) ([]byte, error) {
	call := k.Call("org.kde.KWallet.readPassword", 0, handle, chromeFolder, chromeKey, appId)
	if call.Err != nil {
		return []byte{}, call.Err
	}
	password := call.Body[0].(string)
	return []byte(password), call.Err
}
