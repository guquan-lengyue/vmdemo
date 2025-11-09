package kvm

import "testing"

func TestUsb(t *testing.T) {
	usbList, err := GetUsbList()
	if err != nil {
		t.Errorf("failed to get usb list: %s", err)
	}
	t.Log(usbList)
}
