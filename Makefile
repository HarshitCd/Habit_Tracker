build:
	@go build -o run/habit_tracker

exe:
	@make build
	@./run/habit_tracker $(OPTION)

clean:
	@rm -rf run