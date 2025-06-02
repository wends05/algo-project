extends Node

# game settings
var number_of_vertices: int = 10
var start_vertex: int = 0
var end_vertex: int = 9

var weight_range: int = 15

# api results
var _graph: Dictionary = {}
var _bellman_ford_result: Dictionary = {}


# current game progress
var randomized_levels: Array = []
var progress: Array[int] = [start_vertex]
var start_energy := 0
var current_energy: int = start_energy

signal progress_changed(progress: Array[int])
signal energy_changed(energy: int)

func get_goal_path() -> Array:
	if _bellman_ford_result and "path" in _bellman_ford_result:
		return _bellman_ford_result.path
	else:
		print("Bellman-Ford result is not available.")
		return []

func _ready() -> void:
	randomized_levels = randomize_levels()

	print("Game initialized:")
	print("Randomized Levels:", randomized_levels)

#
# start functions
#

func restart():
	progress = [start_vertex]
	progress_changed.emit(progress)

	current_energy = start_energy
	energy_changed.emit(current_energy)
	

func randomize_levels():
	randomized_levels = range(1, number_of_vertices + 1)
	randomized_levels.shuffle()

	# ensure that the end vertex is not the first one
	if randomized_levels[0] == end_vertex:
		var temp = randomized_levels[0]
		randomized_levels[0] = randomized_levels[1]
		randomized_levels[1] = temp
	
	return randomized_levels

func get_index_of_level(level_number: int) -> int:
	return randomized_levels.find(level_number)

func get_level_at_index(level_index: int) -> int:
	return randomized_levels[level_index]

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

#
# in game functions
#

func forward(to: int, energy_consumed: int):
	progress.append(to)
	progress_changed.emit(progress)

	current_energy = current_energy - energy_consumed
	energy_changed.emit(current_energy)

	print("\n=====\nFORWARD")
	print("Progress:", progress)
	print("Current Energy:", current_energy)
	print("=====\n")

	var level_to = get_level_at_index(to)
	Utils.change_scene("res://scenes/levels/level_%s.tscn" % level_to)

func backward(returned_energy: int):
	progress.pop_back()
	progress_changed.emit(progress)

	current_energy = current_energy + returned_energy
	energy_changed.emit(current_energy)

	var to = progress[-1] if progress.size() > 0 else start_vertex
	
	print("\n=====\nBACKWARD")
	print("Progress:", progress)
	print("To:", to)
	print("Current Energy:", current_energy)
	print("=====\n")
	
	var level_to = get_level_at_index(to)
	Utils.change_scene("res://scenes/levels/level_%s.tscn" % level_to)

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

#
# end game functions
#

func calculate_score():
	var goal_energy = get_goal_energy()

	if not goal_energy:
		print("Goal energy is not set, cannot calculate score.")
		return 0
	var score: int = current_energy / goal_energy * 100
	print("Score calculated: %s" % score)
	return score

func reset_game():
	print("Resetting game...")
	restart()
	_graph = {}
	_bellman_ford_result = {}
	randomize_levels()

func get_goal_energy():
	if _bellman_ford_result and "distance" in _bellman_ford_result:
		return start_energy - _bellman_ford_result.distance
