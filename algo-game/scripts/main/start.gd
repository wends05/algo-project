extends Node3D

var player_spawn_position: Vector3 = Vector3(0, 0, 0)

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	Utils.bellman_ford_success.connect(_bellman_ford_success)
	player_spawn_position = $Player.global_transform.origin

func _bellman_ford_success(data):
	print(data)
	if data:
		Utils.change_scene("res://scenes/levels/level_%d.tscn" % Game.randomized_levels[0])


func _on_start_area_body_entered(body: Node3D) -> void:
	if body is Player:
		Utils.generate_graph()
		$"%Start Label".text = "Loading Game..."

func _on_instructions_area_body_entered(body: Node3D) -> void:
	if body is Player:
		var return_door : Node3D = $"%Return Door"
		var return_door_position : Vector3 = return_door.global_transform.origin
		var return_door_rotation : Vector3 = return_door.global_transform.basis.get_euler()
		
		var teleport_position := Vector3(
			return_door_position.x - 2,
			return_door_position.y + 1.0,
			return_door_position.z + 2
		)

		var teleport_rotation := Vector3(
			return_door_rotation.x,
			return_door_rotation.y + PI,
			return_door_rotation.z
		)

		body.global_transform.origin = teleport_position
		body.global_transform.basis = Basis.from_euler(teleport_rotation)

func _on_return_area_body_entered(body: Node3D) -> void:
	if body is Player:
		body.global_transform.origin = player_spawn_position
		body.global_transform.basis = Basis.IDENTITY
