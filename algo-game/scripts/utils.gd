extends Node


signal graph_generated(graph: Variant)
signal bellman_ford_success(result: Variant)

func change_scene(path: String):
	get_tree().call_deferred("change_scene_to_file", path)
	
var base_url: String = "http://localhost:8080"


func generate_graph():
	var http := HTTPRequest.new()
	http.request_completed.connect(_generate_graph_success)
	add_child(http)

	var url = \
  "%s/graph?number_of_vertices=%s&weight_range=%s" % \
	[
		base_url,
	  Game.number_of_vertices,
	  Game.weight_range
	]

	var err := http.request(url, [], HTTPClient.METHOD_GET)
	if err != OK:
		push_error("Error requesting graph data: %s" % err)
	

func _generate_graph_success(_result, _response_code, _headers, body):
	var json = JSON.new()
	json.parse(body.get_string_from_utf8())
	var response = json.get_data()

	Game.set_graph(response)
	bellman_ford()
	graph_generated.emit(response)

func bellman_ford():
	var http := HTTPRequest.new()
	http.request_completed.connect(_bellman_ford_success)
	add_child(http)

	var url = \
		"%s/bellmanford" % [base_url]
	var graph = Game.get_graph()
	if not graph:
		print("Graph data is not available.")
		return
	var edges = graph["edges"]

	var data = {
		"number_of_vertices": Game.number_of_vertices,
		"start_vertex": Game.start_vertex,
		"end_vertex": Game.end_vertex,
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

	var error = false
	print(JSON.stringify(response, "\t"))
	if response.has("error"):
		push_error("Error in Bellman-Ford response: %s" % response["error"])
		error = true
	
	if error:
		# try again
		print("trying again")
		generate_graph()
		return
	
	Game.set_bellman_ford_result(response)
	bellman_ford_success.emit(response)
