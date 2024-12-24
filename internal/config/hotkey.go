package config

import (
	"errors"
	"strings"

	"github.com/jonesrussell/godo/internal/logger"
	"golang.design/x/hotkey"
)

// Modifier constants for hotkeys
const (
	ModCtrl  = hotkey.Modifier(1 << 2) // Control key
	ModShift = hotkey.Modifier(1 << 0) // Shift key
	ModAlt   = hotkey.Modifier(1 << 3) // Alt/Option key
	ModSuper = hotkey.Modifier(1 << 6) // Super/Windows/Command key
)

// Hotkey represents a keyboard shortcut
type Hotkey string

// HotkeyConfig holds hotkey configuration
type HotkeyConfig struct {
	QuickNote Hotkey `mapstructure:"quick_note"`
}

type HotkeyMapper struct {
	log    logger.Logger
	modMap map[string]hotkey.Modifier
	keyMap map[string]hotkey.Key
}

func NewHotkeyMapper(log logger.Logger) *HotkeyMapper {
	return &HotkeyMapper{
		log: log,
		modMap: map[string]hotkey.Modifier{
			"Ctrl":  ModCtrl,
			"Alt":   ModAlt,
			"Shift": ModShift,
			"Super": ModSuper,
		},
		keyMap: map[string]hotkey.Key{
			"A":      hotkey.KeyA,
			"B":      hotkey.KeyB,
			"C":      hotkey.KeyC,
			"D":      hotkey.KeyD,
			"E":      hotkey.KeyE,
			"F":      hotkey.KeyF,
			"G":      hotkey.KeyG,
			"H":      hotkey.KeyH,
			"I":      hotkey.KeyI,
			"J":      hotkey.KeyJ,
			"K":      hotkey.KeyK,
			"L":      hotkey.KeyL,
			"M":      hotkey.KeyM,
			"N":      hotkey.KeyN,
			"O":      hotkey.KeyO,
			"P":      hotkey.KeyP,
			"Q":      hotkey.KeyQ,
			"R":      hotkey.KeyR,
			"S":      hotkey.KeyS,
			"T":      hotkey.KeyT,
			"U":      hotkey.KeyU,
			"V":      hotkey.KeyV,
			"W":      hotkey.KeyW,
			"X":      hotkey.KeyX,
			"Y":      hotkey.KeyY,
			"Z":      hotkey.KeyZ,
			"Space":  hotkey.KeySpace,
			"Return": hotkey.KeyReturn,
			"Escape": hotkey.KeyEscape,
			"0":      hotkey.Key0,
			"1":      hotkey.Key1,
			"2":      hotkey.Key2,
			"3":      hotkey.Key3,
			"4":      hotkey.Key4,
			"5":      hotkey.Key5,
			"6":      hotkey.Key6,
			"7":      hotkey.Key7,
			"8":      hotkey.Key8,
			"9":      hotkey.Key9,
		},
	}
}

func (h *HotkeyMapper) GetModifier(name string) (hotkey.Modifier, bool) {
	mod, ok := h.modMap[name]
	return mod, ok
}

func (h *HotkeyMapper) GetKey(name string) (hotkey.Key, bool) {
	key, ok := h.keyMap[name]
	return key, ok
}

// Parse converts the hotkey string into modifiers and key
func (h Hotkey) Parse(mapper *HotkeyMapper) ([]hotkey.Modifier, hotkey.Key, error) {
	if h == "" {
		err := errors.New("empty hotkey string")
		mapper.log.Error("empty hotkey string")
		return nil, 0, err
	}

	parts := strings.Split(string(h), "+")
	if len(parts) < 2 {
		err := errors.New("invalid hotkey format")
		mapper.log.Error("invalid hotkey format", "hotkey", h)
		return nil, 0, err
	}

	// Last part is the key
	keyStr := parts[len(parts)-1]
	key, ok := mapper.GetKey(keyStr)
	if !ok {
		err := errors.New("invalid key")
		mapper.log.Error("invalid key", "key", keyStr)
		return nil, 0, err
	}

	// Rest are modifiers
	var mods []hotkey.Modifier
	for _, modStr := range parts[:len(parts)-1] {
		mod, ok := mapper.GetModifier(modStr)
		if !ok {
			err := errors.New("invalid modifier")
			mapper.log.Error("invalid modifier", "modifier", modStr)
			return nil, 0, err
		}
		mods = append(mods, mod)
	}

	return mods, key, nil
}

// IsValid checks if the hotkey string is valid
func (h Hotkey) IsValid(mapper *HotkeyMapper) bool {
	_, _, err := h.Parse(mapper)
	return err == nil
}

// NewHotkey creates a new Hotkey from modifiers and key
func NewHotkey(mapper *HotkeyMapper, mods []hotkey.Modifier, key hotkey.Key) Hotkey {
	var parts []string

	// Add modifiers in the expected order: Ctrl, Shift, Alt
	modNames := []string{"Ctrl", "Shift", "Alt"}
	for _, modName := range modNames {
		mod, _ := mapper.GetModifier(modName)
		for _, m := range mods {
			if m == mod {
				parts = append(parts, modName)
				break
			}
		}
	}

	// Add key
	for keyStr, k := range mapper.keyMap {
		if k == key {
			parts = append(parts, keyStr)
			break
		}
	}

	return Hotkey(strings.Join(parts, "+"))
}

func (h Hotkey) String() string {
	return string(h)
}
