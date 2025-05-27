extends Node3D

class_name Door

var edge: Dictionary = {}

func _ready():
	if edge:
		print("Edge data is now set for this door.")
		
		$WeightLabel.text = "%s" % int(edge["weight"] * -1)
		$ToLabel.text = "%d" % (int(edge["target"]) + 1)
	else:
		print("Edge data is not set for this door.")

func _on_player_detector_body_entered(body: Node3D) -> void:
	if body is Player:
		print("Door entered: ", JSON.stringify(edge, "\t"))
		Game.forward(edge["target"], edge["weight"])
