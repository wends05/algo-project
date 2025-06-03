extends Node3D


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	$"%Statistics Details".text = \
	"""
	Total Energy: %s
	Goal Energy: %s

	Goal Path: %s
	Path Passed: %s

	""" % [
		Game.current_energy,
		Game.get_goal_energy(),
		Game.get_goal_path().map(func(x): return x + 1),
		Game.progress.map(func(x): return x + 1)
	]


func _on_restart_area_body_entered(body: Node3D) -> void:
	if body is Player:
		Game.restart()
		Utils.change_scene("res://scenes/levels/level_%d.tscn" % Game.randomized_levels[0])

func _on_exit_area_body_entered(body: Node3D) -> void:
	if body is Player:
		Game.reset_game()
		Utils.change_scene("res://scenes/main/start.tscn")


func _on_restart_area_2_body_entered(body: Node3D) -> void:
	# reset map too
	if body is Player:
		Utils.bellman_ford_success.connect(_start_game)
		print_debug("Restarting game... with another map")
		Game.reset_game()
		Utils.generate_graph()

func _start_game(data):
	if data:
		Utils.change_scene("res://scenes/levels/level_%d.tscn" % Game.randomized_levels[0])
