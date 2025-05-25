extends Node3D

func _on_area_3d_area_entered(area: Area3D) -> void:
	if area.get_parent() is Player:
		print("IS PLAYER")
		utils.changeScene("res://scenes/levels/level_3.tscn")
