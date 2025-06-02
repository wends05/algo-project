extends Control

@export var player : Player

signal transition_finished

func _ready() -> void:
	player.can_move = false
	check_game_status()
	transition_finished.connect(_on_transition_finished)

func _on_transition_finished():
	player.can_move = true

func check_game_status():
	$%Transition.modulate = Color8(255, 255, 255, 0)

	if Game.calculate_score() == 100:
		$%Label.text = "You win"
		$%Background.color = Color8(255, 255, 255, 150)
		win_game_transition()
	else:
		$%Label.text = "You lose"
		$%Background.color = Color8(0, 0, 0, 150)
		lose_game_transition()

func lose_game_transition():
	var tween = create_tween()
	$%Label.modulate = Color8(255, 255, 255)
	tween.tween_property($%Transition, "modulate:a", 1, 1.2)
	tween.tween_interval(1.2)
	tween.tween_property($%Transition, "modulate:a", 0, 1.2)
	await tween.finished
	$%Transition.visible = false
	transition_finished.emit()

func win_game_transition():
	var tween = create_tween()
	tween.tween_property($%Transition, "modulate:a", 1, 1.2)
	tween.tween_interval(1.2)
	tween.tween_property($%Transition, "modulate:a", 0, 1.2)
	await tween.finished
	$%Transition.visible = false
	transition_finished.emit()
