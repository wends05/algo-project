extends Node3D

class_name Level

var level_index: int

@export var doors_node: Doors
@export var backdoor_node: Backdoor

func _ready() -> void:
	# find index level from the name
	level_index = int(name.split(" ")[1]) - 1
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
		backdoor_node = null
		print("No backdoor for level %d." % level_index)
	
	initialize_doors()


func initialize_doors():
	var edges: Array = Game.get_doors(level_index)
	
	if not edges:
		print("No doors found for level %d." % level_index)
	doors_node.initialize_doors(edges)

func initialize_backdoor():
	var backdoor_edge = Game.get_backdoor_edge(level_index)
	backdoor_node.set_edge(backdoor_edge)
