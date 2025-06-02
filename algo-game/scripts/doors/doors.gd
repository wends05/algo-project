extends Node3D
class_name Doors


func initialize_doors(edges: Array):
	for idx in range(edges.size()):
		var edge = edges[idx]
		
		var door_scene : PackedScene = load("res://scenes/doors/Door.tscn")
		var door_instance : Door = door_scene.instantiate()
		door_instance.name = "Door %d" % edge["target"]
		door_instance.edge = edge

		## find door position based on doors node
		var node_to_place: Node3D = get_child(edges.find(edge))

		if node_to_place:
			door_instance.global_transform = node_to_place.global_transform
		else:
			continue
		add_child(door_instance)
	for child in get_children():
		if child is not Door:
			child.queue_free()
