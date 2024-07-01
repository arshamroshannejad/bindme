package bindme

type Map map[string]string

type Validator struct {
	Errors Map
}

func New() *Validator {
	return &Validator{
		Errors: make(Map),
	}
}

func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) Add(key, value string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = value
	}
}

func (v *Validator) Check(ok bool, key, value string) {
	if !ok {
		v.Add(key, value)
	}
}

func (v *Validator) In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}
