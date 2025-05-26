extends Node3D
class_name Backdoor

var edge: Dictionary = {}

func _on_player_detector_body_entered(body: Node3D) -> void:
	if body is Player:
		Game.backward(edge["weight"])

func set_edge(edge_data: Dictionary) -> void:
	edge = edge_data
	if edge:
		$WeightLabel.text = "%s" % int(edge["weight"])
		$ToLabel.text = "%s" % (int(edge["source"]) + 1)
	else:
		print("Edge data is not set for this backdoor.")
