package custom_libp2p

func (l *CustomLibP2P) GetHelloMessage(userName string) (string, error) {
	return "Hello " + userName + " this is a dummy function!", nil
}
