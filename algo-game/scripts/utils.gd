extends Node

func change_scene(path: String):
	get_tree().change_scene_to_file(path)


func _ready() -> void:
	generate_graph()

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
	if err:
		print("Error requesting graph data: %s" % err)
	
	remove_child(http) # Remove the HTTPRequest node after the request is made

func _generate_graph_success(_result, _response_code, _headers, body):
	var json = JSON.new()
	json.parse(body.get_string_from_utf8())
	var response = json.get_data()
  
	printt("Graph data received: %s" % response)
	global.set_graph(response)

	bellman_ford()
  
  # Emit a signal or call a function to handle the graph data
  # For example, you could emit a signal here to notify other parts of the game
  # that the graph data is ready.
  # emit_signal("graph_data_received", response)


func bellman_ford():
	var http := HTTPRequest.new()
	http.request_completed.connect(_bellman_ford_success)
	add_child(http)
	var url = \
		"http://localhost:8080/bellman_ford"
	
	
	var edges = global.get_graph().find_key("edges")

	var data = {
		"start_vertex": global.start_vertex,
		"end_vertex": global.end_vertex,
		"edges": edges
	}
	# Convert the edges to a format suitable for the request
	for edge in edges:
		data[edge] = global.get_graph()[edge]
	# Convert the data to a JSON string
	var json = JSON.stringify(data)
	print_debug("Graph edges: ", data)

	var err := http.request(url, [], HTTPClient.METHOD_POST, json)
	if err:
		print("Error requesting Bellman-Ford data: %s" % err)

func _bellman_ford_success(_result, _response_code, _headers, body):
	var json = JSON.new()
	json.parse(body.get_string_from_utf8())
	var response = json.get_data()
	
	printt("Bellman-Ford result", response)	
