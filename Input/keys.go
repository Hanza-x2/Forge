package Input

import "github.com/go-gl/glfw/v3.3/glfw"

type Key int

const (
	KeyUnknown Key = iota
	KeySpace
	KeyApostrophe
	KeyComma
	KeyMinus
	KeyPeriod
	KeySlash
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeySemicolon
	KeyEqual
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyLeftBracket
	KeyBackslash
	KeyRightBracket
	KeyGraveAccent
	KeyWorld1
	KeyWorld2
	KeyEscape
	KeyEnter
	KeyTab
	KeyBackspace
	KeyInsert
	KeyDelete
	KeyRight
	KeyLeft
	KeyDown
	KeyUp
	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd
	KeyCapsLock
	KeyScrollLock
	KeyNumLock
	KeyPrintScreen
	KeyPause
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyF25
	KeyKP0
	KeyKP1
	KeyKP2
	KeyKP3
	KeyKP4
	KeyKP5
	KeyKP6
	KeyKP7
	KeyKP8
	KeyKP9
	KeyKPDecimal
	KeyKPDivide
	KeyKPMultiply
	KeyKPSubtract
	KeyKPAdd
	KeyKPEnter
	KeyKPEqual
	KeyLeftShift
	KeyLeftControl
	KeyLeftAlt
	KeyLeftSuper
	KeyRightShift
	KeyRightControl
	KeyRightAlt
	KeyRightSuper
	KeyMenu
	KeyLast = KeyMenu
)

