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
		str(Game.get_goal_path()),
		Game.progress
	]


func _on_restart_area_body_entered(body: Node3D) -> void:
	if body is Player:
		Game.restart()

func _on_exit_area_body_entered(body: Node3D) -> void:
	if body is Player:
		pass
