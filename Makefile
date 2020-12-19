.PHONY: lambda/build lambda/zip

lambda/build:
	GOOS=linux go build -o bin/lambdasavefriendslist cmd/lambdasavefriendslist/main.go
	GOOS=linux go build -o bin/lambdagetfriendslist cmd/lambdagetfriendslist/main.go
	GOOS=linux go build -o bin/lambdadraw cmd/lambdadraw/main.go
	GOOS=linux go build -o bin/lambdageneratesecretkey cmd/lambdageneratesecretkey/main.go
	GOOS=linux go build -o bin/lambdagetsecretfriend cmd/lambdagetsecretfriend/main.go

lambda/zip:
	zip -j bin/lambdaSaveFriendsList.zip bin/lambdasavefriendslist
	zip -j bin/lambdaGetFriendsList.zip bin/lambdagetfriendslist
	zip -j bin/lambdaDraw.zip bin/lambdadraw
	zip -j bin/lambdaGenerateSecretKey.zip bin/lambdageneratesecretkey
	zip -j bin/lambdaGetSecretFriend.zip bin/lambdagetsecretfriend
