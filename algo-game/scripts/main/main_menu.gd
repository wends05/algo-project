extends Control


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	#$VBoxContainer/PlayButton.visible = false
	Utils.bellman_ford_success.connect(_bellman_ford_success)
	$VBoxContainer/Message.visible = false

func _on_button_pressed() -> void:
	Utils.generate_graph()
	$VBoxContainer/PlayButton.visible = false
	$VBoxContainer/Message.visible = true
	$VBoxContainer/Message.text = "Loading..."

func _bellman_ford_success(data):
	print(data)
	if data:
		Utils.change_scene("res://scenes/levels/level_%d.tscn" % Game.randomized_levels[0])
