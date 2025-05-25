extends CharacterBody3D

class_name Player

const SPEED = 5.0
const JUMP_VELOCITY = 4.5
const SENSITIVITY = 0.003

var gravity = ProjectSettings.get_setting("physics/3d/default_gravity")

@onready var head = $Head
@onready var head_camera = $Head/Camera3D

func _ready() -> void:
	# fix camera
	head_camera.rotation_degrees.x = -15.0


func _input(event) -> void:
	if event.is_action_pressed("shift_lock"):
		Global.shift_locked = !Global.shift_locked
func _unhandled_input(event: InputEvent) -> void:
	if event is InputEventMouseMotion:
		if event.button_mask & MOUSE_BUTTON_MASK_LEFT or Global.shift_locked:
			# rotate body
			rotate_y(-event.relative.x * SENSITIVITY)
			
			# rotate head
			head_camera.rotate_x(-event.relative.y * SENSITIVITY)
			
			head_camera.rotation.x = clamp(head_camera.rotation.x, deg_to_rad(-90), deg_to_rad(90))

		if Global.shift_locked:
			Input.mouse_mode = Input.MOUSE_MODE_CAPTURED
		else:
			Input.mouse_mode = Input.MOUSE_MODE_VISIBLE

func _physics_process(delta: float) -> void:
	if not is_on_floor():
		velocity.y -= gravity * delta
		
	if Input.is_action_just_pressed("jump") and is_on_floor():
		velocity.y = JUMP_VELOCITY
		
	
	var input_dir = Input.get_vector("move_left", "move_right", "move_forward", "move_back")
	var direction = (transform.basis * Vector3(input_dir.x, 0, input_dir.y)).normalized()
	
	if direction:
		velocity.x = direction.x * SPEED
		velocity.z = direction.z * SPEED
	else:
		velocity.x = 0
		velocity.z = 0
	
	move_and_slide()