func convertGLFWKey(glfwKey glfw.Key) Key {
	switch glfwKey {
	// Printable keys
	case glfw.KeySpace:
		return KeySpace
	case glfw.KeyApostrophe:
		return KeyApostrophe
	case glfw.KeyComma:
		return KeyComma
	case glfw.KeyMinus:
		return KeyMinus
	case glfw.KeyPeriod:
		return KeyPeriod
	case glfw.KeySlash:
		return KeySlash
	case glfw.Key0:
		return Key0
	case glfw.Key1:
		return Key1
	case glfw.Key2:
		return Key2
	case glfw.Key3:
		return Key3
	case glfw.Key4:
		return Key4
	case glfw.Key5:
		return Key5
	case glfw.Key6:
		return Key6
	case glfw.Key7:
		return Key7
	case glfw.Key8:
		return Key8
	case glfw.Key9:
		return Key9
	case glfw.KeySemicolon:
		return KeySemicolon
	case glfw.KeyEqual:
		return KeyEqual

	// Alphabet keys
	case glfw.KeyA:
		return KeyA
	case glfw.KeyB:
		return KeyB
	case glfw.KeyC:
		return KeyC
	case glfw.KeyD:
		return KeyD
	case glfw.KeyE:
		return KeyE
	case glfw.KeyF:
		return KeyF
	case glfw.KeyG:
		return KeyG
	case glfw.KeyH:
		return KeyH
	case glfw.KeyI:
		return KeyI
	case glfw.KeyJ:
		return KeyJ
	case glfw.KeyK:
		return KeyK
	case glfw.KeyL:
		return KeyL
	case glfw.KeyM:
		return KeyM
	case glfw.KeyN:
		return KeyN
	case glfw.KeyO:
		return KeyO
	case glfw.KeyP:
		return KeyP
	case glfw.KeyQ:
		return KeyQ
	case glfw.KeyR:
		return KeyR
	case glfw.KeyS:
		return KeyS
	case glfw.KeyT:
		return KeyT
	case glfw.KeyU:
		return KeyU
	case glfw.KeyV:
		return KeyV
	case glfw.KeyW:
		return KeyW
	case glfw.KeyX:
		return KeyX
	case glfw.KeyY:
		return KeyY
	case glfw.KeyZ:
		return KeyZ

	// Brackets/accents
	case glfw.KeyLeftBracket:
		return KeyLeftBracket
	case glfw.KeyBackslash:
		return KeyBackslash
	case glfw.KeyRightBracket:
		return KeyRightBracket
	case glfw.KeyGraveAccent:
		return KeyGraveAccent

	// Function keys
	case glfw.KeyEscape:
		return KeyEscape
	case glfw.KeyEnter:
		return KeyEnter
	case glfw.KeyTab:
		return KeyTab
	case glfw.KeyBackspace:
		return KeyBackspace
	case glfw.KeyInsert:
		return KeyInsert
	case glfw.KeyDelete:
		return KeyDelete

	// Arrow keys
	case glfw.KeyRight:
		return KeyRight
	case glfw.KeyLeft:
		return KeyLeft
	case glfw.KeyDown:
		return KeyDown
	case glfw.KeyUp:
		return KeyUp

	// Navigation keys
	case glfw.KeyPageUp:
		return KeyPageUp
	case glfw.KeyPageDown:
		return KeyPageDown
	case glfw.KeyHome:
		return KeyHome
	case glfw.KeyEnd:
		return KeyEnd

	// Lock keys
	case glfw.KeyCapsLock:
		return KeyCapsLock
	case glfw.KeyScrollLock:
		return KeyScrollLock
	case glfw.KeyNumLock:
		return KeyNumLock
	case glfw.KeyPrintScreen:
		return KeyPrintScreen
	case glfw.KeyPause:
		return KeyPause

	// Function keys
	case glfw.KeyF1:
		return KeyF1
	case glfw.KeyF2:
		return KeyF2
	case glfw.KeyF3:
		return KeyF3
	case glfw.KeyF4:
		return KeyF4
	case glfw.KeyF5:
		return KeyF5
	case glfw.KeyF6:
		return KeyF6
	case glfw.KeyF7:
		return KeyF7
	case glfw.KeyF8:
		return KeyF8
	case glfw.KeyF9:
		return KeyF9
	case glfw.KeyF10:
		return KeyF10
	case glfw.KeyF11:
		return KeyF11
	case glfw.KeyF12:
		return KeyF12
	case glfw.KeyF13:
		return KeyF13
	case glfw.KeyF14:
		return KeyF14
	case glfw.KeyF15:
		return KeyF15
	case glfw.KeyF16:
		return KeyF16
	case glfw.KeyF17:
		return KeyF17
	case glfw.KeyF18:
		return KeyF18
	case glfw.KeyF19:
		return KeyF19
	case glfw.KeyF20:
		return KeyF20
	case glfw.KeyF21:
		return KeyF21
	case glfw.KeyF22:
		return KeyF22
	case glfw.KeyF23:
		return KeyF23
	case glfw.KeyF24:
		return KeyF24
	case glfw.KeyF25:
		return KeyF25

	// Keypad keys
	case glfw.KeyKP0:
		return KeyKP0
	case glfw.KeyKP1:
		return KeyKP1
	case glfw.KeyKP2:
		return KeyKP2
	case glfw.KeyKP3:
		return KeyKP3
	case glfw.KeyKP4:
		return KeyKP4
	case glfw.KeyKP5:
		return KeyKP5
	case glfw.KeyKP6:
		return KeyKP6
	case glfw.KeyKP7:
		return KeyKP7
	case glfw.KeyKP8:
		return KeyKP8
	case glfw.KeyKP9:
		return KeyKP9
	case glfw.KeyKPDecimal:
		return KeyKPDecimal
	case glfw.KeyKPDivide:
		return KeyKPDivide
	case glfw.KeyKPMultiply:
		return KeyKPMultiply
	case glfw.KeyKPSubtract:
		return KeyKPSubtract
	case glfw.KeyKPAdd:
		return KeyKPAdd
	case glfw.KeyKPEnter:
		return KeyKPEnter
	case glfw.KeyKPEqual:
		return KeyKPEqual

	// Modifier keys
	case glfw.KeyLeftShift:
		return KeyLeftShift
	case glfw.KeyLeftControl:
		return KeyLeftControl
	case glfw.KeyLeftAlt:
		return KeyLeftAlt
	case glfw.KeyLeftSuper:
		return KeyLeftSuper
	case glfw.KeyRightShift:
		return KeyRightShift
	case glfw.KeyRightControl:
		return KeyRightControl
	case glfw.KeyRightAlt:
		return KeyRightAlt
	case glfw.KeyRightSuper:
		return KeyRightSuper
	case glfw.KeyMenu:
		return KeyMenu

	default:
		return KeyUnknown
	}
}
