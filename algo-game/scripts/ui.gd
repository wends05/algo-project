extends Control

class_name UserInterface

signal paused
signal unpaused

func _process(delta: float) -> void:
	$Bounds/VBoxContainer/RichTextLabel.text = str(Engine.get_frames_per_second())

func pause():
	Global.shift_locked = false
	$"%Pause Screen".visible = true
	get_tree().paused = true
	paused.emit()

func unpause():
	get_tree().paused = false
	$"%Pause Screen".visible = false
	unpaused.emit()

func _on_pause_button_pressed() -> void:
	print("[UI] Pause button pressed")
	pause()

func _on_return_pressed() -> void:
	unpause()


func _on_restart_pressed() -> void:
	Game.restart()


func _on_quit_pressed() -> void:
	Utils.change_scene("res://scenes/main/main_menu.tscn")
