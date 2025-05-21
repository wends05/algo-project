extends Area3D


@export var spawn: Vector3

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	# check if there is a spawn variable
	if spawn == Vector3.ZERO:
		print("No spawn point set for restarter")
	
	# attempt to use the initial position of the player node
	if get_parent_node_3d():
		var player := get_parent_node_3d().find_child("Player")
		
		assert(is_instance_of(player, CharacterBody3D), "the player assigned is not a player")
		
		spawn = player.position
	if !spawn:
		printerr("There is no spawn created")

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass


func _on_body_entered(body: Node3D) -> void:
	print(body)
	if body is Player:
		print_debug("Player detectedw")
		body.position = spawn
