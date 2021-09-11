package main

func main() {
	NewRouter(
		NewBCWebAPI(
			NewBrightnessController(),
		),
	).Start()
}
