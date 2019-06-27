package globals

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"strings"
)

type R2Registry struct{}

// ErrorCode : Used to provide error codes with textual descriptions
// ErrorCode -> Int : Negative if error is problematic.
type ErrorCode struct {
	Code   int
	Reason string
}

// Status : Expands upon true/false return values by providing error information.
type Status struct {
	Valid bool
	Error ErrorCode
}

func (r2 *R2Registry) getAssociatedKey(access uint32) *registry.Key {
	regKey, err := registry.OpenKey(registry.CLASSES_ROOT, "ror2mm\\shell\\open\\command", access)
	if err != nil {
		regKey, _, err = registry.CreateKey(registry.CLASSES_ROOT, "ror2mm\\shell\\open\\command", access)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
	}
	return &regKey
}

// SetAssociatedProtocol : Check if the ror2mm protocol is associated with r2modman
func (r2 *R2Registry) IsAssociatedWithProtocol() *Status {
	regKey := r2.getAssociatedKey(registry.READ)
	if regKey == nil {
		// Prevent popups if failure to create key
		return &Status{false, ErrorCode{-1, "Registry key could not be created"}}
	}
	defer regKey.Close()
	value, _, readErr := regKey.GetStringValue("")
	if readErr != nil {
		return &Status{false, ErrorCode{-2, readErr.Error()}}
	}
	if strings.Contains(strings.ToLower(value), "r2modman.exe") {
		return &Status{true, ErrorCode{0, "Is Associated"}}
	} else {
		return &Status{false, ErrorCode{1, "Is not associated"}}
	}
}

// SetAssociatedProtocol : Associate the ror2mm protocol with r2modman
func (r2 *R2Registry) SetAssociatedProtocol() *Status {
	regKey := r2.getAssociatedKey(registry.SET_VALUE)
	if regKey == nil {
		// Prevent popups if failure to create key
		return &Status{false, ErrorCode{-1, "Registry key could not be created"}}
	}
	defer regKey.Close()
	setErr := regKey.SetStringValue("", "\""+ExecutableLocation+"\" \"%1\"")
	if setErr != nil {
		return &Status{false, ErrorCode{-2, setErr.Error()}}
	}
	return &Status{true, ErrorCode{0, "Key set"}}
}
