package wechat


type defaultHook struct {

}
func (hook defaultHook) ShortURLReadStorage (longURL string) (shortURL string, has bool) {
	return "", false
}
func (hook defaultHook) ShortURLWriteStorage (longURL string, shortURL string) {

}
func DefaultHook() defaultHook {
	return defaultHook{}
}
