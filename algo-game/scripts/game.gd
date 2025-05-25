extends Node

# game settings
var number_of_vertices: int = 7
var start_vertex: int = 0
var end_vertex: int = 6

var weight_range: int = 15
var max_energy: int = 100

# api results
var _graph: Dictionary = {}
var _bellman_ford_result: Dictionary = {}

# current game progress
var progress: Array[int] = [start_vertex]
var current_level: int = 0
var current_energy: int = max_energy

func forward(to: int, energy_consumed: int):
	progress.append(to)
	current_level = to
	current_energy = current_energy - energy_consumed

	print("\n=====\nFORWARD")
	print("Progress:", progress)
	print("Current Level:", current_level)
	print("Current Energy:", current_energy)
	print("=====\n")

	Utils.change_scene("res://scenes/levels/level_%s.tscn" % (to + 1))

func backward(returned_energy: int):
	progress.pop_back()
	current_level = progress[-1] if progress.size() > 0 else start_vertex
	current_energy = current_energy + returned_energy
	
	print("\n=====\nBACKWARD")
	print("Progress:", progress)
	print("Current Level:", current_level)
	print("Current Energy:", current_energy)
	print("=====\n")
	
	Utils.change_scene("res://scenes/levels/level_%s.tscn" % (current_level + 1))

func set_graph(data) -> void:
	for edge in data["edges"]:
		edge["source"] = int(edge["source"])
		edge["target"] = int(edge["target"])
		edge["weight"] = int(edge["weight"])
	_graph = data

func get_graph():
	return _graph

func set_bellman_ford_result(data) -> void:
	var converted_path = []
	for v in data["path"]:
		converted_path.append(int(v))
	data["path"] = converted_path
	_bellman_ford_result = data

func get_bellman_ford_result():
	return _bellman_ford_result

func get_doors(level_index: int):
	var edges = _graph["edges"]
	var level_doors = edges.filter(func(edge): return edge["source"] == level_index)
	return level_doors

func get_backdoor_edge(level_index: int):
	var edges = _graph["edges"]
	var previous_level_idx = Game.progress[-2]
	print("Finding edge from source %s to target %s" % [previous_level_idx, level_index])
	var valid_edges = edges.filter(func(edge): return edge["source"] == previous_level_idx and edge["target"] == level_index)
	print_debug(valid_edges)

	## ensure only one edge is possible

	if valid_edges.size() != 1:
		print("Error: Expected exactly one edge from level %s to %s, found %d edges." % [previous_level_idx, level_index, valid_edges.size()])
		return null
	if valid_edges.size() == 0:
		print("Error: No edge found from level %s to %s." % [previous_level_idx, level_index])
		return null
	print("Found edge from level %s to %s: %s" % [previous_level_idx, level_index, valid_edges[0]])
	if valid_edges.size() > 1:
		print("Warning: More than one edge found from level %s to %s, returning the first one." % [previous_level_idx, level_index])
	# If there are multiple edges, return the first one
	# This is a fallback, ideally we should ensure the graph structure allows only one edge between two vertices.

	return valid_edges[0]
