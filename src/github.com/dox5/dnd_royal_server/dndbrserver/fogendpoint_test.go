package main_test

type fakeGetter struct {
    values map [string]string
}

func (f fakeGetter) FormValue(key string) string {
    if value, ok := f.values[key]; ok {
        return value
    } else {
        return ""
    }
}

func FakeFormValueGetter() fakeGetter {
    return fakeGetter{
        values: make(map [string]string) }
}
