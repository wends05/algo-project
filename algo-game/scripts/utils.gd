extends Node

func change_scene(path: String):
	get_tree().change_scene_to_file(path)


# func _ready() -> void:
# 	generate_graph()

func generate_graph():
	var http := HTTPRequest.new()
	http.request_completed.connect(_generate_graph_success)
	add_child(http)
	var url = \
  "http://localhost:8080/graph?number_of_vertices=%s&weight_range=%s" % \
	[
	  global.number_of_vertices,
	  global.weight_range
	]
	var err := http.request(url, [], HTTPClient.METHOD_GET)
	if err != OK:
		push_error("Error requesting graph data: %s" % err)
	

func _generate_graph_success(_result, _response_code, _headers, body):
	var json = JSON.new()
	json.parse(body.get_string_from_utf8())
	var response = json.get_data()

	global.set_graph(response)

	bellman_ford()

func bellman_ford():
	var http := HTTPRequest.new()
	http.request_completed.connect(_bellman_ford_success)
	add_child(http)
	var url = \
		"http://localhost:8080/bellmanford"
	
	var graph = global.get_graph()
	if not graph:
		print("Graph data is not available.")
		return
	var edges = graph["edges"]

	var data = {
		"number_of_vertices": global.number_of_vertices,
		"start_vertex": global.start_vertex,
		"end_vertex": global.end_vertex,
		"edges": edges
	}

	# Convert the data to a JSON string
	var body = JSON.stringify(data)
	var err := http.request(url, [], HTTPClient.METHOD_POST, body)
	if err != OK:
		printerr("Error requesting Bellman-Ford data: %s" % err)
		push_error("Error requesting graph data: %s" % err)
	print(err)

func _bellman_ford_success(_result, _response_code, _headers, body):
	var json = JSON.new()
	json.parse(body.get_string_from_utf8())
	var response = json.get_data()
	
	global.set_bellman_ford_result(response)
