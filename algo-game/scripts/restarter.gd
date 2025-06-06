extends Area3D


@export var spawn: Vector3

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	# attempt to use the initial position of the player node
	if get_parent_node_3d():
		var player := get_parent_node_3d().find_child("Player")
		
		assert(is_instance_of(player, CharacterBody3D), "the player assigned is not a player")
		
		spawn = player.position
	if !spawn:
		printerr("There is no spawn created")

func _on_body_entered(body: Node3D) -> void:
	if body is Player:
		body.position = spawn
