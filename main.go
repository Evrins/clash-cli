package main

func main() {
	client := NewClient("172.16.44.123:9090", "571a959d945830805ec3c963ca3ed675e9046601")
	client.Traffic()
}
