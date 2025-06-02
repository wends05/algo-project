extends Node3D

class_name EndDoor


func _on_player_detector_body_entered(body: Node3D) -> void:
	if body is Player:
		body.can_move = false
		var ui = body.UI
		ui.fade_in_finished.connect(_change_scene)
		ui.fade_in_transition()


func _change_scene() -> void:
	# This function should be connected to the fade_out_finished signal
	# and will change the scene to the end screen or next level.
	print("[EndDoor] Changing scene to end screen.")
	Utils.change_scene("res://scenes/end_game/level_end.tscn")
