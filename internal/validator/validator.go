package validator

type FieldError struct {
	Field   string
	Message string
}

type Validator struct {
	Errors []FieldError
}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) Add(field string, message string) {
	if v.Get(field) != "" {
		return
	}

	v.Errors = append(v.Errors, FieldError{Field: field, Message: message})
}

func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) Get(field string) string {
	for _, e := range v.Errors {
		if e.Field == field {
			return e.Message
		}
	}
	return ""
}

func (v *Validator) Map() map[string]string {
	m := make(map[string]string, len(v.Errors))
	for _, e := range v.Errors {
		m[e.Field] = e.Message
	}
	return m
}
