extends Node


var shift_locked: bool = false
var in_game: bool = false

var number_of_vertices: int = 7
var start_vertex = 0
var end_vertex = 0

var weight_range: int = 14

var _graph = {}


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	pass

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(_delta: float) -> void:
	pass

func http_request_complete(_result, _response_code, _headers, body):
	var json = JSON.new()
	json.parse(body.get_string_from_utf8())
	var response = json.get_data()
	

	printt("Graph data received: %s" % response)

func set_graph(data):
	_graph = data
	print_debug("Graph edited: ", self._graph)

func get_graph():
	return _graph
