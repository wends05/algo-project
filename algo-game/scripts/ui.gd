extends Control

class_name UserInterface

signal paused
signal unpaused

func _ready():
	Game.progress_changed.connect(_on_progress_changed)
	_on_progress_changed(Game.progress)

	Game.energy_changed.connect(_on_energy_changed)
	_on_energy_changed(Game.current_energy)

func _on_progress_changed(progress: Array[int]) -> void:
	var progress_text = "Progress: %s" % [str(progress.map(
		func(idx): return idx + 1
	))]
	print("[UI] Progress changed:", progress_text)
	$"%Progress Label".text = progress_text


func _on_energy_changed(energy: int) -> void:
	$"%Energy Label".text = "⚡ %s " % energy

func _process(_delta: float) -> void:
	$"%FPS Text".text = str(Engine.get_frames_per_second())

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
	unpause()
	Game.restart()
	Utils.change_scene("res://scenes/levels/level_%d.tscn" % Game.randomized_levels[0])

func _on_quit_pressed() -> void:
	unpause()
	Game.reset_game()
	Utils.change_scene("res://scenes/main/start.tscn")

signal fade_in_finished
signal fade_out_finished

func fade_in_transition() -> void:
	var tween = create_tween()
	$Transition/Background.color = Color(0, 0, 0, 0)
	$Transition.visible = true
	tween.tween_property($"Transition/Background", "color", Color8(0,0,0,120), 1.2)
	await tween.finished
	fade_in_finished.emit()

func fade_out_transition() -> void:
	var tween = create_tween()
	tween.tween_property($"Transition/Background", "color", Color8(0,0,0,0), 1.2)
	await tween.finished
	$Transition.visible = false
	fade_out_finished.emit()
