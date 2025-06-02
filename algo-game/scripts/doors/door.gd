extends Node3D

class_name Door

var edge: Dictionary = {}

func _ready():
	if edge:
		print("Edge data is now set for this door.")

		var to_label := (int(edge["target"]) + 1)
		var to_end := to_label == Game.end_vertex + 1
		$WeightLabel.text = "%s" % int(edge["weight"] * -1)

		if to_end:
			$ToLabel.text = "end"
			$"%Light".visible = true
		else:
			$ToLabel.text = "%s" % to_label
	else:
		print("Edge data is not set for this door.")

func _on_player_detector_body_entered(body: Node3D) -> void:
	if body is Player:
		print("Door entered: ", JSON.stringify(edge, "\t"))
		Game.forward(edge["target"], edge["weight"])
