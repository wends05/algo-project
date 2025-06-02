extends Node3D

class_name Level

var level_index: int
var is_end_level: bool = false

@export var doors_node: Doors
@export var backdoor_node: Backdoor

func _ready() -> void:
	# find index level from the name
	var level_number = int(name.split(" ")[1])
	level_index = Game.get_index_of_level(level_number)
	print("Level index: ", level_index)
	

	# find the doors
	doors_node = $Doors
	if not doors_node:
		print("Doors node not found in level %d." % level_index)
		return
	
	if level_index != 0:
		backdoor_node = $Backdoor
		if not backdoor_node and level_index > 0:
			print("Backdoor node not found in level %d." % level_index)
			return
		initialize_backdoor()
	else:
		backdoor_node.queue_free()
		print("No backdoor for level %d." % level_index)
	
	if level_index != Game.end_vertex:
		is_end_level = false
		initialize_doors()
	else:
		is_end_level = true
		print("[Level] Level is the end")
	
	update_node_visibilities()
	
func update_node_visibilities():
	if is_end_level:
		
		for node in get_tree().get_nodes_in_group("removed_end"):
			node.queue_free()
		doors_node.queue_free()
	else:
		for node in get_tree().get_nodes_in_group("removed_game"):
			node.queue_free()
	pass

func initialize_doors():
	var edges: Array = Game.get_doors(level_index)
	
	if not edges:
		print("No doors found for level %d." % level_index)
	doors_node.initialize_doors(edges)

func initialize_backdoor():
	var backdoor_edge = Game.get_backdoor_edge(level_index)
	backdoor_node.set_edge(backdoor_edge)
