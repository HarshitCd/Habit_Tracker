#! /bin/bash

mkdir -p ~/.config/habit_tracker
cp ./config/config.toml ~/.config/habit_tracker/config.toml
cp ./template/template.toml ~/.config/habit_tracker/tracker.toml

go build -o ~/mybin/habit_tracker

alias habit_tracker=~/mybin/habit_tracker
echo "alias habit_tracker=~/mybin/habit_tracker" >> ~/.zshrc