extends Node


var shift_locked: bool = false
var in_game: bool = false

var number_of_vertices: int = 7
var start_vertex: int = 0
var end_vertex: int = 6

var weight_range: int = 14

var _graph = {}

var _bellman_ford_result = {}



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
