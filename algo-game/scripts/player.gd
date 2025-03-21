extends CharacterBody3D

# How fast the player moves in meters per second.
@export var speed = 14
# The downward acceleration when in the air, in meters per second squared.
@export var fall_acceleration = 75

@export var jump_velocity = 20

var gravity = ProjectSettings.get_setting("physics/3d/default_gravity")

var target_velocity = Vector3.ZERO

func _physics_process(delta):
	# We create a local variable to store the input direction.
	
	if not is_on_floor():
		velocity.y -=  gravity * delta
		
	if Input.is_action_just_pressed("jump") and is_on_floor():
		velocity.y = jump_velocity
		
	
	var input_dir := Input.get_vector("ui_left", "ui_right", "move_forward", "move_back")
	var direction := (transform.basis * Vector3(input_dir.x, 0, input_dir.y)).normalized()
	
	if direction:
		velocity.x = direction.x * speed
		velocity.z = direction.z * speed
	else:
		velocity.x = move_toward(velocity.x, 0, speed)
		velocity.z = move_toward(velocity.z, 0, speed)
		
