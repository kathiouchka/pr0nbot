build:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o pbot

transfer: build
	scp -i awskey.pem pbot ec2-user@ec2-35-180-41-226.eu-west-3.compute.amazonaws.com:~/pbot

run: transfer
	ssh -i awskey.pem ec2-user@ec2-35-180-41-226.eu-west-3.compute.amazonaws.com "~/./pbot -t \$$pronbotkey &"