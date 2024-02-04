package main

import "chsir-zy/msg-push-center/impl"

func main() {
	// hub := impl.NewHub()
	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	impl.ServeWs(hub, w, r)
	// })

	// http.ListenAndServe(":8080", nil)

	impl.NewServer()
}
