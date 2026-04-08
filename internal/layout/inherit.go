package layout

import (
	"reflect"

	"github.com/Viswesh934/gotei/internal/style"
)

// inheritAllProperties forces text nodes to inherit ALL inheritable properties from parent
// per W3C spec: text nodes inherit all cascading properties
func inheritAllProperties(currentStyle, parentStyle style.Style) style.Style {
	current := reflect.ValueOf(&currentStyle).Elem()
	parent := reflect.ValueOf(parentStyle)

	for fieldName, isInherit := range style.InheritableFields {
		if !isInherit {
			continue // Skip non-inheritable properties
		}

		currentField := current.FieldByName(fieldName)
		parentField := parent.FieldByName(fieldName)

		if !currentField.IsValid() || !parentField.IsValid() {
			continue
		}

		// For text nodes, inherit all inheritable properties unconditionally
		if currentField.CanSet() {
			currentField.Set(parentField)
		}
	}

	return currentStyle
}

func shouldInherit(currentStyle, parentStyle style.Style) style.Style {
	// Use W3C inheritable properties map from style package
	current := reflect.ValueOf(&currentStyle).Elem()
	parent := reflect.ValueOf(parentStyle)

	for fieldName, isInherit := range style.InheritableFields {
		if !isInherit {
			continue // Skip non-inheritable properties
		}

		currentField := current.FieldByName(fieldName)
		parentField := parent.FieldByName(fieldName)

		if !currentField.IsValid() || !parentField.IsValid() {
			continue
		}

		// Get default value for this field from DefaultStyle
		defaultStyle := reflect.ValueOf(style.DefaultStyle)
		defaultField := defaultStyle.FieldByName(fieldName)

		// Inherit if current equals default AND parent differs from default
		if reflect.DeepEqual(currentField.Interface(), defaultField.Interface()) &&
			!reflect.DeepEqual(parentField.Interface(), defaultField.Interface()) {
			currentField.Set(parentField)
		}
	}

	return currentStyle
}
