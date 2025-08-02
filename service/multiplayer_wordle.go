package service

type MultiplayerHost struct {
	Send chan string
}

type MultiplayerClinet struct {
	Send chan string
}
