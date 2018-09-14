# go-roulette

Go-roulette is represents api for any roulette-based things like (mini games, chatroullete etc.). It works on websockets, tokens and golang:)

Warning! All **"games"** in this readme file actually is not necessarily a games. It's modules that implements `game` interface.

### Files
All logic is contained in:`lib/` and `games/names.go`.<br>
Api is in `server.go`. It's very easy and well documented.<br>
Client library is in `public/assets/js/go-roulette.js`

### Run

Just change port in `router.Run()` in `server.go` and build this file.<br>
Example games will be accesible on `http://hostname:port/`

### How to use?

You can add this whole module(except example files) to your server and build.

Server-side module files should be in `games/%gamename%` folder and contain at least one file with same package name, that implements `game` interface.

Don't forget to add special case in `games/names.go` like in example:

```go 
case "chatroulette":
	return &chatroulette.Chat{}
 ```
 
Also create `public/%gamename%/main.html`, which you will be able to access on `http://hostname:port/%gamename%`.

For messaging with server use `public/assets/js/go-roulette.js`. It is easy to use and my examples shows how to work with it. I'm pretty bad in javascript, so it would be great if somebody rewritten this library considering all standards,